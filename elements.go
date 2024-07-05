package main

import (
	"context"

	"honnef.co/go/js/dom/v2"
)

type InputElement struct {
	BasicElement
}

func Input(children ...GueElement) *InputElement {
	return &InputElement{
		BasicElement: newBasicElement("input", children),
	}
}

type TextElement struct {
	BasicElement
	Factory func(ctx context.Context) string
}

func (t *TextElement) Render(parent dom.Element) {
	e := dom.GetWindow().Document().CreateTextNode("")
	ctx := context.Background()
	unsup := Effect(ctx, func(ctx context.Context) {
		e.SetTextContent(t.Factory(ctx))
	})

	t.RemoveFunc = func() {
		unsup()
		e.SetTextContent("")
	}

	parent.AppendChild(e)
}

func TextR(f func(ctx context.Context) string) *TextElement {
	return &TextElement{
		Factory: f,
	}
}

func Text(value string) *TextElement {
	return &TextElement{
		Factory: func(ctx context.Context) string {
			return value
		},
	}
}

type ButtonElement struct {
	BasicElement
}

func Button(handler func(dom.Event), children ...GueElement) *ButtonElement {
	return &ButtonElement{
		BasicElement: newBasicElement("button", append(children, Event("click", handler))),
	}
}

type AsyncState int

const (
	AsyncStateLoading AsyncState = iota
	AsyncStateError
	AsyncStateIdle
)

type AsyncDataResult[T any] struct {
	Value *Signal[T]
	State *Signal[AsyncState]
}

func AsyncData[T any](ctx context.Context, fetcher func() (T, error)) AsyncDataResult[T] {
	var zeroT T
	data := AsyncDataResult[T]{
		Value: NewSignal[T](zeroT),
		State: NewSignal[AsyncState](AsyncStateLoading),
	}
	go func() {
		res, err := fetcher()
		if err != nil {
			data.State.Set(ctx, AsyncStateError)
			return
		}
		Batch(ctx, func(ctx context.Context) {
			data.Value.Set(ctx, res)
			data.State.Set(ctx, AsyncStateIdle)
		})
	}()

	return data
}

type switchCase[T comparable] struct {
	value   T
	handler func() GueElement
}

type SwitchElement[T comparable] struct {
	BasicElement
	factory func(ctx context.Context) T
	cases   []switchCase[T]
}

func (t *SwitchElement[T]) When(c T, handler func() GueElement) *SwitchElement[T] {
	t.cases = append(t.cases, switchCase[T]{c, handler})
	return t
}

func (t *SwitchElement[T]) Render(parent dom.Element) {
	Effect(context.Background(), func(ctx context.Context) {
		t.Remove(parent)
		value := t.factory(ctx)
		for _, sc := range t.cases {
			if sc.value == value {
				e := sc.handler()
				e.Render(parent)
				t.RemoveFunc = func() {
					e.Remove(parent)
				}
				return
			}
		}
	})
}

func Switch[T comparable](value func(ctx context.Context) T) *SwitchElement[T] {
	return &SwitchElement[T]{
		factory: value,
	}
}
