package material_utils

import (
	"testing"

	"main/globals"
)

func TestSortMaterialByLength(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Material
		expected []globals.Material
	}{
		{
			name: "Multiple materials",
			input: []globals.Material{
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "B", Length: 5},
				{MaterialCode: "C", Length: 15},
			},
			expected: []globals.Material{
				{MaterialCode: "B", Length: 5},
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "C", Length: 15},
			},
		},
		{
			name: "All materials same length",
			input: []globals.Material{
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "B", Length: 10},
				{MaterialCode: "C", Length: 10},
			},
			expected: []globals.Material{
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "B", Length: 10},
				{MaterialCode: "C", Length: 10},
			},
		},
		{
			name:     "Empty input",
			input:    []globals.Material{},
			expected: []globals.Material{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := SortMaterialByLength(tt.input)
				for i, mat := range got {
					if mat != tt.expected[i] {
						t.Errorf("SortMaterialByLength() = %v, want %v", got, tt.expected)
						break
					}
				}
			})
	}
}

func TestSortMaterialByCode(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Material
		target   string
		expected []globals.Material
		err      bool
	}{
		{
			name: "Materials with matching code",
			input: []globals.Material{
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "B", Length: 5},
				{MaterialCode: "A", Length: 15},
			},
			target: "A",
			expected: []globals.Material{
				{MaterialCode: "A", Length: 10},
				{MaterialCode: "A", Length: 15},
			},
			err: false,
		},
		{
			name:     "Materials with no matching code",
			input:    []globals.Material{{MaterialCode: "A", Length: 10}},
			target:   "B",
			expected: []globals.Material{},
			err:      true,
		},
		{
			name:     "Empty input",
			input:    []globals.Material{},
			target:   "A",
			expected: []globals.Material{},
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := SortMaterialByCode(tt.input, tt.target)
				if (err != nil) != tt.err {
					t.Errorf("SortMaterialByCode() error = %v, wantErr %v", err, tt.err)
					return
				}
				for i, mat := range got {
					if mat != tt.expected[i] {
						t.Errorf("SortMaterialByCode() = %v, want %v", got, tt.expected)
						break
					}
				}
			})
	}
}
