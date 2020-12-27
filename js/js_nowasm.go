// +build !wasm
// +build !tinygo

package js

var Window = &browserWindow{value: value{}}

type value struct{}

// Func is the interface that describes a wrapped Go function to be called by
// JavaScript.
type Func interface {
	Value

	// Release frees up resources allocated for the function. The function must
	// not be invoked after calling Release.
	Release()
}

type EventHandler struct {
	Event   string
	JSvalue Func
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

func NewValue(v interface{}) value {
	return value{}
}

func (v value) Bool() bool {
	panic("wasm required")
}

func (v value) Call(m string, args ...interface{}) Value {
	panic("wasm required")
}

func (v value) Float() float64 {
	panic("wasm required")
}

func (v value) Get(p string) Value {
	panic("wasm required")
}

func (v value) Index(i int) Value {
	panic("wasm required")
}

func (v value) InstanceOf(t Value) bool {
	panic("wasm required")
}

func (v value) Int() int {
	panic("wasm required")
}

func (v value) Invoke(args ...interface{}) Value {
	panic("wasm required")
}

func (v value) IsNaN() bool {
	panic("wasm required")
}

func (v value) IsNull() bool {
	panic("wasm required")
}

func (v value) IsUndefined() bool {
	panic("wasm required")
}

func (v value) JSValue() Value {
	panic("wasm required")
}

func (v value) Length() int {
	panic("wasm required")
}

func (v value) New(args ...interface{}) Value {
	panic("wasm required")
}

func (v value) Set(p string, x interface{}) {
	panic("wasm required")
}

func (v value) SetIndex(i int, x interface{}) {
	panic("wasm required")
}

func (v value) String() string {
	panic("wasm required")
}

func (v value) Truthy() bool {
	panic("wasm required")
}

func (v value) Type() Type {
	panic("wasm required")
}

func null() Value {
	panic("wasm required")
}

func undefined() Value {
	panic("wasm required")
}

func valueOf(x interface{}) Value {
	panic("wasm required")
}

func FuncOf(fn func(this Value, args []Value) interface{}) Func {
	panic("wasm required")
}

type browserWindow struct {
	value
}

func (w browserWindow) Size() (width, height int) {
	panic("wasm required")
}

func (w browserWindow) CursorPosition() (x, y int) {
	panic("wasm required")
}

func (w browserWindow) SetCursorPosition(x, y int) {
	panic("wasm required")
}

func (w *browserWindow) GetElementByID(id string) Value {
	panic("wasm required")
}

func (w *browserWindow) ScrollToID(id string) {
	panic("wasm required")
}

func (w *browserWindow) AddEventListener(event string, h EventHandlerFunc) func() {
	panic("wasm required")
}

func (w *browserWindow) Location() Location {
	panic("wasm required")
}

func copyBytesToGo(dst []byte, src Value) int {
	panic("wasm required")
}

func copyBytesToJS(dst Value, src []byte) int {
	panic("wasm required")
}

func (l *Location) Reload() {
	panic("wasm required")
}
