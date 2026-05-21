package comb

import (
	"reflect"
	"testing"
)

func TestCombinations(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		k     int
		want  [][]string
	}{
		{
			name:  "3 choose 2",
			items: []string{"A", "B", "C"},
			k:     2,
			want: [][]string{
				{"A", "B"},
				{"A", "C"},
				{"B", "C"},
			},
		},
		{
			name:  "4 choose 2",
			items: []string{"A", "B", "C", "D"},
			k:     2,
			want: [][]string{
				{"A", "B"},
				{"A", "C"},
				{"A", "D"},
				{"B", "C"},
				{"B", "D"},
				{"C", "D"},
			},
		},
		{
			name:  "3 choose 1",
			items: []string{"A", "B", "C"},
			k:     1,
			want: [][]string{
				{"A"},
				{"B"},
				{"C"},
			},
		},
		{
			name:  "3 choose 3",
			items: []string{"A", "B", "C"},
			k:     3,
			want: [][]string{
				{"A", "B", "C"},
			},
		},
		{
			name:  "3 choose 0",
			items: []string{"A", "B", "C"},
			k:     0,
			want: [][]string{
				{},
			},
		},
		{
			name:  "3 choose 4",
			items: []string{"A", "B", "C"},
			k:     4,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got [][]string
			for _, combo := range Combinations(tt.items, tt.k) {
				got = append(got, combo)
			}
			if len(tt.want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combinations() = %v, want %v", got, tt.want)
			}
		})
	}
}
