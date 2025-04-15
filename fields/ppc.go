package fields

import (
	"fmt"
	"slices"
	"strconv"
)

type PPC struct {
}

func (ppu *PPC) Name() string {
	return "ppc"
}

func (p *PPC) Label() string {
	return "PPC"
}

func (p *PPC) Value(values map[string]any) string {
	return fmt.Sprintf("%.2f", p.value(values))
}

func (p *PPC) value(values map[string]any) float64 {
	var err error
	cost := 0.0
	units := 0.0
	protein := 0.0

	if val, ok := values["cost"].(string); ok {
		cost, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
	}
	if val, ok := values["protein"].(string); ok {
		protein, err = strconv.ParseFloat(val, 64)
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

	if cost == 0 {
		return 0
	}

	return (protein * units) / cost
}

func (p *PPC) Sort(entries []map[string]any) []int {
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
	addField(&PPC{}, 7)
}
