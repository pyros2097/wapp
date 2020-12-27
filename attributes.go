package app

import (
	"reflect"
	"strconv"

	"github.com/pyros2097/wapp/js"
)

type Attribute struct {
	Key   string
	Value string
}

func ID(v string) Attribute {
	return Attribute{"id", v}
}

func Style(v string) Attribute {
	return Attribute{"style", v}
}

func Accept(v string) Attribute {
	return Attribute{"accept", v}
}

func AutoComplete(v bool) Attribute {
	return Attribute{"autocomplete", strconv.FormatBool(v)}
}

func Checked(v bool) Attribute {
	return Attribute{"checked", strconv.FormatBool(v)}
}

func Disabled(v bool) Attribute {
	return Attribute{"disabled", strconv.FormatBool(v)}
}

func Name(v string) Attribute {
	return Attribute{"name", v}
}

func Type(v string) Attribute {
	return Attribute{"type", v}
}

func Value(v string) Attribute {
	return Attribute{"value", v}
}

func Placeholder(v string) Attribute {
	return Attribute{"placeholder", v}
}

func Src(v string) Attribute {
	return Attribute{"src", v}
}

func ViewBox(v string) Attribute {
	return Attribute{"viewBox", v}
}

func X(v string) Attribute {
	return Attribute{"x", v}
}

func Y(v string) Attribute {
	return Attribute{"y", v}
}

type CssAttribute struct {
	classes string
}

func Css(d string) CssAttribute {
	return CssAttribute{classes: d}
}

func CssIf(v bool, d string) CssAttribute {
	if v {
		return CssAttribute{classes: d}
	}
	return CssAttribute{}
}

type OnClickAttribute struct {
	cb func()
}

func OnClick(cb func()) OnClickAttribute {
	return OnClickAttribute{cb: cb}
}

type OnChangeAttribute struct {
	cb js.EventHandlerFunc
}

func OnChange(cb js.EventHandlerFunc) OnChangeAttribute {
	return OnChangeAttribute{cb: cb}
}

type OnInputAttribute struct {
	cb js.EventHandlerFunc
}

func OnInput(cb js.EventHandlerFunc) OnInputAttribute {
	return OnInputAttribute{cb: cb}
}

type HelmetTitle string
type HelmetDescription string
type HelmetAuthor string
type HelmetKeywords string

func mergeAttributes(parent *Element, uis ...interface{}) {
	elems := []UI{}
	for _, v := range uis {
		switch c := v.(type) {
		case Attribute:
			parent.setAttr(c.Key, c.Value)
		case CssAttribute:
			if vv, ok := parent.attrs["classes"]; ok {
				parent.setAttr("class", vv+" "+c.classes)
			} else {
				parent.setAttr("class", c.classes)
			}
		case OnClickAttribute:
			parent.setEventHandler("click", func(e js.Event) {
				c.cb()
			})
		case OnChangeAttribute:
			parent.setEventHandler("change", c.cb)
		case OnInputAttribute:
			parent.setEventHandler("input", c.cb)
		case HelmetTitle:
			helmet.Title = string(c)
		case HelmetDescription:
			helmet.Description = string(c)
		case HelmetAuthor:
			helmet.Author = string(c)
		case HelmetKeywords:
			helmet.Keywords = string(c)
		case UI:
			elems = append(elems, c)
		case nil:
			// dont need to add nil items
		default:
			panic("unknown type in render: " + reflect.TypeOf(v).String())
		}
	}
	parent.setBody(elems)
}
