package globals

import (
	"encoding/json"
	"fmt"
	"os"
)

type Part struct {
	PartNumber       string
	MaterialCode     string
	Length           float32
	Quantity         uint16
	CuttingOperation string
}

type Material struct {
	MaterialCode     string
	Length           float32
	CuttingOperation string
	Quantity         uint16
}

type CutMaterial struct {
	Job          string
	MaterialCode string
	Parts        map[string]uint16
	Quantity     uint16
	StockLength  float32
	Length       float32
}

type CutMaterialPart struct {
	CutMaterialID           int
	Job                     string
	JobId                   int
	CutMaterialMaterialCode string
	CutMaterialQuantity     int
	StockLength             float64
	CutMaterialLength       float64
	PartID                  int
	PartNumber              string
	PartMaterialCode        string
	PartLength              float64
	PartQty                 int
}

type CutMaterialTotals struct {
	MaterialCode     string  `json:"material_code"`
	StockLength      float64 `json:"stock_length"`
	Length           float64 `json:"remaining_length"`   // length AS remaining_length
	TotalQuantity    int     `json:"total_quantity"`     // SUM(quantity)
	TotalStockLength float64 `json:"total_stock_length"` // stock_length * SUM(quantity)
	TotalUsedLength  float64 `json:"total_used_length"`  // stock_length - length * SUM(quantity)
}

// type LengthSortable interface {
// 	GetLength() float32
// }
//
// func (material Material) GetLength() float32 {
// 	return material.StockLength
// }
//
// func (cutMaterial CutMaterial) GetLength() float32 {
// 	return cutMaterial.RemainingLength
// }

// type ByLength []LengthSortable
//
// func (a ByLength) Len() int           { return len(a) }
// func (a ByLength) Less(i, j int) bool { return a[i].GetLength() > a[j].GetLength() }
// func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//
//	func SortResultsByLength(slice []LengthSortable) {
//		sort.Sort(ByLength(slice))
//	}

type LocalJobsList struct {
	JobNumber string
	Customer  string
	Star      int
}

type SettingsConfig struct {
	Kerf float32 `json:"kerf"`
}

var Settings SettingsConfig

func LoadSettings() (Settings SettingsConfig, err error) {
	// fmt.Println("Loading settings...")
	var filename = "./globals/settings.json"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening settings.json:", err)
		return Settings, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Settings)
	if err != nil {
		return Settings, err
	}
	// fmt.Println("Settings:", Settings)

	return Settings, nil
}

func SaveSettings(settings SettingsConfig) error {
	file, err := os.Create("settings.json")
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("Saving settings to settings.json...")
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: For pretty-printing JSON
	return encoder.Encode(settings)
}

// var Kerf float32 = .0625

var Parts = []Part{
	{
		PartNumber:       "Part1",
		MaterialCode:     "HSS3X3X.25",
		Length:           1.0,
		Quantity:         10,
		CuttingOperation: "Saw",
	},
	{
		PartNumber:       "Part2",
		MaterialCode:     "HSS3X3X.25",
		Length:           2,
		Quantity:         15,
		CuttingOperation: "Saw",
	},
	{
		PartNumber:       "Part3",
		MaterialCode:     "HSS3X3X.375",
		Length:           18.6,
		Quantity:         20,
		CuttingOperation: "Tube",
	},
	{
		PartNumber:       "Part4",
		MaterialCode:     "HSS3X3X.25",
		Length:           1.6,
		Quantity:         20,
		CuttingOperation: "Saw",
	},
}

var Materials = []Material{
	{
		MaterialCode:     "HSS3X3X.25",
		Length:           40,
		CuttingOperation: "",
		Quantity:         2,
	},
	{
		MaterialCode:     "HSS3X3X.25",
		Length:           10,
		CuttingOperation: "",
		Quantity:         20,
	},
	{
		MaterialCode:     "HSS3X3X.375",
		Length:           100,
		CuttingOperation: "",
		Quantity:         20,
	},
}

type JobType struct {
	Job      string
	Customer string
	Star     int
}

var JobInfo = JobType{
	Job:      "TEST1",
	Customer: "ATS",
}
