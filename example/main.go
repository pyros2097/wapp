package main

import (
	"strconv"

	. "github.com/pyros2097/wapp"
)

func Index(c *RenderContext) UI {
	count, setCount := c.UseInt(0)
	inc := func() { setCount(count() + 1) }
	dec := func() { setCount(count() - 1) }
	return Col(
		Row(
			Row(Css("text-6xl"),
				Text("Counter"),
			),
		),
		Row(
			Row(Css("text-6xl m-20 cursor-pointer select-none"), OnClick(dec),
				Text("-"),
			),
			Row(Css("text-6xl m-20"),
				Text(strconv.Itoa(count())),
			),
			Row(Css("text-6xl m-20 cursor-pointer select-none"), OnClick(inc),
				Text("+"),
			),
		),
	)
}

func main() {
	Route("/", Index)
	Run()
}
