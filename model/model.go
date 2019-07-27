package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/ibidpayfazz/go-fazz-db/database/generic"
)

type ModelInterface interface {
	GetColumn() string
	GetPK() string
	GetLenColumn() int64
	SetColumn(column string)
	Save() (interface{}, error)
	SetModel(table string, data interface{})
}

type modelStruct struct {
	storage generic.GenericRepositoryInterface
	model   interface{}
	table   string
	key     string
	result  interface{}
}

func NewStorage(storage generic.GenericRepositoryInterface) *modelStruct {

	return &modelStruct{
		storage: storage,
	}
}

func (m *modelStruct) SetModel(table string, model interface{}) *modelStruct {

	// model UUID || model serial || model other
	// interfaceModel utk 3 model diatas
	// bikin fungsin interface yang kemungkinan 3 model ini pakai --> getPK --> data pk
	// check
	// m.data = data

	// switch data.(type) {
	// case *product.Product:
	// 	// generate uuid
	// 	Model := data.(*product.Product)
	// 	pkField := reflect.ValueOf(Model).Elem().FieldByName(m.GetPK())
	// 	fmt.Println(reflect.TypeOf(Model.ID), reflect.TypeOf(uuid.UUID{}))
	// 	if pkField.Type() == reflect.TypeOf(uuid.UUID{}) {
	// 		uid := uuid.New()
	// 		pkField.Set(reflect.ValueOf(uid))
	// 	}

	// 	data = Model
	// 	break

	// case product.Botol:
	// 	// serial
	// 	Model := data.(*product.Botol)
	// 	data = Model
	// 	break
	// case product.Gelas:
	// 	// others
	// 	Model := data.(*product.Gelas)
	// 	data = Model
	// 	break
	// }

	m.table = table
	m.model = model

	is := reflect.Indirect(reflect.ValueOf(model))
	ty := is.Type()

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		value := is.Field(i)

		if field.Type.String() == "uuid.UUID" && field.Tag.Get("db") == m.GetPK() {
			uid, err := uuid.NewV4()
			fmt.Println(uid)
			if err != nil {

			}
			value.Set(reflect.ValueOf(uid))
		}
	}

	return m
}

func (m *modelStruct) GetModel() interface{} {
	return m.model
}

func (m *modelStruct) GetTable() interface{} {
	return m.table
}

func (m *modelStruct) GetColumn() string {
	return ""
}

func (m *modelStruct) GetPK() string {
	v := reflect.ValueOf(m.model)
	is := reflect.Indirect(v)
	ty := is.Type()

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		dbName := field.Tag.Get("db")
		dbType := field.Tag.Get("dbType")
		if strings.Contains(dbType, "PRIMARY KEY") {
			return dbName
		}
	}

	return ""
}

func (m *modelStruct) GetLenColumn() int64 {
	return 0
}

func (m *modelStruct) SetColumn(colomn string) {
}

func (m *modelStruct) Save(returning ...string) (interface{}, error) {

	m.storage.SetTable(m.table)

	if len(returning) > 0 {
		return m.storage.Insert(m.model, returning...)
	} else {
		def := []string{m.GetPK()}
		res, _ := m.storage.Insert(m.model, def...)
		return res, nil
	}
}
