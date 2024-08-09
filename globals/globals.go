package globals

type Part struct {
	PartNumber       string
	MaterialCode     string
	Length           float32
	Quantity         uint16
	CuttingOperation string
}

type Material struct {
	MaterialCode     string
	StockLength      float32
	CuttingOperation string
	Quantity         uint16
}

type CutMaterials struct {
	MaterialCode    string
	Parts           []Part
	Quantity        uint16
	RemainingLength float32
}

var Kerf float32 = .0625

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
		Length:           12.3,
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
		StockLength:      240,
		CuttingOperation: "Saw",
		Quantity:         2,
	},
	{
		MaterialCode:     "HSS3X3X.25",
		StockLength:      40,
		CuttingOperation: "Saw",
		Quantity:         20,
	},
}
