package app

import "reflect"

func Html(elems ...UI) *Element {
	return &Element{tag: "html", body: elems}
}

func Head(elems ...UI) *Element {
	basic := []UI{
		&Element{tag: "meta", selfClosing: true, attrs: map[string]string{"charset": "UTF-8"}},
		&Element{tag: "meta", selfClosing: true, attrs: map[string]string{"http-equiv": "Content-Type", "content": "text/html;charset=utf-8"}},
		&Element{tag: "meta", selfClosing: true, attrs: map[string]string{"http-equiv": "encoding", "content": "utf-8"}},
	}
	return &Element{tag: "head", body: append(basic, elems...)}
}

func Body(elems ...UI) *Element {
	return &Element{tag: "body", body: elems}
}

func Title(v string) *Element {
	return &Element{tag: "title", body: []UI{Text(v)}}
}

func Meta(name, content string) *Element {
	e := &Element{
		tag:         "meta",
		selfClosing: true,
	}
	e.setAttr("name", name)
	e.setAttr("content", content)
	return e
}

func Link(rel, href string) *Element {
	e := &Element{
		tag:         "link",
		selfClosing: true,
	}
	e.setAttr("rel", rel)
	e.setAttr("href", href)
	return e
}

func Script(str string) *Element {
	return &Element{
		tag:  "script",
		body: []UI{Text(str)},
	}
}

func Div(uis ...interface{}) *Element {
	return NewElement("div", false, uis...)
}

func A(uis ...interface{}) *Element {
	return NewElement("a", false, uis...)
}

func P(uis ...interface{}) *Element {
	return NewElement("p", false, uis...)
}

func Span(uis ...interface{}) *Element {
	return NewElement("span", false, uis...)
}

func Input(uis ...interface{}) *Element {
	return NewElement("input", false, uis...)
}

func Image(uis ...interface{}) *Element {
	return NewElement("image", true, uis...)
}

func Button(uis ...interface{}) *Element {
	return NewElement("button", false, uis...)
}

func Svg(uis ...interface{}) *Element {
	return NewElement("svg", false, uis...)
}

func SvgText(uis ...interface{}) *Element {
	return NewElement("text", false, uis...)
}

func Ul(uis ...interface{}) *Element {
	return NewElement("ul", false, uis...)
}

func Li(uis ...interface{}) *Element {
	return NewElement("li", false, uis...)
}

func Row(uis ...interface{}) UI {
	return Div(append([]interface{}{Css("flex flex-row justify-center items-center")}, uis...)...)
}

func Col(uis ...interface{}) UI {
	return Div(append([]interface{}{Css("flex flex-col justify-center items-center")}, uis...)...)
}

func If(expr bool, a UI) UI {
	if expr {
		return a
	}
	return nil
}

func IfElse(expr bool, a UI, b UI) UI {
	if expr {
		return a
	}
	return b
}

func Map(source interface{}, f func(i int) UI) []UI {
	src := reflect.ValueOf(source)
	if src.Kind() != reflect.Slice {
		panic("range loop source is not a slice: " + src.Type().String())
	}
	body := make([]UI, 0, src.Len())
	for i := 0; i < src.Len(); i++ {
		body = append(body, f(i))
	}
	return body
}

func Map2(source interface{}, f func(v interface{}, i int) UI) []UI {
	src := reflect.ValueOf(source)
	if src.Kind() != reflect.Slice {
		panic("range loop source is not a slice: " + src.Type().String())
	}
	body := make([]UI, 0, src.Len())
	for i := 0; i < src.Len(); i++ {
		body = append(body, f(src.Index(i), i))
	}
	return body
}

// func (r RangeLoop) Slice(f func(int) UI) RangeLoop {
// 	src := reflect.ValueOf(r.source)
// 	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
// 		panic("range loop source is not a slice or array: " + src.Type().String())
// 	}

// 	body := make([]UI, 0, src.Len())
// 	for i := 0; i < src.Len(); i++ {
// 		body = append(body, FilterUIElems(f(i))...)
// 	}

// 	r.body = body
// 	return r
// }

// // Map sets the loop content by repeating the given function for the number
// // of elements in the source. Elements are ordered by keys.
// //
// // It panics if the range source is not a map or if map keys are not strings.
// func (r RangeLoop) Map(f func(string) UI) RangeLoop {
// 	src := reflect.ValueOf(r.source)
// 	if src.Kind() != reflect.Map {
// 		panic("range loop source is not a map: " + src.Type().String())
// 	}

// 	if keyType := src.Type().Key(); keyType.Kind() != reflect.String {
// 		panic("range loop source keys are not strings: " + src.Type().String() + keyType.String())
// 	}

// 	body := make([]UI, 0, src.Len())
// 	keys := make([]string, 0, src.Len())

// 	for _, k := range src.MapKeys() {
// 		keys = append(keys, k.String())
// 	}
// 	sort.Strings(keys)

// 	for _, k := range keys {
// 		body = append(body, FilterUIElems(f(k))...)
// 	}

// 	r.body = body
// 	return r
// }
