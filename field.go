package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type FieldSortable interface {
	Sort(entries []map[string]any) []int
}

type FieldUpgradable interface {
	Upgrade(fields map[string]any) (string, string)
}

type FieldVirtual interface {
	Name() string
	Label() string
	Value(map[string]any) string
}

type Field interface {
	FieldVirtual
	FormItem(map[string]any) *widget.FormItem
	ModifyEntry(*map[string]any, fyne.CanvasObject)
}
