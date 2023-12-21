//go:build js
// +build js

package main

import (
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

type GueElement interface {
	Render(parent dom.Element)
	Remove(parent dom.Element)
}

type EventElement struct {
	Name    string
	Handler func(dom.Event)
	jsFunc  js.Func
}

func (t *EventElement) Render(parent dom.Element) {
	t.jsFunc = parent.AddEventListener(t.Name, true, t.Handler)
}

func (t *EventElement) Remove(parent dom.Element) {
	parent.RemoveEventListener(t.Name, true, t.jsFunc)
}

func Event(name string, handler func(dom.Event)) *EventElement {
	return &EventElement{
		Name:    name,
		Handler: handler,
	}
}

type BasicElement struct {
	Tag        string
	Children   []GueElement
	RemoveFunc func()
}

func newBasicElement(tag string, children []GueElement) BasicElement {
	return BasicElement{
		Tag:      tag,
		Children: children,
	}
}

func (t *BasicElement) Remove(parent dom.Element) {
	if t.RemoveFunc != nil {
		t.RemoveFunc()
	}
}

func (t *BasicElement) Render(parent dom.Element) {
	e := dom.GetWindow().Document().CreateElement(t.Tag)

	for _, ch := range t.Children {
		ch.Render(e)
	}

	t.RemoveFunc = func() {
		for _, ch := range t.Children {
			ch.Remove(e)
		}
		e.Remove()
	}

	parent.AppendChild(e)
}

type DivElement struct {
	BasicElement
}

func Div(children ...GueElement) *DivElement {
	return &DivElement{BasicElement: newBasicElement("div", children)}
}
