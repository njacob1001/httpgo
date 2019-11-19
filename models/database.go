package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Datastore interface
type Datastore interface {
	GetData() ([]*Basura, error)
	GetDataBy(isUrgent bool) ([]*Basura, error)
	InsertData(info *BasuraFromArduino) error
	// InsertData(temp string, humed string) error
}

// DB database
type DB struct {
	*sql.DB
}

// MDB struct
type MDB struct {
	*mongo.Client
}

// NewMongoClient ok
func NewMongoClient() (*MDB, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return &MDB{client}, nil
}

// NewDB create a connection with a postgressdatabase
func NewDB(dataSourceName string) (*DB, error) {

	initializeTable := `
	CREATE TABLE IF NOT EXISTS
	data (
		id SERIAL PRIMARY KEY,
		temperatura numeric DEFAULT 0.0,
		humedad numeric DEFAULT 0.0,
		fecha TIMESTAMP DEFAULT NOW() 
	)
`

	defaultValue := `
	INSERT INTO 
	data 
	(temperatura, humedad, fecha) 
	VALUES 
	(1.2, 3.4, '2017-03-31 09:30:20-07')
	ON CONFLICT DO NOTHING
`

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(initializeTable); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if _, err := db.Exec(defaultValue); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &DB{db}, nil
}
