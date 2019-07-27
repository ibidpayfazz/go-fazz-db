package mongo

import (
	"context"
	"errors"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/ibidpayfazz/go-fazz-db/database/generic"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repositoryStruct struct {
	table string
}

var db *mongo.Database

func NewMongoRepository(dbConnectionString string) generic.GenericRepositoryInterface {

	u, err := url.Parse(dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	path := strings.Split(u.Path, "/")

	// Set client options
	clientOptions := options.Client().ApplyURI(dbConnectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if len(path) > 1 {
		db = client.Database(path[1])
	} else {
		log.Fatal(errors.New("no database selected"))
	}

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return &repositoryStruct{}
}

func (s *repositoryStruct) SetTable(table string) {
	s.table = table
}

func (s *repositoryStruct) Insert(elem interface{}, returning ...string) (interface{}, error) {
	var result []interface{}
	ss := strings.Split(reflect.TypeOf(elem).String(), ".")
	tn := strings.ToLower(ss[len(ss)-1])

	collection := db.Collection(tn)
	res, err := collection.InsertOne(context.TODO(), elem)
	if err != nil {
		log.Fatal(err)
	}

	if len(returning) > 0 {
		for _, v := range returning {
			if strings.ToLower(v) == "id" {
				result := append(result, res.InsertedID.(primitive.ObjectID).Hex())
				return result, nil
			} else {
				return result, nil
			}
		}

		return res, err
	}

	return res, err
}

func (s *repositoryStruct) Update() {
}

func (s *repositoryStruct) Delete() {
}
