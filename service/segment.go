package service

type Segment struct {
	Num      int
	Duration float64
}

type ByNum []*Segment

func (a ByNum) Len() int           { return len(a) }
func (a ByNum) Less(i, j int) bool { return a[i].Num < a[j].Num }
func (a ByNum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
