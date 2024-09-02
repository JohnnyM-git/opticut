package optimizer

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"

	"main/globals"
)

func TestPartsComplete(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Part
		expected bool
	}{
		{
			name: "All Parts Complete",
			input: []globals.Part{
				{PartNumber: "Part1", MaterialCode: "HSS3X3X.25", Length: 1, Quantity: 8, CutQuantity: 8},
				{PartNumber: "Part2", MaterialCode: "HSS3X3X.25", Length: 2, Quantity: 3, CutQuantity: 3},
			},
			expected: true,
		},
		{
			name: "Some Parts Incomplete",
			input: []globals.Part{
				{PartNumber: "Part1", MaterialCode: "HSS3X3X.25", Length: 1, Quantity: 8, CutQuantity: 7},
				{PartNumber: "Part2", MaterialCode: "HSS3X3X.25", Length: 2, Quantity: 3, CutQuantity: 3},
			},
			expected: false,
		},
		{
			name:     "No Parts",
			input:    []globals.Part{},
			expected: true,
		},
		{
			name: "All Parts Incomplete",
			input: []globals.Part{
				{PartNumber: "Part1", MaterialCode: "HSS3X3X.25", Length: 1, Quantity: 8, CutQuantity: 5},
				{PartNumber: "Part2", MaterialCode: "HSS3X3X.25", Length: 2, Quantity: 3, CutQuantity: 2},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := PartsComplete(tt.input)
				if got != tt.expected {
					t.Errorf("PartsComplete() = %v, want %v", got, tt.expected)
				}
			})
	}
}

func TestNoPartsFit(t *testing.T) {
	tests := []struct {
		name            string
		parts           []globals.Part
		currentMaterial *globals.CutMaterial
		expected        bool
	}{
		{
			name: "All Parts Already Cut",
			parts: []globals.Part{
				{
					PartNumber:       "Part1",
					MaterialCode:     "HSS3X3X.25",
					Length:           1,
					Quantity:         8,
					CutQuantity:      8,
					CuttingOperation: "Cut",
				},
			},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           1,
				CuttingOperation: "Cut",
			},
			expected: true, // All Parts have been cut but not yet hit the Parts Complete
		},
		{
			name: "Parts Fit Exactly",
			parts: []globals.Part{
				{
					PartNumber:       "Part1",
					MaterialCode:     "HSS3X3X.25",
					Length:           1,
					Quantity:         8,
					CutQuantity:      7,
					CuttingOperation: "Cut",
				},
			},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           1,
				CuttingOperation: "Cut",
			},
			expected: false, // Parts fit exactly on the material
		},
		{
			name: "No Parts Fit",
			parts: []globals.Part{
				{
					PartNumber:       "Part1",
					MaterialCode:     "HSS3X3X.25",
					Length:           2,
					Quantity:         5,
					CutQuantity:      5,
					CuttingOperation: "Cut",
				},
			},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           1,
				CuttingOperation: "Cut",
			},
			expected: true, // No parts fit because the length is greater
		},
		{
			name: "Part Does Not Match Cutting Operation",
			parts: []globals.Part{
				{
					PartNumber:       "Part1",
					MaterialCode:     "HSS3X3X.25",
					Length:           1,
					Quantity:         5,
					CutQuantity:      5,
					CuttingOperation: "Tube Laser",
				},
			},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           1,
				CuttingOperation: "Cut",
			},
			expected: true, // No parts fit due to a different cutting operation
		},
		{
			name: "All Parts Fit",
			parts: []globals.Part{
				{
					PartNumber:       "Part1",
					MaterialCode:     "HSS3X3X.25",
					Length:           1,
					Quantity:         8,
					CutQuantity:      7,
					CuttingOperation: "Cut",
				},
				{
					PartNumber:       "Part2",
					MaterialCode:     "HSS3X3X.25",
					Length:           2,
					Quantity:         3,
					CutQuantity:      2,
					CuttingOperation: "Cut",
				},
			},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           2,
				CuttingOperation: "Cut",
			},
			expected: false, // Parts fit on the material
		},
		{
			name:  "No Parts and No Material",
			parts: []globals.Part{},
			currentMaterial: &globals.CutMaterial{
				MaterialCode:     "HSS3X3X.25",
				Length:           2,
				CuttingOperation: "Cut",
			},
			expected: true, // No parts to fit
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NoPartsFit(tt.parts, tt.currentMaterial)
				if got != tt.expected {
					t.Errorf("NoPartsFit() = %v, want %v", got, tt.expected)
				}
			})
	}
}

