package fields

import (
	"slices"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Field interface {
	Name() string
	Label() string
	Value(map[string]any) string
}

type FieldEditable interface {
	Field
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
			if fe, ok := field.(FieldEditable); ok {
				items = append(items, fe.FormItem(values))
			}
		}
	}
	return items
}

func GetFormModifiers() []func(values *map[string]any, e fyne.CanvasObject) {
	var items []func(values *map[string]any, e fyne.CanvasObject)
	for _, name := range fieldSlice {
		field := fieldMap[name]
		if field != nil {
			if fe, ok := field.(FieldEditable); ok {
				items = append(items, fe.ModifyEntry)
			}
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

func GetFieldAsFloat(name string, values map[string]any) float64 {
	if field, ok := fieldMap[name]; ok {
		if val, ok := values[field.Name()].(string); ok {
			if num, err := strconv.ParseFloat(val, 64); err == nil {
				return num
			}
		}
	}
	return 0.0
}

func addField(field Field) {
	fieldMap[field.Name()] = field
	fieldSlice = append(fieldSlice, field.Name())
	slices.SortStableFunc(fieldSlice, func(a, b string) int {
		return strings.Compare(a, b)
	})
}
