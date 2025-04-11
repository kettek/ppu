package fields

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Field interface {
	Name() string
	Label() string
	Value(map[string]any) string
	FormItem(map[string]any) *widget.FormItem
	ModifyEntry(*map[string]any, fyne.CanvasObject)
}

var fieldMap = map[string]Field{}
var fieldSlice []string

func GetFormItems(values map[string]any) []*widget.FormItem {
	var items []*widget.FormItem
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			items = append(items, field.FormItem(values))
		}
	}
	return items
}

func GetFormModifiers() []func(values *map[string]any, e fyne.CanvasObject) {
	var items []func(values *map[string]any, e fyne.CanvasObject)
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			items = append(items, field.ModifyEntry)
		}
	}
	return items
}

func GetNames() []string {
	return fieldSlice
}

func GetLabels() []string {
	var labels []string
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			labels = append(labels, field.Label())
		}
	}
	return labels
}

func GetFields() []Field {
	var fields []Field
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			fields = append(fields, field)
		}
	}
	return fields
}

func addField(field Field) {
	fieldMap[field.Name()] = field
	fieldSlice = append(fieldSlice, field.Name())
	slices.SortStableFunc(fieldSlice, func(a, b string) int {
		return strings.Compare(a, b)
	})
}
