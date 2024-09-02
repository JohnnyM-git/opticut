package globals

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
	Job              string
	MaterialCode     string
	Parts            map[string]PartQTY
	Quantity         uint16
	StockLength      float64
	Length           float64
	CuttingOperation string
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
	CuttingOperation string  `json:"cutting_operation"`
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
	CuttingOperation        string  `json:"cutting_operation"`
}

type CutMaterialTotals struct {
	Id               int     `json:"id"`
	MaterialCode     string  `json:"material_code"`
	StockLength      float64 `json:"stock_length"`
	Length           float64 `json:"remaining_length"`   // length AS remaining_length
	TotalQuantity    uint16  `json:"total_quantity"`     // SUM(quantity)
	TotalStockLength float64 `json:"total_stock_length"` // stock_length * SUM(quantity)
	TotalUsedLength  float64 `json:"total_used_length"`  // stock_length - length * SUM(quantity)
	CuttingOperation string  `json:"cutting_operation"`
}

type LocalJobsList struct {
	JobNumber string
	Customer  string
	Star      int
}

type ExcelSettings struct {
	FilesPath          string
	PartSheet          string
	PartHeaderRows     uint16
	MaterialSheet      string
	MaterialHeaderRows uint16
}

type ExcelFileData struct {
	Job       JobType
	Parts     []Part     `json:"parts"`
	Materials []Material `json:"materials"`
}

type SettingsConfig struct {
	Kerf  float64       `json:"kerf"`
	Units string        `json:"units"`
	Excel ExcelSettings `json:"excelSettings"`
}

type JobType struct {
	JobId    int
	Job      string
	Customer string
	Star     int
}
