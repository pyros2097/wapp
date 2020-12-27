// +build wasm
// +build tinygo

package app

import (
	njs "syscall/js"

	"github.com/pyros2097/wapp/js"
)

func (e *Element) setJsEventHandler(k string, h js.EventHandler) {
	// h.JSvalue = dd
	// e.events[k] = h
	e.JSValue().Call("addEventListener", k, njs.FuncOf(func(this njs.Value, args []njs.Value) interface{} {
		println("2")
		println("3")
		return nil
	}))
}

func (e *Element) delJsEventHandler(k string, h js.EventHandler) {
	// e.JSValue().Call("removeEventListener", k, h.JSvalue)
	// delete(e.events, k)
}

func HH(this njs.Value, args []njs.Value) interface{} {
	println("HH")
	return nil
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
