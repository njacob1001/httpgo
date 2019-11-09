package models

import (
	"database/sql"
	"log"
)

// Datastore interface
type Datastore interface {
	GetData() ([]*WeatherInfoRequest, error)
	InsertData(temp string, humed string) error
}

// DB database
type DB struct {
	*sql.DB
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
