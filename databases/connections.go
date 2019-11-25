package databases

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	clients (
		id SERIAL PRIMARY KEY,
		cash INT DEFAULT 0,
		username varchar(40),
		password varchar(40)
	)
`

	// 	defaultValue := `
	// 	INSERT INTO
	// 	clients
	// 	(cash, username, password)
	// 	VALUES
	// 	(100000, 'redes', 'redes')
	// 	ON CONFLICT DO NOTHING
	// `

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
	// if _, err := db.Exec(defaultValue); err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	fmt.Println("Connected to Posgres!")
	return &DB{db}, nil
}
