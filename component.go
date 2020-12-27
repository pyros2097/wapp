package app

import (
	"reflect"
	"strings"

	"github.com/pyros2097/wapp/js"
)

var contextMap = map[int]*RenderContext{}
var contextIndex = 0

func getCurrentContext() *RenderContext {
	return contextMap[contextIndex]
}

type RenderFunc func(ctx *RenderContext) UI

func (r RenderFunc) JSValue() js.Value {
	c := getCurrentContext()
	return c.root.JSValue()
}

func (r RenderFunc) Mounted() bool {
	c := getCurrentContext()
	return c.root != nil && c.root.Mounted() &&
		r.self() != nil
}

func (r RenderFunc) Render() UI {
	c := getCurrentContext()
	c.index = 0
	c.eindex = 0
	println("render")
	return r(c)
}

func (r RenderFunc) Update() {
	// dispatch(func() {
	if !r.Mounted() {
		return
	}
	println("update")
	if err := r.updateRoot(); err != nil {
		panic(err)
	}
	// })
}

func (r RenderFunc) name() string {
	name := reflect.TypeOf(r.self()).String()
	name = strings.ReplaceAll(name, "main.", "")
	return name
}

func (r RenderFunc) self() UI {
	c := getCurrentContext()
	return c.this
}

func (r RenderFunc) setSelf(n UI) {
	c := getCurrentContext()
	if n != nil {
		println("xnew render context")
		c := NewRenderContext()
		switch d := n.(type) {
		case RenderFunc:
			c.this = d
		default:
			panic("unknown type")
		}
		return
	}

	c.this = nil
}

func (r RenderFunc) attributes() map[string]string {
	return nil
}

func (r RenderFunc) eventHandlers() map[string]js.EventHandler {
	return nil
}

func (r RenderFunc) parent() UI {
	c := getCurrentContext()
	return c.parentElem
}

func (r RenderFunc) setParent(p UI) {
	c := getCurrentContext()
	c.parentElem = p
}

func (r RenderFunc) children() []UI {
	c := getCurrentContext()
	return []UI{c.root}
}

func (r RenderFunc) mount() error {
	c := getCurrentContext()
	if r.Mounted() {
		panic("mounting component failed already mounted " + r.name())
	}
	root := r.Render()
	if err := mount(root); err != nil {
		panic("mounting component failed " + r.name())
	}
	root.setParent(c.this)
	c.root = root
	return nil
}

func (r RenderFunc) dismount() {
	c := getCurrentContext()
	// for _, v := range c.effectsUnsub {
	// 	if v != nil {
	// 		v()
	// 	}
	// }
	dismount(c.root)
	delete(contextMap, c.contextMapIndex)
	contextIndex--
}

func (r RenderFunc) update(n UI) error {
	if r.self() == n || !r.Mounted() {
		return nil
	}

	if n.name() != n.name() {
		panic("updating ui element failed replace different element type current-name: " + r.name() + " updated-name: " + n.name())
	}

	aval := reflect.Indirect(reflect.ValueOf(r.self()))
	bval := reflect.Indirect(reflect.ValueOf(n))
	compotype := reflect.ValueOf(r).Elem().Type()

	for i := 0; i < aval.NumField(); i++ {
		a := aval.Field(i)
		b := bval.Field(i)

		if a.Type() == compotype {
			continue
		}

		if !a.CanSet() {
			continue
		}

		if !reflect.DeepEqual(a.Interface(), b.Interface()) {
			a.Set(b)
		}
	}

	return r.updateRoot()
}

func (r RenderFunc) updateRoot() error {
	c := getCurrentContext()
	a := c.root
	println("updateRoot")
	b := r.Render()

	err := update(a, b)
	if isErrReplace(err) {
		err = r.replaceRoot(b)
	}

	if err != nil {
		panic("updating component failed " + r.name())
	}

	return nil
}

func (r RenderFunc) replaceRoot(n UI) error {
	c := getCurrentContext()
	old := c.root
	new := n

	if err := mount(new); err != nil {
		panic("replacing component root failed name: " + r.name() + " root-name: " + old.name() + "new-root-name: " + new.name())
	}

	var parent UI
	for {
		parent = r.parent()
		_, isElem := parent.(*Element)
		if parent == nil || isElem {
			break
		}
	}

	if parent == nil {
		panic("replacing component root failed name: " + r.name() + " component does not have html element parents")
	}

	c.root = new
	new.setParent(r.self())

	oldjs := old.JSValue()
	newjs := n.JSValue()
	parent.JSValue().Call("replaceChild", newjs, oldjs)

	dismount(old)
	return nil
}

type RenderContext struct {
	contextMapIndex int
	parentElem      UI
	root            UI
	this            RenderFunc
	index           int
	values          map[int]interface{}
	eindex          int
	effects         map[int][]interface{}
	effectsUnsub    map[int]func()
}

func NewRenderContext() *RenderContext {
	c := &RenderContext{
		contextMapIndex: contextIndex,
		values:          map[int]interface{}{},
		effects:         map[int][]interface{}{},
		effectsUnsub:    map[int]func(){},
	}
	contextMap[contextIndex] = c
	// contextIndex++
	return c
}

func (c *RenderContext) UseState(initial interface{}) (func() interface{}, func(v interface{})) {
	i := c.index
	c.index++
	if _, ok := c.values[i]; !ok {
		c.values[i] = initial
	}
	return func() interface{} {
			return c.values[i]
		}, func(v interface{}) {
			c.values[i] = v
			// special check so that the backend doesn't crash
			if c.this != nil {
				c.this.Update()
			}
		}
}

func (c *RenderContext) UseInt(initial int) (func() int, func(v int)) {
	getState, setState := c.UseState(initial)
	return func() int {
			v, exist := getState().(int)
			if !exist {
				panic("v is not int")
			}
			return v
		}, func(v int) {
			setState(v)
		}
}

func (c *RenderContext) UseEffect(f func() func(), deps ...interface{}) {
	i := c.eindex
	c.eindex++
	if _, ok := c.effects[i]; !ok {
		println("initial deps")
		c.effects[i] = deps
		c.effectsUnsub[i] = f()
		return
	}
	hasChanged := false
	for di, ndv := range deps {
		odv := c.effects[i][di]
		if odv != ndv {
			c.effects[i] = deps
			hasChanged = true
			break
		}
	}
	println("hasChanged", hasChanged)
	if hasChanged {
		f()
	}
}

func (c *RenderContext) UseAtom(a *Atom) interface{} {
	c.UseEffect(func() func() {
		return a.Subscribe(func(v interface{}) {
			c.this.Update()
		})
	})
	return a.Get()
}

type Subscriber func(v interface{})

type Atom struct {
	value       interface{}
	subscribers []Subscriber
}

func NewAtom(v interface{}) *Atom {
	return &Atom{
		value: v,
	}
}

func (a *Atom) Subscribe(v Subscriber) func() {
	a.subscribers = append(a.subscribers, v)
	i := len(a.subscribers)
	return func() {
		a.subscribers = append(a.subscribers[:i], a.subscribers[i+1:]...)
	}
}

func (a *Atom) Get() interface{} {
	return a.value
}

func (a *Atom) Set(v interface{}) {
	a.value = v
	for _, s := range a.subscribers {
		s(v)
	}
}
