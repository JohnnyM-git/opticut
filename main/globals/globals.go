package globals

import (
	"encoding/json"
	"fmt"
	"os"
)

type Part struct {
	PartNumber       string
	MaterialCode     string
	Length           float64
	Quantity         uint16
	CuttingOperation string
	CutQuantity      uint16
}

type Material struct {
	MaterialCode     string
	Length           float64
	CuttingOperation string
	Quantity         uint16
}

type PartQTY struct {
	CurrentQty uint16
	TotalQty   uint16
}

type CutMaterial struct {
	Job          string
	MaterialCode string
	Parts        map[string]PartQTY
	Quantity     uint16
	StockLength  float64
	Length       float64
}

type CutMaterialPart struct {
	CutMaterialID    int     `json:"cut_material_id"`
	PartID           int     `json:"part_id"`
	PartNumber       string  `json:"part_number"`
	PartMaterialCode string  `json:"part_material_code"`
	PartLength       float64 `json:"part_length"`
	PartCutLength    float64 `json:"part_cut_length"`
	PartQty          int     `json:"part_qty"`
	TotalPartQty     int     `json:"total_part_qty"`
}

type CutMaterials struct {
	CutMaterialID           int     `json:"cut_material_id"`
	Job                     string  `json:"job"`
	JobId                   int     `json:"job_id"`
	CutMaterialMaterialCode string  `json:"cut_material_material_code"`
	CutMaterialQuantity     int     `json:"cut_material_quantity"`
	StockLength             float64 `json:"stock_length"`
	CutMaterialLength       float64 `json:"cut_material_length"`
	TotalUsedLength         float64 `json:"total_used_length"`
	TotalPartsCutOnMaterial int     `json:"total_parts_cut_on_material"`
}

type CutMaterialTotals struct {
	Id               int     `json:"id"`
	MaterialCode     string  `json:"material_code"`
	StockLength      float64 `json:"stock_length"`
	Length           float64 `json:"remaining_length"`   // length AS remaining_length
	TotalQuantity    uint16  `json:"total_quantity"`     // SUM(quantity)
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

type ExcelSettings struct {
	SheetName  string
	HeaderRows uint16
}

type SettingsConfig struct {
	Kerf  float64       `json:"kerf"`
	Units string        `json:"units"`
	Excel ExcelSettings `json:"excelSettings"`
}

var Settings SettingsConfig

func LoadSettings() (SettingsConfig, error) {
	fmt.Println("Loading settings...")
	var filename = "./globals/settings.json"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening settings.json:", err)
		return SettingsConfig{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Settings)
	if err != nil {
		return SettingsConfig{}, err
	}
	fmt.Println("Settings:", Settings)

	return Settings, nil
}

func SaveSettings(settings SettingsConfig) error {
	var filename = "./globals/settings.json"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("Saving settings to settings.json...")
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: For pretty-printing JSON
	return encoder.Encode(settings)
}

var Parts = []Part{
	{
		PartNumber:       "Part1",
		MaterialCode:     "HSS3X3X.25",
		Length:           1,
		Quantity:         8,
		CuttingOperation: "Saw",
	},
	{
		PartNumber:       "Part2",
		MaterialCode:     "HSS3X3X.25",
		Length:           2,
		Quantity:         3,
		CuttingOperation: "Saw",
	},
	// {
	// 	PartNumber:       "Part3",
	// 	MaterialCode:     "HSS3X3X.375",
	// 	Length:           18.6,
	// 	Quantity:         20,
	// 	CuttingOperation: "Tube",
	// },
	// {
	// 	PartNumber:       "Part4",
	// 	MaterialCode:     "HSS3X3X.25",
	// 	Length:           1.6,
	// 	Quantity:         20,
	// 	CuttingOperation: "Saw",
	// },
}

var Materials = []Material{
	// {
	// 	MaterialCode:     "HSS3X3X.25",
	// 	Length:           40,
	// 	CuttingOperation: "",
	// 	Quantity:         2,
	// },
	{
		MaterialCode:     "HSS3X3X.25",
		Length:           10,
		CuttingOperation: "",
		Quantity:         20,
	},
	// {
	// 	MaterialCode:     "HSS3X3X.375",
	// 	Length:           100,
	// 	CuttingOperation: "",
	// 	Quantity:         20,
	// },
}

type JobType struct {
	JobId    int
	Job      string
	Customer string
	Star     int
}

var JobInfo = JobType{
	Job:      "TEST1",
	Customer: "ATS",
}
