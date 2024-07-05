package main

import (
	"context"
	"encoding/json"
	"fmt"

	"honnef.co/go/js/dom/v2"
	fetch "marwan.io/wasm-fetch"
)

func main() {
	c := make(chan struct{})

	root := dom.GetWindow().Document().QuerySelector("body")
	root.SetAttribute("id", "root")

	ctx := context.Background()

	count := NewSignal(0)

	data := AsyncData(ctx, func() (string, error) {
		resp, err := fetch.Fetch("https://jsonplaceholder.typicode.com/todos/1", &fetch.Opts{})
		if err != nil {
			return "", err
		}

		d := make(map[string]interface{})
		err = json.Unmarshal(resp.Body, &d)
		if err != nil {
			return "", err
		}

		return d["title"].(string), nil
	})

	inputValue := NewSignal("")

	d := Div(
		Button(
			func(e dom.Event) {
				count.Set(ctx, count.Get(ctx)-1)
			},
			Text("-"),
		),
		TextR(func(ctx context.Context) string { return fmt.Sprintf("Count: %d", count.Get(ctx)) }),
		Button(
			func(e dom.Event) {
				count.Set(ctx, count.Get(ctx)+1)
			},
			Text("+"),
		),
		Div(
			Input(
				Event("input", func(e dom.Event) {
					inputValue.Set(context.Background(), e.Target().(*dom.HTMLInputElement).Value())
				}),
			),
		),
		TextR(inputValue.Get),
		Div(
			Switch(data.State.Get).
				When(AsyncStateLoading, func() GueElement { return Text("Loading...") }).
				When(AsyncStateError, func() GueElement { return Text("Oops error!") }).
				When(AsyncStateIdle, func() GueElement {
					return TextR(data.Value.Get)
				}),
		),
	)

	d.Render(root)

	<-c
}
