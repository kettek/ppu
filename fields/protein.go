package fields

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Protein struct {
}

func (protein *Protein) Name() string {
	return "protein"
}

func (p *Protein) Label() string {
	return "Protein"
}

func (p *Protein) Value(values map[string]any) string {
	if val, ok := values["protein"].(string); ok {
		return val
	}
	return ""
}

func (p *Protein) FormItem(values map[string]any) *widget.FormItem {
	e := widget.NewEntry()
	e.SetText(p.Value(values))
	return &widget.FormItem{
		Text:   "Protein",
		Widget: e,
	}
}

func (p *Protein) ModifyEntry(values *map[string]any, e fyne.CanvasObject) {
	if *values == nil {
		*values = make(map[string]any)
	}
	(*values)[p.Name()] = e.(*widget.Entry).Text
}

func init() {
	addField(&Protein{})
}
