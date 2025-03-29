package main

import (
	"strings"
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

func filterEntriesByTags(entries []*Entry, tags []string) []*Entry {
	var entries2 []*Entry

	for _, e := range entries {
		matches := true
		for _, t := range tags {
			if len(t) == 0 {
				continue
			}
			t = strings.ToLower(t)
			missing := true
			tags := append(e.Tags, strings.ToLower(e.Name))
			for _, t2 := range tags {
				t2 = strings.ToLower(t2)
				if strings.Contains(t2, t) {
					//if t == t2 {
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

func stringToTags(s string) []string {
	return strings.Split(s, " ")
}
