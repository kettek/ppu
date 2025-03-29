package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var appTabs *container.AppTabs

func main() {
	a := app.New()
	w := a.NewWindow("PPU")

	var listedEntries []*Entry
	var listEntry *Entry

	var search *widget.Entry
	var results *widget.List
	var toolbar *fyne.Container

	selectEntry := func(entry *Entry) {
		listEntry = entry
		if entry == nil {
			toolbar.Objects[1].(*widget.Button).Disable()
			toolbar.Objects[2].(*widget.Button).Disable()
		} else {
			toolbar.Objects[1].(*widget.Button).Enable()
			toolbar.Objects[2].(*widget.Button).Enable()
		}
	}

	refreshResults := func() {
		listedEntries = filterEntriesByTags(entries, stringToTags(search.Text))
		results.UnselectAll()
		results.Refresh()
		selectEntry(nil)
	}

	search = widget.NewEntry()
	search.OnChanged = func(s string) {
		refreshResults()
	}
	results = widget.NewList(
		func() int {
			return len(listedEntries)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewVBoxLayout(),
				container.New(layout.NewGridLayout(4),
					widget.NewLabel("name"),
					widget.NewLabel("unit"),
					widget.NewLabel("cost"),
					widget.NewLabel("ppu"),
				),
				widget.NewLabel("tags"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			entry := listedEntries[i]
			vbox := o.(*fyne.Container).Objects
			items := vbox[0].(*fyne.Container).Objects
			tags := vbox[1].(*widget.Label)
			items[0].(*widget.Label).SetText(entry.Name)

			ppu := math.Round(entry.Cost/entry.Units*100) / 100

			items[1].(*widget.Label).SetText(fmt.Sprintf("%g", entry.Units))
			items[2].(*widget.Label).SetText(fmt.Sprintf("%g", entry.Cost))
			items[3].(*widget.Label).SetText(fmt.Sprintf("%.2f", ppu))

			tags.SetText(strings.Join(entry.Tags, " "))
		},
	)
	results.OnSelected = func(id widget.ListItemID) {
		selectEntry(listedEntries[id])
	}

	toolbar = container.NewHBox(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			var popup *widget.PopUp
			name := widget.NewEntry()
			units := widget.NewEntry()
			cost := widget.NewEntry()
			format := widget.NewSelectEntry([]string{
				"kg",
				"count",
			})
			tags := widget.NewSelectEntry(getAllTags())
			form := &widget.Form{
				Items: []*widget.FormItem{
					{Text: "Name", Widget: name},
					{Text: "Tags", Widget: tags},
					{Text: "Units", Widget: units},
					{Text: "Format", Widget: format},
					{Text: "Cost", Widget: cost},
				},
				OnSubmit: func() {
					units, _ := strconv.ParseFloat(units.Text, 64)
					cost, _ := strconv.ParseFloat(cost.Text, 64)

					entries = append(entries, &Entry{
						Name:   name.Text,
						Tags:   stringToTags(tags.Text),
						Cost:   cost,
						Units:  units,
						Format: UnitFormat(format.SelectedText()),
					})
					refreshResults()
					popup.Hide()
				},
				OnCancel: func() {
					popup.Hide()
				},
			}
			popup = widget.NewModalPopUp(form, w.Canvas())
			popup.Show()
		}),
		widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			var popup *widget.PopUp
			name := widget.NewEntry()
			name.SetText(listEntry.Name)
			units := widget.NewEntry()
			units.SetText(fmt.Sprintf("%g", listEntry.Units))
			cost := widget.NewEntry()
			cost.SetText(fmt.Sprintf("%g", listEntry.Cost))
			format := widget.NewSelectEntry([]string{
				"kg",
				"count",
			})
			format.SetText(string(listEntry.Format))
			tags := widget.NewSelectEntry(getAllTags())
			tags.SetText(strings.Join(listEntry.Tags, " "))
			form := &widget.Form{
				Items: []*widget.FormItem{
					{Text: "Name", Widget: name},
					{Text: "Tags", Widget: tags},
					{Text: "Units", Widget: units},
					{Text: "Format", Widget: format},
					{Text: "Cost", Widget: cost},
				},
				OnSubmit: func() {
					units, _ := strconv.ParseFloat(units.Text, 64)
					cost, _ := strconv.ParseFloat(cost.Text, 64)

					listEntry.Name = name.Text
					listEntry.Tags = stringToTags(tags.Text)
					listEntry.Cost = cost
					listEntry.Units = units
					listEntry.Format = UnitFormat(format.SelectedText())
					refreshResults()
					popup.Hide()
				},
				OnCancel: func() {
					popup.Hide()
				},
			}
			popup = widget.NewModalPopUp(form, w.Canvas())
			popup.Show()
		}),
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			for i, ent := range entries {
				if ent == listEntry {
					entries = append(entries[:i], entries[i+1:]...)
					refreshResults()
					return
				}
			}
		}),
	)

	selectEntry(nil) // Just to ensure UI is sync'd

	container := container.NewBorder(search, toolbar, nil, nil, results)

	w.SetContent(container)

	w.ShowAndRun()
}
