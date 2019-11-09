package models

import (
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
func (db *DB) GetData() ([]*WeatherInfoRequest, error) {
	rows, err := db.Query("SELECT * FROM data")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	servers := make([]*WeatherInfoRequest, 0)
	for rows.Next() {
		server := new(WeatherInfoRequest)
		err := rows.Scan(&server.ID, &server.Temperatura, &server.Humedad, &server.Fecha)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, server)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// var allData []WeatherInfoRequest
	// rows, err := db.Query("SELECT * FROM data")
	// fmt.Println(rows)
	// fmt.Println(allData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for rows.Next() {
	// 	if err := rows.Scan(&allData); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// defer rows.Close()
	// response := FullResponse{
	// 	Status: "ok",
	// 	Data:   servers,
	// }
	// // iterate over the result and print out the titles
	// for rows.Next() {
	// 	var title string
	// 	if err := rows.Scan(&title); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(title)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	// return response, err
	return servers, nil
}

// InsertData function
func (db *DB) InsertData(temp string, humed string) error {
	_, err := db.Exec(`
		INSERT INTO data (temperatura, humedad)  VALUES (
			$1, $2, $3
		)
	`, temp, humed)
	if err != nil {
		log.Println(err)
	}

	return nil
}
