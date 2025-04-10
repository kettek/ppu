package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/kettek/ppu/fields"
)

var appTabs *container.AppTabs

func main() {
	a := app.NewWithID("net.kettek.ppu")
	w := a.NewWindow("PPU")

	// Try to load in entries.
	jsonData := a.Preferences().StringWithFallback("entries", "[]")
	if err := json.Unmarshal([]byte(jsonData), &entries); err != nil {
		fmt.Println("Error loading entries:", err)
		entries = []*Entry{}
	}

	writeEntries := func() {
		data, err := json.Marshal(entries)
		if err != nil {
			fmt.Println("Error writing entries:", err)
			return
		}
		a.Preferences().SetString("entries", string(data))
	}

	var listedEntries []*Entry
	var listEntry *Entry

	var search *widget.Entry
	var headers *fyne.Container
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
		listedEntries = sortEntriesByPPU(filterEntriesByTags(entries, stringToTags(search.Text)))
		results.UnselectAll()
		results.Refresh()
		selectEntry(nil)
	}

	search = widget.NewEntry()
	search.OnChanged = func(s string) {
		refreshResults()
	}

	headerLabels := []fyne.CanvasObject{
		widget.NewLabel("name"),
		widget.NewLabel("unit"),
		widget.NewLabel("cost"),
		widget.NewLabel("ppu"),
	}
	for _, label := range fields.GetLabels() {
		headerLabels = append(headerLabels, widget.NewLabel(label))
	}
	headers = container.New(layout.NewGridLayout(len(headerLabels)), headerLabels...)

	results = widget.NewList(
		func() int {
			return len(listedEntries)
		},
		func() fyne.CanvasObject {

			labels := []fyne.CanvasObject{
				widget.NewLabel("name"),
				widget.NewLabel("unit"),
				widget.NewLabel("cost"),
				widget.NewLabel("ppu"),
			}

			for _, label := range fields.GetLabels() {
				labels = append(labels, widget.NewLabel(label))
			}

			return container.New(layout.NewVBoxLayout(),
				container.New(layout.NewGridLayout(len(labels)),
					labels...,
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

			for i, field := range fields.GetFields() {
				items[i+4].(*widget.Label).SetText(field.Value(entry.Values))
			}

			tags.SetText(strings.Join(entry.Tags, ", "))
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
			format.SetText("count")
			tags := widget.NewSelectEntry(getAllTags())
			t := time.Now()
			var calendarPopup *widget.PopUp
			var dateButton *widget.Button
			calendar := xwidget.NewCalendar(t, func(t2 time.Time) {
				t = t2
				dateButton.SetText(t.Local().Format("2006-01-02"))
			})
			dateButton = widget.NewButton(t.Local().Format("2006-01-02"), func() {
				if calendarPopup == nil {
					calendarPopup = widget.NewPopUp(calendar, w.Canvas())
				}
				calendarPopup.ShowAtRelativePosition(fyne.NewPos(0, 0), dateButton)
			})

			formItems := []*widget.FormItem{
				{Text: "Name", Widget: name},
				{Text: "Tags", Widget: tags},
				{Text: "Units", Widget: units},
				{Text: "Format", Widget: format},
				{Text: "Cost", Widget: cost},
				{Text: "Date", Widget: dateButton},
			}
			entry := &Entry{}
			formFieldItems := fields.GetFormItems(entry.Values)
			formFieldModifiers := fields.GetFormModifiers()
			formItems = append(formItems, formFieldItems...)

			form := &widget.Form{
				Items: formItems,
				OnSubmit: func() {
					units, _ := strconv.ParseFloat(units.Text, 64)
					cost, _ := strconv.ParseFloat(cost.Text, 64)

					entry.Name = name.Text
					entry.Tags = stringToTags(tags.Text)
					entry.Cost = cost
					entry.Units = units
					entry.Format = UnitFormat(format.SelectedText())
					entry.Date = t

					for i, modifier := range formFieldModifiers {
						if i < len(formFieldItems) {
							modifier(&entry.Values, formFieldItems[i].Widget)
						}
					}

					entries = append(entries, entry)
					writeEntries()
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
			if listEntry == nil {
				return
			}
			entries = append(entries, &Entry{
				Name:   listEntry.Name,
				Tags:   listEntry.Tags,
				Cost:   listEntry.Cost,
				Units:  listEntry.Units,
				Format: listEntry.Format,
				Values: listEntry.Values,
			})
			writeEntries()
			refreshResults()
		}),
		widget.NewButtonWithIcon("", theme.StorageIcon(), func() {
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
			tags.SetText(strings.Join(listEntry.Tags, ", "))
			t := listEntry.Date
			if t.IsZero() {
				t = time.Now()
			}
			var calendarPopup *widget.PopUp
			var dateButton *widget.Button
			calendar := xwidget.NewCalendar(t, func(t2 time.Time) {
				t = t2
				dateButton.SetText(t.Local().Format("2006-01-02"))
			})
			dateButton = widget.NewButton(t.Local().Format("2006-01-02"), func() {
				if calendarPopup == nil {
					calendarPopup = widget.NewPopUp(calendar, w.Canvas())
				}
				calendarPopup.ShowAtRelativePosition(fyne.NewPos(0, 0), dateButton)
			})

			formItems := []*widget.FormItem{
				{Text: "Name", Widget: name},
				{Text: "Tags", Widget: tags},
				{Text: "Units", Widget: units},
				{Text: "Format", Widget: format},
				{Text: "Cost", Widget: cost},
				{Text: "Date", Widget: dateButton},
			}
			formFieldItems := fields.GetFormItems(listEntry.Values)
			formFieldModifiers := fields.GetFormModifiers()
			formItems = append(formItems, formFieldItems...)

			form := &widget.Form{
				Items: formItems,
				OnSubmit: func() {
					units, _ := strconv.ParseFloat(units.Text, 64)
					cost, _ := strconv.ParseFloat(cost.Text, 64)

					listEntry.Name = name.Text
					listEntry.Tags = stringToTags(tags.Text)
					listEntry.Cost = cost
					listEntry.Units = units
					listEntry.Format = UnitFormat(format.SelectedText())
					listEntry.Date = t

					for i, modifier := range formFieldModifiers {
						if i < len(formFieldItems) {
							modifier(&listEntry.Values, formFieldItems[i].Widget)
						}
					}

					writeEntries()
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
					writeEntries()
					refreshResults()
					return
				}
			}
		}),
	)

	refreshResults()

	container := container.NewBorder(container.NewVBox(search, headers), toolbar, nil, nil, results)

	w.SetContent(container)

	w.ShowAndRun()
}
