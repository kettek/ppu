package main

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Field interface {
	Name() string
	Label() string
	Value(entry *Entry) string
	FormItem(*Entry) *widget.FormItem
	ModifyEntry(entry *Entry, e fyne.CanvasObject)
}

type ProteinField struct {
}

func (protein *ProteinField) Name() string {
	return "protein"
}

func (p *ProteinField) Label() string {
	return "Protein"
}

func (p *ProteinField) Value(entry *Entry) string {
	if val, ok := entry.Values["protein"].(string); ok {
		return val
	}
	return ""
}

func (p *ProteinField) FormItem(entry *Entry) *widget.FormItem {
	e := widget.NewEntry()
	e.SetText(p.Value(entry))
	return &widget.FormItem{
		Text:   "Protein",
		Widget: e,
	}
}

func (p *ProteinField) ModifyEntry(entry *Entry, e fyne.CanvasObject) {
	if entry.Values == nil {
		entry.Values = make(map[string]any)
	}
	entry.Values[p.Name()] = e.(*widget.Entry).Text
}

var fieldMap = map[string]Field{}
var fieldSlice []string

func getFieldFormItems(entry *Entry) []*widget.FormItem {
	var items []*widget.FormItem
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			items = append(items, field.FormItem(entry))
		}
	}
	return items
}

func getFieldFormModifiers() []func(entry *Entry, e fyne.CanvasObject) {
	var items []func(entry *Entry, e fyne.CanvasObject)
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			items = append(items, field.ModifyEntry)
		}
	}
	return items
}

func addField(name string, field Field) {
	fieldMap[name] = field
	fieldSlice = append(fieldSlice, name)
	slices.SortStableFunc(fieldSlice, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

func init() {
	pf := &ProteinField{}
	addField(pf.Name(), pf)
}
