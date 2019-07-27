package main

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/ibidpayfazz/go-fazz-db/database/postgres"
	"github.com/ibidpayfazz/go-fazz-db/model"

	_ "github.com/lib/pq"
)

func main() {

	type Test struct {
		ID    uuid.UUID `json:"id" db:"id" dbType:"uuid PRIMARY KEY"`
		Name  string    `json:"name" db:"name" dbType:"varchar(255)"`
		Price int       `json:"price" db:"price" dbType:"integer"`
	}

	storage := postgres.NewPostgresRepository(fmt.Sprintf(`user=%s password=%s dbname=%s sslmode=disable`, "postgres", "postgres", "postgres"))
	//storage := mongo.NewMongoRepository(fmt.Sprintf(`mongodb://localhost:27017/%s`, "mongose"))

	m, _ := model.NewStorage(storage).SetModel("Testing", &Test{Name: "Testing"}).Save()

	fmt.Println(m)

}
