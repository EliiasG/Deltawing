package web

import (
	"syscall/js"

	"github.com/eliiasg/deltawing/input"
)

type keyboard struct {
	keyPressedHandler   func(input.Key)
	keyReleaseddHandler func(input.Key)
	keyHeldHandler      func(input.Key)
	keyTypedHandler     func(rune)
}

func makeKeyboard() *keyboard {
	k := &keyboard{}
	doc := js.Global().Get("document")
	doc.Call("addEventListener", "keydown", js.FuncOf(k.jsKeypressed), false)
	return k
}

func (k *keyboard) jsKeypressed(_ js.Value, args []js.Value) any {
	key := args[0]
	k.handleKey(key)
	k.handleChar(key)
	return nil
}

func (k *keyboard) handleKey(key js.Value) {
	button := keyVal(key.Get("code").String())
	// key
	if key.Get("repeat").Bool() {
		if k.keyHeldHandler != nil {
			k.keyHeldHandler(button)
		}
	} else if k.keyPressedHandler != nil {
		k.keyPressedHandler(button)
	}
}

func (k *keyboard) handleChar(key js.Value) {
	if k.keyTypedHandler == nil {
		return
	}
	charVal := key.Get("key").String()
	// dont fire on empty
	if len(charVal) == 0 {
		return
	}
	var r rune
	for i, val := range charVal {
		r = val
		// only fire if single rune
		if i > 0 {
			return
		}
	}
	k.keyTypedHandler(r)
}

func (k *keyboard) SetKeyPressedHandler(handler func(input.Key)) {
	k.keyPressedHandler = handler
}

func (k *keyboard) SetKeyReleasedHandler(handler func(input.Key)) {
	k.keyReleaseddHandler = handler
}

func (k *keyboard) SetKeyHeldHandler(handler func(input.Key)) {
	k.keyPressedHandler = handler
}

func (k *keyboard) SetKeyTypedHandler(handler func(rune)) {
	k.keyTypedHandler = handler
}

var keyMap = map[string]input.Key{
	"Space":        input.KeySpace,
	"Quote":        input.KeyApostrophe,
	"Comma":        input.KeyComma,
	"Minus":        input.KeyMinus,
	"Period":       input.KeyPeriod,
	"Slash":        input.KeySlash,
	"Digit0":       input.Key0,
	"Digit1":       input.Key1,
	"Digit2":       input.Key2,
	"Digit3":       input.Key3,
	"Digit4":       input.Key4,
	"Digit5":       input.Key5,
	"Digit6":       input.Key6,
	"Digit7":       input.Key7,
	"Digit8":       input.Key8,
	"Digit9":       input.Key9,
	"Semicolon":    input.KeySemicolon,
	"Equal":        input.KeyEqual,
	"KeyA":         input.KeyA,
	"KeyB":         input.KeyB,
	"KeyC":         input.KeyC,
	"KeyD":         input.KeyD,
	"KeyE":         input.KeyE,
	"KeyF":         input.KeyF,
	"KeyG":         input.KeyG,
	"KeyH":         input.KeyH,
	"KeyI":         input.KeyI,
	"KeyJ":         input.KeyJ,
	"KeyK":         input.KeyK,
	"KeyL":         input.KeyL,
	"KeyM":         input.KeyM,
	"KeyN":         input.KeyN,
	"KeyO":         input.KeyO,
	"KeyP":         input.KeyP,
	"KeyQ":         input.KeyQ,
	"KeyR":         input.KeyR,
	"KeyS":         input.KeyS,
	"KeyT":         input.KeyT,
	"KeyU":         input.KeyU,
	"KeyV":         input.KeyV,
	"KeyW":         input.KeyW,
	"KeyX":         input.KeyX,
	"KeyY":         input.KeyY,
	"KeyZ":         input.KeyZ,
	"BracketLeft":  input.KeyLeftBracket,
	"Backslash":    input.KeyBackslash,
	"BracketRight": input.KeyRightBracket,
	"BackQuote":    input.KeyGraveAccent,
	// world 1 and 2 missing
	"Escape":     input.KeyEscape,
	"Enter":      input.KeyEnter,
	"Tab":        input.KeyTab,
	"Backspace":  input.KeyBackspace,
	"Insert":     input.KeyInsert,
	"Delet":      input.KeyDelete,
	"ArrowRight": input.KeyRight,
	"ArrowLeft":  input.KeyLeft,
	"ArrowDown":  input.KeyDown,
	"ArrowUp":    input.KeyUp,
	"PageUp":     input.KeyPageUp,
	"PageDown":   input.KeyPageDown,
	"Home":       input.KeyHome,
	"End":        input.KeyEnd,
	"CapsLock":   input.KeyCapsLock,
	"ScrollLock": input.KeyScrollLock,
	"NumLock":    input.KeyNumLock,
	// printscreen missing
	"Pause": input.KeyPause,
	"F1":    input.KeyF1,
	"F2":    input.KeyF2,
	"F3":    input.KeyF3,
	"F4":    input.KeyF4,
	"F5":    input.KeyF5,
	"F6":    input.KeyF6,
	"F7":    input.KeyF7,
	"F8":    input.KeyF8,
	"F9":    input.KeyF9,
	"F10":   input.KeyF10,
	"F11":   input.KeyF11,
	"F12":   input.KeyF12,
	"F13":   input.KeyF13,
	"F14":   input.KeyF14,
	"F15":   input.KeyF15,
	"F16":   input.KeyF16,
	"F17":   input.KeyF17,
	"F18":   input.KeyF18,
	"F19":   input.KeyF19,
	"F20":   input.KeyF20,
	"F21":   input.KeyF21,
	"F22":   input.KeyF22,
	"F23":   input.KeyF23,
	"F24":   input.KeyF24,
	// f25 missing
	"Numpad0":        input.KeyKp0,
	"Numpad1":        input.KeyKp1,
	"Numpad2":        input.KeyKp2,
	"Numpad3":        input.KeyKp3,
	"Numpad4":        input.KeyKp4,
	"Numpad5":        input.KeyKp5,
	"Numpad6":        input.KeyKp6,
	"Numpad7":        input.KeyKp7,
	"Numpad8":        input.KeyKp8,
	"Numpad9":        input.KeyKp9,
	"NumpadDecimal":  input.KeyKpDecimal,
	"NumpadDivide":   input.KeyKpDivide,
	"NumpadMultiply": input.KeyKpMultiply,
	"NumpadSubtract": input.KeyKpSubtract,
	"NumpadAdd":      input.KeyKpAdd,
	"NumpadEnter":    input.KeyKpEnter,
	"NumpadEqual":    input.KeyKpEqual,
	"ShiftLeft":      input.KeyLeftShift,
	"ControlLeft":    input.KeyLeftControl,
	"AltLeft":        input.KeyLeftAlt,
	"ShiftRight":     input.KeyRightShift,
	"ControlRight":   input.KeyRightControl,
	"AltRight":       input.KeyRightAlt,
	// super buttons missing
}

func keyVal(key string) input.Key {
	r, ok := keyMap[key]
	if !ok {
		return input.KeyUnknown
	}
	return r
}