func TestSortCutMaterialsByLength(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.CutMaterial
		expected []globals.CutMaterial
	}{
		{
			name: "Sort CutMaterials by Length",
			input: []globals.CutMaterial{
				{Length: 12.5},
				{Length: 5.0},
				{Length: 7.75},
			},
			expected: []globals.CutMaterial{
				{Length: 5.0},
				{Length: 7.75},
				{Length: 12.5},
			},
		},
		{
			name: "Already Sorted",
			input: []globals.CutMaterial{
				{Length: 1.0},
				{Length: 2.0},
				{Length: 3.0},
			},
			expected: []globals.CutMaterial{
				{Length: 1.0},
				{Length: 2.0},
				{Length: 3.0},
			},
		},
		{
			name:     "Empty Slice",
			input:    []globals.CutMaterial{},
			expected: []globals.CutMaterial{},
		},
		{
			name: "Single Element",
			input: []globals.CutMaterial{
				{Length: 10.0},
			},
			expected: []globals.CutMaterial{
				{Length: 10.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				SortCutMaterialsByLength(tt.input)
				for i, v := range tt.input {
					if v.Length != tt.expected[i].Length {
						t.Errorf("SortCutMaterialsByLength() = %v, want %v", tt.input, tt.expected)
						break
					}
				}
			})
	}
}

func TestSortMaterialsByLength(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.Material
		expected []globals.Material
	}{
		{
			name: "Sort Materials by Length",
			input: []globals.Material{
				{Length: 12.5},
				{Length: 5.0},
				{Length: 7.75},
			},
			expected: []globals.Material{
				{Length: 5.0},
				{Length: 7.75},
				{Length: 12.5},
			},
		},
		{
			name: "Already Sorted",
			input: []globals.Material{
				{Length: 1.0},
				{Length: 2.0},
				{Length: 3.0},
			},
			expected: []globals.Material{
				{Length: 1.0},
				{Length: 2.0},
				{Length: 3.0},
			},
		},
		{
			name:     "Empty Slice",
			input:    []globals.Material{},
			expected: []globals.Material{},
		},
		{
			name: "Single Element",
			input: []globals.Material{
				{Length: 10.0},
			},
			expected: []globals.Material{
				{Length: 10.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				SortMaterialsByLength(tt.input)
				for i, v := range tt.input {
					if v.Length != tt.expected[i].Length {
						t.Errorf("SortMaterialsByLength() = %v, want %v", tt.input, tt.expected)
						break
					}
				}
			})
	}
}

