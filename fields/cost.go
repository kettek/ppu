package fields

import (
	"fmt"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Cost struct {
}

func (cost *Cost) Name() string {
	return "cost"
}

func (p *Cost) Label() string {
	return "Cost"
}

func (p *Cost) Value(values map[string]any) string {
	if val, ok := values["cost"].(string); ok {
		return val
	}
	return ""
}

func (p *Cost) FormItem(values map[string]any) *widget.FormItem {
	e := widget.NewEntry()
	e.SetText(p.Value(values))
	return &widget.FormItem{
		Text:   "Cost",
		Widget: e,
	}
}

func (p *Cost) ModifyEntry(values *map[string]any, e fyne.CanvasObject) {
	if *values == nil {
		*values = make(map[string]any)
	}
	(*values)[p.Name()] = e.(*widget.Entry).Text
}

func (p *Cost) Sort(entries []map[string]any) []int {
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

func (p *Cost) Upgrade(values map[string]any) (string, string) {
	if val, ok := values["cost"].(float64); ok {
		return "cost", fmt.Sprintf("%g", val)
	} else if val, ok := values["Cost"].(float64); ok {
		return "cost", fmt.Sprintf("%g", val)
	}
	return "", ""
}

func init() {
	addField(&Cost{}, 5)
}
