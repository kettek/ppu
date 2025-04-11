package fields

import (
	"fmt"
	"slices"
	"strconv"
)

type PPU struct {
}

func (ppu *PPU) Name() string {
	return "ppu"
}

func (p *PPU) Label() string {
	return "PPU"
}

func (p *PPU) Value(values map[string]any) string {
	return fmt.Sprintf("%g", p.value(values))
}

func (p *PPU) value(values map[string]any) float64 {
	var err error
	cost := 0.0
	units := 0.0

	if val, ok := values["cost"].(string); ok {
		cost, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
	}
	if val, ok := values["unit"].(string); ok {
		units, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
	}

	if units == 0 {
		return 0
	}

	return cost / units
}

func (p *PPU) Sort(entries []map[string]any) []int {
	var unsorted []int
	for i := range len(entries) {
		unsorted = append(unsorted, i)
	}
	slices.SortFunc(unsorted, func(i, j int) int {
		anum := p.value(entries[i])
		bnum := p.value(entries[j])
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
	addField(&PPU{})
}