func TestArePartsEqual(t *testing.T) {
	tests := []struct {
		name     string
		parts1   map[string]globals.PartQTY
		parts2   map[string]globals.PartQTY
		expected bool
	}{
		{
			name: "Equal maps",
			parts1: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
				"part2": {CurrentQty: 3, TotalQty: 6},
			},
			parts2: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
				"part2": {CurrentQty: 3, TotalQty: 6},
			},
			expected: true,
		},
		{
			name: "Different sizes",
			parts1: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
			},
			parts2: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
				"part2": {CurrentQty: 3, TotalQty: 6},
			},
			expected: false,
		},
		{
			name: "Different quantities",
			parts1: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
				"part2": {CurrentQty: 3, TotalQty: 6},
			},
			parts2: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
				"part2": {CurrentQty: 2, TotalQty: 6},
			},
			expected: false,
		},
		{
			name: "Different keys",
			parts1: map[string]globals.PartQTY{
				"part1": {CurrentQty: 5, TotalQty: 10},
			},
			parts2: map[string]globals.PartQTY{
				"part2": {CurrentQty: 5, TotalQty: 10},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := arePartsEqual(tt.parts1, tt.parts2)
				if result != tt.expected {
					t.Errorf("arePartsEqual() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func TestMergeDuplicateCutMaterialsInPlace(t *testing.T) {
	tests := []struct {
		name     string
		input    []globals.CutMaterial
		expected []globals.CutMaterial
	}{
		{
			name: "No duplicates",
			input: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part1": {CurrentQty: 2, TotalQty: 5},
					},
				},
				{
					MaterialCode: "B",
					StockLength:  20.0,
					Length:       10.0,
					Quantity:     2,
					Parts: map[string]globals.PartQTY{
						"part2": {CurrentQty: 1, TotalQty: 3},
					},
				},
			},
			expected: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part1": {CurrentQty: 2, TotalQty: 5},
					},
				},
				{
					MaterialCode: "B",
					StockLength:  20.0,
					Length:       10.0,
					Quantity:     2,
					Parts: map[string]globals.PartQTY{
						"part2": {CurrentQty: 1, TotalQty: 3},
					},
				},
			},
		},
		{
			name: "With duplicates",
			input: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part1": {CurrentQty: 2, TotalQty: 5},
					},
				},
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part1": {CurrentQty: 2, TotalQty: 5},
					},
				},
				{
					MaterialCode: "B",
					StockLength:  20.0,
					Length:       10.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part2": {CurrentQty: 1, TotalQty: 3},
					},
				},
			},
			expected: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     2, // Merged quantity
					Parts: map[string]globals.PartQTY{
						"part1": {CurrentQty: 2, TotalQty: 5}, // Merged quantities
					},
				},
				{
					MaterialCode: "B",
					StockLength:  20.0,
					Length:       10.0,
					Quantity:     1,
					Parts: map[string]globals.PartQTY{
						"part2": {CurrentQty: 1, TotalQty: 3},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Execute the function
				mergeDuplicateCutMaterialsInPlace(&tt.input)

				// Verify the result
				if !reflect.DeepEqual(tt.input, tt.expected) {
					t.Errorf("mergeDuplicateCutMaterialsInPlace() = %+v, want %+v", tt.input, tt.expected)
				}
			})
	}
}

func TestCheckForMaterial(t *testing.T) {
	tests := []struct {
		name              string
		part              globals.Part
		results           []globals.CutMaterial
		materials         []globals.Material
		job               string
		expectedIndex     int
		expectedError     string // Change to string for direct comparison
		expectedResults   []globals.CutMaterial
		expectedMaterials []globals.Material
	}{
		{
			name: "No results, no materials",
			part: globals.Part{
				MaterialCode: "A",
				Length:       5.0,
			},
			results:           []globals.CutMaterial{},
			materials:         []globals.Material{},
			job:               "job1",
			expectedIndex:     0,
			expectedError:     `A: no material at the correct length of 5 found needed`,
			expectedResults:   []globals.CutMaterial{},
			expectedMaterials: []globals.Material{},
		},
		{
			name: "Matching CutMaterial found",
			part: globals.Part{
				MaterialCode: "A",
				Length:       5.0,
			},
			results: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts:        map[string]globals.PartQTY{},
				},
			},
			materials:     []globals.Material{},
			job:           "job1",
			expectedIndex: 0,
			expectedError: "",
			expectedResults: []globals.CutMaterial{
				{
					MaterialCode: "A",
					StockLength:  10.0,
					Length:       5.0,
					Quantity:     1,
					Parts:        map[string]globals.PartQTY{},
				},
			},
			expectedMaterials: []globals.Material{},
		},
		{
			name: "Matching Material found",
			part: globals.Part{
				MaterialCode: "A",
				Length:       5.0,
			},
			results: []globals.CutMaterial{},
			materials: []globals.Material{
				{
					MaterialCode: "B",
					Length:       10.0,
					Quantity:     1,
				},
			},
			job:           "job1",
			expectedIndex: 0,
			expectedError: "",
			expectedResults: []globals.CutMaterial{
				{
					Job:          "job1",
					MaterialCode: "B",
					StockLength:  10.0,
					Length:       10.0,
					Quantity:     1,
					Parts:        map[string]globals.PartQTY{},
				},
			},
			expectedMaterials: []globals.Material{
				{
					MaterialCode: "B",
					Length:       10.0,
					Quantity:     0,
				},
			},
		},
		{
			name: "No matching material",
			part: globals.Part{
				MaterialCode: "A",
				Length:       15.0,
			},
			results: []globals.CutMaterial{},
			materials: []globals.Material{
				{
					MaterialCode: "C",
					Length:       10.0,
					Quantity:     1,
				},
			},
			job:             "job1",
			expectedIndex:   0,
			expectedError:   `A: no material at the correct length of 15 found needed`,
			expectedResults: []globals.CutMaterial{},
			expectedMaterials: []globals.Material{
				{
					MaterialCode: "C",
					Length:       10.0,
					Quantity:     1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Capture the output
				var buffer bytes.Buffer
				log.SetOutput(&buffer)
				defer log.SetOutput(os.Stderr)

				// Execute the function
				index, err := checkForMaterial(&tt.part, &tt.results, &tt.materials, tt.job)

				// Verify results
				if index != tt.expectedIndex {
					t.Errorf("checkForMaterial() index = %d, want %d", index, tt.expectedIndex)
				}

				if err != nil && err.Error() != tt.expectedError {
					t.Errorf("checkForMaterial() error = %v, want %v", err, tt.expectedError)
				} else if err == nil && tt.expectedError != "" {
					t.Errorf("checkForMaterial() error = nil, want %v", tt.expectedError)
				}

				if !reflect.DeepEqual(tt.results, tt.expectedResults) {
					t.Errorf("checkForMaterial() results = %v, want %v", tt.results, tt.expectedResults)
				}

				if !reflect.DeepEqual(tt.materials, tt.expectedMaterials) {
					t.Errorf("checkForMaterial() materials = %v, want %v", tt.materials, tt.expectedMaterials)
				}
			})
	}
}

