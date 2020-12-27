// +build !wasm

package app

import "github.com/pyros2097/wapp/js"

func (e *Element) setJsEventHandler(k string, h js.EventHandler) {
}

func (e *Element) delJsEventHandler(k string, h js.EventHandler) {
}

func (e *Element) mount() error {
	return nil
}
