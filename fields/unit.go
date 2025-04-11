package fields

import (
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Unit struct {
}

func (unit *Unit) Name() string {
	return "unit"
}

func (p *Unit) Label() string {
	return "Unit"
}

func (p *Unit) Value(values map[string]any) string {
	if val, ok := values["unit"].(string); ok {
		return val
	}
	return ""
}

func (p *Unit) FormItem(values map[string]any) *widget.FormItem {
	e := widget.NewEntry()
	e.SetText(p.Value(values))
	return &widget.FormItem{
		Text:   "Unit",
		Widget: e,
	}
}

func (p *Unit) ModifyEntry(values *map[string]any, e fyne.CanvasObject) {
	if *values == nil {
		*values = make(map[string]any)
	}
	(*values)[p.Name()] = e.(*widget.Entry).Text
}

func (p *Unit) Sort(entries []map[string]any) []int {
	var unsorted []int
	for i := range len(entries) {
		unsorted = append(unsorted, i)
	}
	slices.SortFunc(unsorted, func(i, j int) int {
		var anum float64
		var bnum float64
		var err error
		if a, ok := entries[i][p.Name()].(string); ok {
			anum, err = strconv.ParseFloat(a, 64)
			if err != nil {
				return 0
			}
		}
		if b, ok := entries[j][p.Name()].(string); ok {
			bnum, err = strconv.ParseFloat(b, 64)
			if err != nil {
				return 0
			}
		}
		if anum < bnum {
			return -1
		}
		if anum > bnum {
			return 1
		}
		return 0
	})
	return unsorted
}

func init() {
	addField(&Unit{})
}
