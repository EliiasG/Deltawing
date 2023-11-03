package web

import (
	"github.com/eliiasg/deltawing/internal/setup/web"
	"github.com/eliiasg/deltawing/web/app"
)

// DOES NOT RETURN!
func InitApp(init func(app.App), update func()) {
	web.InitApp(init, update)
}
