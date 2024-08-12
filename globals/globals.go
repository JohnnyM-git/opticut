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
// func SortResultsByLength(slice []LengthSortable) {
// 	sort.Sort(ByLength(slice))
// }

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
		Length:           40,
		CuttingOperation: "Saw",
		Quantity:         2,
	},
	{
		MaterialCode:     "HSS3X3X.25",
		Length:           40,
		CuttingOperation: "Saw",
		Quantity:         20,
	},
}
