// +build wasm
// +build tinygo

package js

import (
	"syscall/js"
)

var Window = &browserWindow{value: value{Value: js.Global()}}

type value struct {
	js.Value
}

type EventHandler struct {
	Event   string
	JSvalue js.Func
	Value   EventHandlerFunc
}

func NewEventHandler(e string, v EventHandlerFunc) EventHandler {
	return EventHandler{
		Event: e,
		Value: v,
	}
}

func (h EventHandler) Equal(o EventHandler) bool {
	return h.Event == o.Event
	// h.Value == o.Value
	// fmt.Sprintf("%p", h.Value) == fmt.Sprintf("%p", o.Value)
}

func NewValue(v js.Value) Value {
	return value{Value: v}
}

func (v value) Call(m string, args ...interface{}) Value {
	args = cleanArgs(args...)
	return val(v.Value.Call(m, args...))
}

func (v value) Get(p string) Value {
	return val(v.Value.Get(p))
}

func (v value) Set(p string, x interface{}) {
	if wrapper, ok := x.(Wrapper); ok {
		x = jsval(wrapper.JSValue())
	}
	v.Value.Set(p, x)
}

func (v value) Index(i int) Value {
	return val(v.Value.Index(i))
}

func (v value) InstanceOf(t Value) bool {
	return v.Value.InstanceOf(jsval(t))
}

func (v value) Invoke(args ...interface{}) Value {
	return val(v.Value.Invoke(args...))
}

func (v value) JSValue() Value {
	return v
}

func (v value) New(args ...interface{}) Value {
	args = cleanArgs(args...)
	return val(v.Value.New(args...))
}

func (v value) Type() Type {
	return Type(v.Value.Type())
}

func null() Value {
	return val(js.Null())
}

func undefined() Value {
	return val(js.Undefined())
}

func valueOf(x interface{}) Value {
	switch t := x.(type) {
	case value:
		x = t.Value

	case *browserWindow:
		x = t.Value

	case Event:
		return valueOf(t.Value)
	}

	return val(js.ValueOf(x))
}

type browserWindow struct {
	value

	cursorX int
	cursorY int
}

func (w *browserWindow) Size() (width int, height int) {
	getSize := func(axis string) int {
		size := w.Get("inner" + axis)
		if !size.Truthy() {
			size = w.
				Get("document").
				Get("documentElement").
				Get("client" + axis)
		}
		if !size.Truthy() {
			size = w.
				Get("document").
				Get("body").
				Get("client" + axis)
		}
		if size.Type() != TypeNumber {
			return 0
		}
		return size.Int()
	}

	return getSize("Width"), getSize("Height")
}

func (w *browserWindow) CursorPosition() (x, y int) {
	return w.cursorX, w.cursorY
}

func (w *browserWindow) SetCursorPosition(x, y int) {
	w.cursorX = x
	w.cursorY = y
}

func (w *browserWindow) GetElementByID(id string) Value {
	return w.Get("document").Call("getElementById", id)
}

func (w *browserWindow) ScrollToID(id string) {
	if elem := w.GetElementByID(id); elem.Truthy() {
		elem.Call("scrollIntoView")
	}
}

func (w *browserWindow) AddEventListener(event string, h EventHandler) func() {
	panic("not implemented")
}

func (w *browserWindow) Location() *Location {
	return &Location{
		value: w.Get("location"),
		Href: w.
			Get("location").
			Get("href").
			String(),
		Pathname: w.
			Get("location").
			Get("pathname").
			String(),
	}
}

func (l *Location) Reload() {
	l.value.Call("reload")
}

func val(v js.Value) Value {
	return value{Value: v}
}

func jsval(v Value) js.Value {
	switch v := v.(type) {
	case value:
		return v.Value

	case *browserWindow:
		return v.Value

	case Event:
		return jsval(v.Value)

	default:
		// + reflect.TypeOf(v).String()
		println("syscall/js value conversion failed type: ", v)
		return js.Undefined()
	}
}

// JSValue returns the underlying syscall/js value of the given Javascript
// value.
func JSValue(v Value) js.Value {
	return jsval(v)
}

func copyBytesToGo(dst []byte, src Value) int {
	return js.CopyBytesToGo(dst, jsval(src))
}

func copyBytesToJS(dst Value, src []byte) int {
	return js.CopyBytesToJS(jsval(dst), src)
}

func cleanArgs(args ...interface{}) []interface{} {
	for i, a := range args {

		args[i] = cleanArg(a)
	}

	return args
}

func cleanArg(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(v))
		for key, val := range v {
			m[key] = cleanArg(val)
		}
		return m

	case []interface{}:
		s := make([]interface{}, len(v))
		for i, val := range v {
			s[i] = cleanArgs(val)
		}

	case Wrapper:
		return jsval(v.JSValue())
	}

	return v

}