func TestCheckForMaterialV2(t *testing.T) {
	tests := []struct {
		name              string
		part              globals.Part
		currentMaterial   globals.CutMaterial
		materials         []globals.Material
		jobInfo           globals.JobType
		expectedError     string
		expectedMaterial  globals.CutMaterial
		expectedMaterials []globals.Material
	}{
		{
			name: "Materials list empty",
			part: globals.Part{
				MaterialCode: "A",
				Length:       5.0,
			},
			currentMaterial:   globals.CutMaterial{},
			materials:         []globals.Material{},
			jobInfo:           globals.JobType{Job: "job1"},
			expectedError:     "Materials list empty",
			expectedMaterial:  globals.CutMaterial{},
			expectedMaterials: []globals.Material{},
		},
		{
			name: "Matching material found, quantity not max",
			part: globals.Part{
				MaterialCode:     "A",
				Length:           5.0,
				CuttingOperation: "Cutting",
			},
			currentMaterial: globals.CutMaterial{},
			materials: []globals.Material{
				{
					MaterialCode: "A",
					Length:       10.0,
					Quantity:     10,
				},
			},
			jobInfo:       globals.JobType{Job: "job1"},
			expectedError: "",
			expectedMaterial: globals.CutMaterial{
				Job:              "job1",
				MaterialCode:     "A",
				Parts:            map[string]globals.PartQTY{},
				Quantity:         1,
				StockLength:      10.0,
				Length:           10.0,
				CuttingOperation: "Cutting",
			},
			expectedMaterials: []globals.Material{
				{
					MaterialCode: "A",
					Length:       10.0,
					Quantity:     9, // Decremented
				},
			},
		},
		{
			name: "Matching material found, quantity at max",
			part: globals.Part{
				MaterialCode:     "A",
				Length:           5.0,
				CuttingOperation: "Cutting",
			},
			currentMaterial: globals.CutMaterial{},
			materials: []globals.Material{
				{
					MaterialCode: "A",
					Length:       10.0,
					Quantity:     9999, // Max quantity
				},
			},
			jobInfo:       globals.JobType{Job: "job1"},
			expectedError: "",
			expectedMaterial: globals.CutMaterial{
				Job:              "job1",
				MaterialCode:     "A",
				Parts:            map[string]globals.PartQTY{},
				Quantity:         1,
				StockLength:      10.0,
				Length:           10.0,
				CuttingOperation: "Cutting",
			},
			expectedMaterials: []globals.Material{
				{
					MaterialCode: "A",
					Length:       10.0,
					Quantity:     9999, // Unchanged
				},
			},
		},
		{
			name: "No matching material found",
			part: globals.Part{
				MaterialCode: "A",
				Length:       15.0,
			},
			currentMaterial: globals.CutMaterial{},
			materials: []globals.Material{
				{
					MaterialCode: "B",
					Length:       10.0,
					Quantity:     1,
				},
			},
			jobInfo:          globals.JobType{Job: "job1"},
			expectedError:    "material with the correct length not found",
			expectedMaterial: globals.CutMaterial{},
			expectedMaterials: []globals.Material{
				{
					MaterialCode: "B",
					Length:       10.0,
					Quantity:     1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Initialize the currentMaterial to avoid comparison issues
				var currentMaterial globals.CutMaterial = tt.currentMaterial

				// Execute the function
				err := checkForMaterialV2(&tt.part, &currentMaterial, &tt.materials, tt.jobInfo)

				// Verify error
				if tt.expectedError != "" {
					if err == nil || err.Error() != tt.expectedError {
						t.Errorf("checkForMaterialV2() error = %v, want %v", err, tt.expectedError)
					}
				} else if err != nil {
					t.Errorf("checkForMaterialV2() unexpected error = %v", err)
				}

				// Verify updated material
				if !reflect.DeepEqual(currentMaterial, tt.expectedMaterial) {
					t.Errorf("checkForMaterialV2() currentMaterial = %v, want %v", currentMaterial, tt.expectedMaterial)
				}

				// Verify updated materials
				if !reflect.DeepEqual(tt.materials, tt.expectedMaterials) {
					t.Errorf("checkForMaterialV2() materials = %v, want %v", tt.materials, tt.expectedMaterials)
				}
			})
	}
}

// func TestShouldAppendCurrentMaterial(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		parts           []globals.Part
// 		currentMaterial globals.CutMaterial
// 		expected        bool
// 	}{
// 		{
// 			name: "All Parts Too Long",
// 			parts: []globals.Part{
// 				{
// 					PartNumber:       "P1",
// 					Length:           6.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 				{
// 					PartNumber:       "P2",
// 					Length:           8.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 			},
// 			currentMaterial: globals.CutMaterial{
// 				MaterialCode:     "A",
// 				Length:           5.0,
// 				CuttingOperation: "Cutting",
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "All Parts Don't Match Cutting Operation",
// 			parts: []globals.Part{
// 				{
// 					PartNumber:       "P1",
// 					Length:           6.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 				{
// 					PartNumber:       "P2",
// 					Length:           8.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 			},
// 			currentMaterial: globals.CutMaterial{
// 				MaterialCode:     "A",
// 				Length:           5.0,
// 				CuttingOperation: "Saw",
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "One Part Fits Length and Operation",
// 			parts: []globals.Part{
// 				{
// 					PartNumber:       "P1",
// 					Length:           6.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 				{
// 					PartNumber:       "P2",
// 					Length:           3.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      1,
// 					Quantity:         2,
// 				},
// 			},
// 			currentMaterial: globals.CutMaterial{
// 				MaterialCode:     "A",
// 				Length:           5.0,
// 				CuttingOperation: "Cutting",
// 			},
// 			expected: false,
// 		},
// 		{
// 			name: "Empty Current Material",
// 			parts: []globals.Part{
// 				{
// 					PartNumber:       "P1",
// 					Length:           6.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      2,
// 					Quantity:         1,
// 				},
// 				{
// 					PartNumber:       "P2",
// 					Length:           3.0,
// 					CuttingOperation: "Cutting",
// 					MaterialCode:     "A",
// 					CutQuantity:      2,
// 					Quantity:         1,
// 				},
// 			},
// 			currentMaterial: globals.CutMaterial{},
// 			expected:        false,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(
// 			tt.name, func(t *testing.T) {
// 				result := shouldAppendCurrentMaterial(tt.parts, &tt.currentMaterial)
// 				if result != tt.expected {
// 					t.Errorf("ShouldAppendCurrentMaterial() = %v, want %v", result, tt.expected)
// 				}
// 			})
// 	}
//
// }
