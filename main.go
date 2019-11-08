package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "worker"
	dbname   = "redesuao"
)

// Response respuesta de la api
type Response struct {
	Temperatura float32 `json:"temperatura"`
	Humedad     float32 `json:"humedad"`
}

// TotalResponse respuesta de la api
type TotalResponse struct {
	Temperatura float32 `json:"temperatura"`
	Humedad     float32 `json:"humedad"`
	Resultado   string  `json:"resultado"`
}

var global Response

func handleGet(w http.ResponseWriter, r *http.Request) {

	js, err := json.Marshal(global)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&global); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(global)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	temp := "frio"
	humedad := "seco"

	if global.Temperatura >= 26 {
		temp = "caluroso"
	}
	if global.Humedad >= 50 {
		humedad = "humedo"
	}

	fullResponse := TotalResponse{
		Temperatura: global.Temperatura,
		Humedad:     global.Humedad,
		Resultado:   temp + " y " + humedad,
	}

	js, err := json.Marshal(fullResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

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

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(initializeTable); err != nil {
		panic(err)
	}

	if _, err := db.Exec(defaultValue); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	r := chi.NewRouter()
	s := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	r.Get("/restserver/app/ws", handleGet)
	r.Post("/restserver/app/ws", handlePost)
	r.Get("/restserver/app/weather", handleWeather)
	s.ListenAndServe()
}
