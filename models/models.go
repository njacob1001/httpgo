package models

import (
	"fmt"
	"log"
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

// GetData function
func (db *DB) GetData() (FullResponse, error) {
	var allData []WeatherInfoRequest
	rows, err := db.Query("SELECT * FROM data")
	fmt.Println(rows)
	fmt.Println(allData)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&allData); err != nil {
			log.Fatal(err)
		}
	}

	defer rows.Close()
	response := FullResponse{
		Status: "ok",
		Data:   allData,
	}
	// iterate over the result and print out the titles
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Fatal(err)
		}
		fmt.Println(title)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return response, err
}

// InsertData function
func (db *DB) InsertData(info *WeatherInfo) error {
	_, err := db.Exec(`
		INSERT INTO data (temperatura, humedad)  VALUES (
			$1, $2, $3
		)
	`, info.Temperatura, info.Humedad)
	if err != nil {
		log.Println(err)
	}

	return nil
}
