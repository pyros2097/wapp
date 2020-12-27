// +build !tinygo
// +build wasm

package app

import "github.com/pyros2097/wapp/js"

func (e *Element) setJsEventHandler(k string, h js.EventHandler) {
	jshandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// dispatch(func() {
		if !e.self().Mounted() {
			return nil
		}
		e := js.Event{
			Src:   this,
			Value: args[0],
		}
		trackMousePosition(e)
		h.Value(e)
		// })
		return nil
	})
	h.JSvalue = jshandler
	e.events[k] = h
	e.JSValue().Call("addEventListener", k, jshandler)
}

func (e *Element) delJsEventHandler(k string, h js.EventHandler) {
	e.JSValue().Call("removeEventListener", k, h.JSvalue)
	h.JSvalue.Release()
	delete(e.events, k)
}

func (e *Element) mount() error {
	if e.Mounted() {
		panic("mounting elem failed already mounted " + e.name())
	}

	v := js.Window.Get("document").Call("createElement", e.tag)
	if !v.Truthy() {
		panic("mounting component failed create javascript node returned nil " + e.name())
	}
	e.jsvalue = v

	for k, v := range e.attrs {
		e.setJsAttr(k, v)
	}

	for k, v := range e.events {
		e.setJsEventHandler(k, v)
	}

	for _, c := range e.children() {
		if err := e.appendChild(c, true); err != nil {
			panic("mounting component failed appendChild " + e.name())
		}
	}

	return nil
}
