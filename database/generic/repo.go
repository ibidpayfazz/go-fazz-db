package generic

type genericStruct struct {
	table string
}

type GenericRepositoryInterface interface {
	SetTable(table string)
	Insert(elem interface{}, returning ...string) (interface{}, error)
	Update()
	Delete()
}


