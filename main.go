package main

import (
	"encoding/json"
	"fmt"
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

	// Let's also do an "upgrade" step.
	var entryFields []map[string]any
	if err := json.Unmarshal([]byte(jsonData), &entryFields); err == nil {
		for i, entry := range entryFields {
			if entries[i].Values == nil {
				entries[i].Values = make(map[string]any)
			}
			for _, field := range fields.GetFields() {
				if fu, ok := field.(FieldUpgradable); ok {
					if k, v := fu.Upgrade(entry); k != "" {
						entries[i].Values[k] = v
					}
				}
			}
		}
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

	var sortName string
	var sortFunc func([]map[string]any) []int
	var reversed bool
	var sort = func() {
		if sortFunc == nil {
			return
		}
		fieldMaps := make([]map[string]any, len(listedEntries))
		for i, entry := range listedEntries {
			fieldMaps[i] = entry.Values
		}
		sortedEntries := make([]*Entry, len(listedEntries))
		for i, index := range sortFunc(fieldMaps) {
			sortedEntries[i] = listedEntries[index]
		}
		listedEntries = sortedEntries
		if reversed {
			for i, j := 0, len(listedEntries)-1; i < j; i, j = i+1, j-1 {
				listedEntries[i], listedEntries[j] = listedEntries[j], listedEntries[i]
			}
		}
	}

	refreshResults := func() {
		listedEntries = filterEntriesByTags(entries, stringToTags(search.Text))
		results.UnselectAll()
		sort()
		results.Refresh()
		selectEntry(nil)
	}

	search = widget.NewEntry()
	search.OnChanged = func(s string) {
		refreshResults()
	}

	headerLabels := []fyne.CanvasObject{}
	for _, field := range fields.GetFields() {
		label := field.Label()
		if fs, ok := field.(FieldSortable); ok {
			button := widget.NewButton(label, func() {
				if sortName == field.Name() {
					reversed = !reversed
				}
				sortFunc = fs.Sort
				sortName = field.Name()
				sort()
				results.Refresh()
			})
			headerLabels = append(headerLabels, button)
		} else {
			headerLabels = append(headerLabels, widget.NewLabel(label))
		}
	}
	headers = container.New(layout.NewGridLayout(len(headerLabels)), headerLabels...)

	results = widget.NewList(
		func() int {
			return len(listedEntries)
		},
		func() fyne.CanvasObject {

			labels := []fyne.CanvasObject{}

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

			for i, field := range fields.GetFields() {
				items[i].(*widget.Label).SetText(field.Value(entry.Values))
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
				{Text: "Tags", Widget: tags},
				{Text: "Format", Widget: format},
				{Text: "Date", Widget: dateButton},
			}
			entry := &Entry{}
			formFieldItems := fields.GetFormItems(entry.Values)
			formFieldModifiers := fields.GetFormModifiers()
			formItems = append(formItems, formFieldItems...)

			form := &widget.Form{
				Items: formItems,
				OnSubmit: func() {
					entry.Tags = stringToTags(tags.Text)
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
			entry := &Entry{
				Values: make(map[string]any),
			}
			for k, v := range listEntry.Values {
				entry.Values[k] = v
			}
			entry.Tags = make([]string, len(listEntry.Tags))
			copy(entry.Tags, listEntry.Tags)
			entry.Format = listEntry.Format
			entry.Date = time.Now()
			entries = append(entries, entry)
			writeEntries()
			refreshResults()
		}),
		widget.NewButtonWithIcon("", theme.StorageIcon(), func() {
			var popup *widget.PopUp
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
				{Text: "Tags", Widget: tags},
				{Text: "Format", Widget: format},
				{Text: "Date", Widget: dateButton},
			}
			formFieldItems := fields.GetFormItems(listEntry.Values)
			formFieldModifiers := fields.GetFormModifiers()
			formItems = append(formItems, formFieldItems...)

			form := &widget.Form{
				Items: formItems,
				OnSubmit: func() {
					listEntry.Tags = stringToTags(tags.Text)
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
