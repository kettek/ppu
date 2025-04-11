package fields

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Name struct {
}

func (name *Name) Name() string {
	return "name"
}

func (p *Name) Label() string {
	return "Name"
}

func (p *Name) Value(values map[string]any) string {
	if val, ok := values["name"].(string); ok {
		return val
	}
	return ""
}

func (p *Name) FormItem(values map[string]any) *widget.FormItem {
	e := widget.NewEntry()
	e.SetText(p.Value(values))
	return &widget.FormItem{
		Text:   "Name",
		Widget: e,
	}
}

func (p *Name) ModifyEntry(values *map[string]any, e fyne.CanvasObject) {
	if *values == nil {
		*values = make(map[string]any)
	}
	(*values)[p.Name()] = e.(*widget.Entry).Text
}

func (p *Name) Sort(entries []map[string]any) []int {
	var unsorted []int
	for i := range len(entries) {
		unsorted = append(unsorted, i)
	}
	slices.SortFunc(unsorted, func(i, j int) int {
		return strings.Compare(entries[i][p.Name()].(string), entries[j][p.Name()].(string))
	})
	return unsorted
}

func (p *Name) Upgrade(values map[string]any) (string, string) {
	if val, ok := values["name"].(string); ok {
		return "name", val
	} else if val, ok := values["Name"].(string); ok {
		return "name", val
	}
	return "", ""
}

func init() {
	addField(&Name{}, 0)
}
