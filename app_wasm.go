// +build wasm

package app

import (
	"io"

	"github.com/pyros2097/wapp/js"
)

func Run() {
	handle, _, _ := globalRouter.Lookup("GET", js.Window.Location().Pathname)
	if handle == nil {
		renderFunc = globalRouter.NotFound
	} else {
		renderFunc, _ = handle.(RenderFunc)
	}
	defer func() {
		err := recover()
		// show alert
		panic(err)
	}()

	initBody()
	initContent()
	if err := body.replaceChildAt(0, renderFunc); err != nil {
		panic("replacing content failed")
	}
	content = renderFunc

	for {
		select {
		case f := <-uiChan:
			f()
		}
	}
}

func Reload() {
	dispatch(func() {
		js.Window.Location().Reload()
	})
}

func Route(path string, render RenderFunc) {
	globalRouter.GET(path, render)
}

func initBody() {
	body = &Element{
		jsvalue: js.Window.Get("document").Get("body"),
		tag:     "body",
	}
	body.setSelf(body)
}

func initContent() {
	content := &Element{
		jsvalue: body.JSValue().Get("firstElementChild"),
		tag:     "div",
	}

	content.setSelf(content)
	content.setParent(body)
	body.body = append(body.body, content)
}

// func onPopState(this js.Value, args []js.Value) interface{} {
// 	dispatch(func() {
// 		// navigate(Window().URL(), false)
// 	})
// 	return nil
// }
// func isExternalNavigation(u *url.URL) bool {
// 	return u.Host != "" && u.Host != js.Window().URL().Host
// }
// func isFragmentNavigation(u *url.URL) bool {
// 	return u.Fragment != ""
// }

func writePage(ui UI, w io.Writer) {
}
