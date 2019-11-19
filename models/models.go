package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WeatherInfoRequest for the struct
type WeatherInfoRequest struct {
	ID          int16   `json:"id"`
	Temperatura float32 `json:"temperatura"`
	Humedad     float32 `json:"humedad"`
	Fecha       string  `json:"fecha"`
}

// WeatherInfo of struct
type WeatherInfo struct {
	Temperatura float32 `json:"temperatura"`
	Humedad     float32 `json:"humedad"`
	Fecha       string  `json:"fecha"`
}

// FullResponse of struct
type FullResponse struct {
	Status string               `json:"status"`
	Data   []WeatherInfoRequest `json:"data"`
}

// BasuraFromArduino struct
type BasuraFromArduino struct {
	Nocivo        bool    `json:"nocivo"`
	Identificador string  `json:"identificador"`
	Peso          float32 `json:"peso"`
}

// Basura struct
type Basura struct {
	Identificador string  `json:"identificador"`
	Peso          float32 `json:"peso"`
	Nocivo        bool    `json:"nocivo"`
	Timestamp     string  `json:"timestamp"`
	Grado         string  `json:"grado"`
}

// GetData function
func (db *MDB) GetData() ([]*Basura, error) {
	findOptions := options.Find()
	cur, err := db.Database("redesuao").Collection("basuras").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var results []*Basura
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Basura
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return results, nil
}

// GetDataBy function
func (db *MDB) GetDataBy(isUrgent bool) ([]*Basura, error) {
	findOptions := options.Find()
	filter := bson.D{{
		"grado",
		bson.D{{
			"$in",
			bson.A{"recoger"},
		}},
	}}
	if isUrgent {
		filter = bson.D{{
			"grado",
			bson.D{{
				"$in",
				bson.A{"urgente"},
			}},
		}}
	}
	cur, err := db.Database("redesuao").Collection("basuras").Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var results []*Basura
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Basura
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	return results, nil
}

// InsertData function
func (db *MDB) InsertData(info *BasuraFromArduino) error {
	typeData := "no recoger"
	t := time.Now()
	stringTime := t.Format("2006-01-02T15:04:05-0700")
	if info.Nocivo || info.Peso > 5000 {
		typeData = "urgente"
	} else if info.Peso > 2000 && info.Peso <= 5000 {
		fmt.Print("Recpger")
		typeData = "recoger"
	}
	toSave := Basura{
		Identificador: info.Identificador,
		Peso:          info.Peso,
		Nocivo:        info.Nocivo,
		Timestamp:     stringTime,
		Grado:         typeData,
	}
	_, err := db.Database("redesuao").Collection("basuras").InsertOne(context.TODO(), toSave)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// // InsertData function
// func (db *MDB) InsertData(temp string, humed string) error {
// 	_, err := db.Exec(`
// 		INSERT INTO data (temperatura, humedad)  VALUES (
// 			$1, $2
// 		)
// 	`, temp, humed)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	return nil
// }
