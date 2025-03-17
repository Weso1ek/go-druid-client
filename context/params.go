package context

type InputParams struct {
	Pm          int
	PmCategory  int
	Site        []int
	Category    []int
	Source      []int
	Channel     []int
	DateStart   int32
	DateEnd     int32
	Granulation string
}
