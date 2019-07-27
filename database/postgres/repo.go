package postgres

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ibidpayfazz/go-fazz-db/database/generic"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var db *sqlx.DB

type repositoryStruct struct {
	table string
	key   string
}

// NewPostgresRepository .
func NewPostgresRepository(dbConnectionString string) generic.GenericRepositoryInterface {
	var err error
	db, err = sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return &repositoryStruct{}
}

func (s *repositoryStruct) SetTable(table string) {
	s.table = table
}

func (s *repositoryStruct) Insert(elem interface{}, returning ...string) (interface{}, error) {
	tn, ev, ip, tg, ct := "", reflect.ValueOf(elem), []string{}, []string{}, []string{}

	if len(s.table) > 0 {
		tn = s.table
	} else {
		ss := strings.Split(reflect.TypeOf(elem).String(), ".")
		tn = strings.ToLower(ss[len(ss)-1])
	}

	indirect := reflect.Indirect(ev)
	ty := indirect.Type()

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		dbTag := field.Tag.Get("db")
		dbType := field.Tag.Get("dbType")
		ip = append(ip, fmt.Sprintf(":%s", dbTag))
		tg = append(tg, dbTag)
		ct = append(ct, fmt.Sprintf("%s %s", dbTag, dbType))

	}

	st := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tn, strings.Join(tg, ","), strings.Join(ip, ","))

	if len(returning) > 0 {
		st = st + fmt.Sprintf(` RETURNING %s`, strings.Join(returning, ","))
	} else {
		st = st + fmt.Sprintf(` RETURNING (%s)`, getKeyQuery(tn))
	}

	stmt, err := db.PrepareNamed(st)

	if err != nil {
		_, isPqError := err.(*pq.Error)

		if !isPqError {

			return nil, err

		} else if err.(*pq.Error).Code == "42P01" {

			err = s.createTable(tn, ct)
			if err != nil {

				return nil, err
			}

			return s.Insert(elem, returning...)
		}
	}

	res, _ := stmt.Query(elem)

	if err != nil {
		return nil, err
	}

	data := make([]interface{}, len(returning))

	if len(returning) > 0 {
		addrs := make([]interface{}, len(returning))

		for x := range data {
			addrs[x] = &data[x]
		}
		if res.Next() {
			res.Scan(addrs...)
		}
	}

	for i, v := range data {
		switch v.(type) {
		case []byte:
			data[i] = string(v.([]byte))
			break
		}
	}

	return data, err

	// var ids string

	// if res.Next() {
	// 	res.Scan(&ids)
	// }

	// return ids, err
}

func (s *repositoryStruct) Update() {
	fmt.Println("UPDATE colomn FROM table VALUES ()")
}

func (s *repositoryStruct) Delete() {
	fmt.Println("DELETE table VALUES ()")
}

func (s *repositoryStruct) createTable(tableName string, colomn []string) error {
	var isExist string
	err := db.QueryRow(`SELECT to_regclass($1) as name`, tableName).Scan(&isExist)

	if err == nil {
		return errors.New(tableName + "is already exist")
	}

	schema := fmt.Sprintf(`CREATE TABLE %s (%s);`, tableName, strings.Join(colomn, ","))
	// execute a query on the server
	_, err = db.Exec(schema)

	if err != nil {
		return err
	}

	return nil
}

func selectFields(elemType reflect.Type) string {
	dbFields := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			dbFields = append(dbFields, fmt.Sprintf("\"%s\"", dbTag))
		}
	}
	return strings.Join(dbFields, ",")
}

func getKeyQuery(table string) string {
	return fmt.Sprintf(`SELECT pg_attribute.attname as key FROM pg_index, pg_class, pg_attribute, pg_namespace WHERE 
	pg_class.oid = to_regclass('%s') AND 
	indrelid = pg_class.oid AND 
	nspname = 'public' AND 
	pg_class.relnamespace = pg_namespace.oid AND 
	pg_attribute.attrelid = pg_class.oid AND 
	pg_attribute.attnum = any(pg_index.indkey)
   AND indisprimary`, table)
}
