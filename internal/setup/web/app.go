//go:build wasm
// +build wasm

package web

import (
	"fmt"
	"runtime/debug"
	"strings"
	"syscall/js"

	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/render/gl"
	"github.com/eliiasg/deltawing/input"
	"github.com/eliiasg/deltawing/internal/rendering/webgl"
	"github.com/eliiasg/deltawing/web/app"
	"github.com/eliiasg/glow/enum"
)

type webApp struct {
	time       float64
	renderer   *gl.Renderer
	kb         *keyboard
	updateFunc func()
	jsUpdate   js.Func
	canvas     js.Value
}

func InitApp(init func(app.App), update func()) {
	initApp(init, update)
	select {}
}

// to resolve defer
func initApp(init func(app.App), update func()) {
	// get variables
	glob := js.Global()
	doc := glob.Get("document")
	body := doc.Get("body")
	win := glob.Get("window")
	// remove border
	body.Get("style").Call("setProperty", "margin", 0)
	body.Get("style").Call("setProperty", "overflow", "hidden")
	// init gl
	canvas := doc.Call("createElement", "canvas")
	g := initGl(doc, canvas)
	if g.IsNull() {
		panic("Failed to init webgl2, please try changing or updating your browser")
	}
	// init renderer
	r := gl.NewRenderer(
		func() uint16 {
			return uint16(win.Get("innerWidth").Int())
		},
		func() uint16 {
			return uint16(win.Get("innerHeight").Int())
		},
		webgl.MakeContext(g),
		"#version 300 es\nprecision highp float;\nprecision highp int;",
		true,
	)
	// init app
	a := &webApp{
		renderer:   r,
		updateFunc: update,
		kb:         &keyboard{},
		canvas:     canvas,
	}
	// catch init error
	defer a.requestUpdate()
	a.jsUpdate = js.FuncOf(a.update)
	init(a)
}

func initGl(doc js.Value, canvas js.Value) js.Value {
	doc.Get("body").Call("appendChild", canvas)
	params := make(map[string]any)
	params["antialias"] = false
	g := canvas.Call("getContext", "webgl2", params)
	g.Call("disable", enum.CULL_FACE)
	g.Call("enable", enum.BLEND)
	g.Call("enable", enum.DEPTH_TEST)
	g.Call("enable", enum.SAMPLE_COVERAGE)
	g.Call("enable", enum.SAMPLE_ALPHA_TO_COVERAGE)
	g.Call("depthFunc", enum.GREATER)
	g.Call("clearDepth", 0)
	return g
}

func (w *webApp) requestUpdate() {
	r := recover()
	glob := js.Global()
	if r != nil {
		message := fmt.Sprintf("%v\n%v", r, string(debug.Stack()))
		glob.Get("document").Get("body").Set("innerHTML", strings.ReplaceAll(message, "\n", "<br>"))
		panic(r)
	}
	win := glob.Get("window")
	win.Call("requestAnimationFrame", w.jsUpdate)
}

func (w *webApp) update(_ js.Value, args []js.Value) any {
	defer w.requestUpdate()
	win := js.Global().Get("window")
	width, height := win.Get("innerWidth").Int(), win.Get("innerHeight").Int()
	w.canvas.Set("width", width)
	w.canvas.Set("height", height)
	w.time = args[0].Float() * 0.001
	w.renderer.PrimaryRenderTarget().Resize(uint16(width), uint16(height))
	w.updateFunc()
	w.renderer.PrimaryRenderTarget().BlitTo(w.renderer.RealPrimaryRenderTarget(), 0, 0)
	return nil
}

func (w *webApp) Renderer() render.Renderer {
	return w.renderer
}

func (w *webApp) Time() float64 {
	return w.time
}

func (w *webApp) Keyboard() input.Keyboard {
	return w.kb
}

func (w *webApp) Mouse() input.Mouse {
	panic("unimplemented")
}

func (w *webApp) Controller() input.Controller {
	panic("unimplemented")
}
