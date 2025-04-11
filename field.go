package main

type FieldSortable interface {
	Sort(entries []map[string]any) []int
}
