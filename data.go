package main

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

var entries []*Entry

type UnitFormat string

var (
	UnitFormatKG    = "kg"
	UnitFormatCount = "count"
)

// Entry would be a... something
type Entry struct {
	Name   string  // Store name, etc.
	Cost   float64 // Expressed in... whatever?
	Units  float64 // kg?
	Format UnitFormat
	Tags   []string
	Date   time.Time
	Values map[string]any
}

func getAllTags() []string {
	var tags []string
	for _, entry := range entries {
		for _, t := range entry.Tags {
			found := false
			for _, t2 := range tags {
				if t == t2 {
					found = true
					break
				}
			}
			if !found {
				tags = append(tags, t)
			}
		}
	}
	return tags
}

var months = []string{
	"january", "february", "march", "april", "may", "june",
	"july", "august", "september", "october", "november", "december",
}

var shortMonths = []string{
	"jan", "feb", "mar", "apr", "may", "jun",
	"jul", "aug", "sep", "oct", "nov", "dec",
}

func matchExtended(e *Entry, ext string) bool {
	now := time.Now()
	if ext == "recent" {
		if e.Date.After(now.Add(-time.Hour * 24 * 7)) {
			return true
		}
	} else if slices.Contains(months, ext) {
		if strings.ToLower(e.Date.Month().String()) == ext {
			return true
		}
	} else if slices.Contains(shortMonths, ext) {
		short := strings.ToLower(e.Date.Month().String()[0:3])
		if short == ext {
			return true
		}
	} else {
		number, err := strconv.Atoi(ext)
		if err == nil {
			if e.Date.Year() == number {
				return true
			}
			if e.Date.Month() == time.Month(number) {
				return true
			}
		}
	}
	return false
}

func filterEntriesByTags(entries []*Entry, tags []string) []*Entry {
	var entries2 []*Entry

	for _, e := range entries {
		matches := true
		for _, t := range tags {
			if len(t) == 0 {
				continue
			}
			t = strings.ToLower(t)

			recent := matchExtended(e, t)
			if recent {
				continue
			}

			missing := true
			tags := append(e.Tags, strings.ToLower(e.Name))
			for _, t2 := range tags {
				t2 = strings.ToLower(t2)
				if strings.Contains(t2, t) {
					missing = false
					break
				}
			}
			if missing {
				matches = false
				break
			}
		}
		if matches {
			entries2 = append(entries2, e)
		}
	}

	return entries2
}

func sortEntriesByPPU(entries []*Entry) []*Entry {
	// Sort by price per unit
	for i := 0; i < len(entries)-1; i++ {
		for j := 0; j < len(entries)-i-1; j++ {
			if entries[j].Cost/entries[j].Units > entries[j+1].Cost/entries[j+1].Units {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}
	return entries
}

func stringToTags(s string) []string {
	split := strings.Split(s, ",")
	for i, t := range split {
		split[i] = strings.TrimSpace(t)
	}
	return split
}
