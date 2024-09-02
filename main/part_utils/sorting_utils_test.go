package part_utils

import (
	"reflect"
	"testing"

	"main/globals"
)

func TestSortPartByLength(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Part
		expected []globals.Part
	}{
		{
			name: "Already sorted parts",
			input: []globals.Part{
				{PartNumber: "1", Length: 10.5},
				{PartNumber: "2", Length: 9.2},
				{PartNumber: "3", Length: 7.3},
			},
			expected: []globals.Part{
				{PartNumber: "1", Length: 10.5},
				{PartNumber: "2", Length: 9.2},
				{PartNumber: "3", Length: 7.3},
			},
		},
		{
			name: "Unsorted parts",
			input: []globals.Part{
				{PartNumber: "1", Length: 7.3},
				{PartNumber: "2", Length: 10.5},
				{PartNumber: "3", Length: 9.2},
			},
			expected: []globals.Part{
				{PartNumber: "2", Length: 10.5},
				{PartNumber: "3", Length: 9.2},
				{PartNumber: "1", Length: 7.3},
			},
		},
		{
			name: "All parts with the same length",
			input: []globals.Part{
				{PartNumber: "1", Length: 10.0},
				{PartNumber: "2", Length: 10.0},
				{PartNumber: "3", Length: 10.0},
			},
			expected: []globals.Part{
				{PartNumber: "1", Length: 10.0},
				{PartNumber: "2", Length: 10.0},
				{PartNumber: "3", Length: 10.0},
			},
		},
		{
			name: "Single part",
			input: []globals.Part{
				{PartNumber: "1", Length: 8.5},
			},
			expected: []globals.Part{
				{PartNumber: "1", Length: 8.5},
			},
		},
		{
			name:     "Empty input",
			input:    []globals.Part{},
			expected: []globals.Part{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := SortPartByLength(tt.input)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("SortPartByLength() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func TestSortPartsByCode(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Part
		expected [][]globals.Part
	}{
		{
			name: "Multiple parts with different codes",
			input: []globals.Part{
				{PartNumber: "1", MaterialCode: "A", Length: 10.5},
				{PartNumber: "2", MaterialCode: "B", Length: 7.3},
				{PartNumber: "3", MaterialCode: "A", Length: 8.0},
				{PartNumber: "4", MaterialCode: "B", Length: 9.2},
			},
			expected: [][]globals.Part{
				{
					{PartNumber: "1", MaterialCode: "A", Length: 10.5},
					{PartNumber: "3", MaterialCode: "A", Length: 8.0},
				},
				{
					{PartNumber: "4", MaterialCode: "B", Length: 9.2},
					{PartNumber: "2", MaterialCode: "B", Length: 7.3},
				},
			},
		},
		{
			name: "Parts with the same code and length",
			input: []globals.Part{
				{PartNumber: "1", MaterialCode: "A", Length: 10.0},
				{PartNumber: "2", MaterialCode: "A", Length: 10.0},
				{PartNumber: "3", MaterialCode: "B", Length: 9.0},
			},
			expected: [][]globals.Part{
				{
					{PartNumber: "1", MaterialCode: "A", Length: 10.0},
					{PartNumber: "2", MaterialCode: "A", Length: 10.0},
				},
				{
					{PartNumber: "3", MaterialCode: "B", Length: 9.0},
				},
			},
		},
		{
			name: "Single part",
			input: []globals.Part{
				{PartNumber: "1", MaterialCode: "A", Length: 8.5},
			},
			expected: [][]globals.Part{
				{
					{PartNumber: "1", MaterialCode: "A", Length: 8.5},
				},
			},
		},
		{
			name:     "Empty input",
			input:    []globals.Part{},
			expected: [][]globals.Part{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := SortPartsByCode(tt.input)

				// Check that the grouped parts match, ignoring the order of the groups.
				for _, expectedGroup := range tt.expected {
					found := false
					for _, resultGroup := range result {
						if reflect.DeepEqual(resultGroup, expectedGroup) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("SortPartsByCode() group %v not found in result %v", expectedGroup, result)
					}
				}

				// Check if the overall lengths match.
				if len(result) != len(tt.expected) {
					t.Errorf("SortPartsByCode() = %v, want %v", result, tt.expected)
				}
			})
	}
}
