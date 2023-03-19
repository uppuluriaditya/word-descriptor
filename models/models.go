package models

type Record struct {
	Shabdam    string
	Vibhakthi  string
	Ekavachan  string
	Dvivachan  string
	Bahuvachan string
	Nirdesh    string
	Lingam     string
}

type NounForm struct {
	Forms []Record
}

func (n NounForm) String() string {
	return n.Forms[0].Shabdam
}
