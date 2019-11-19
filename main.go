package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/njacob1001/httpgo/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "worker"
	dbname   = "redesuao"
)

// Env estruct
type Env struct {
	db models.Datastore
}

// RespOk struct
type Respok struct {
	Ok bool `json:"ok"`
}

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

func (env *Env) handleGet(w http.ResponseWriter, r *http.Request) {
	searchType := chi.URLParam(r, "type")
	var resp []*models.Basura
	var err error
	if searchType == "urgent" {
		resp, err = env.db.GetDataBy(true)
	} else {
		resp, err = env.db.GetDataBy(false)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// js, err := json.Marshal(global)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func (env *Env) handlePost(w http.ResponseWriter, r *http.Request) {
	var Body models.BasuraFromArduino
	if err := json.NewDecoder(r.Body).Decode(&Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err := env.db.InsertData(&Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	respuesta := Respok{
		Ok: true,
	}

	js, err := json.Marshal(respuesta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// errr := env.db.InsertData(fmt.Sprintf("%f", global.Temperatura), fmt.Sprintf("%f", global.Humedad))
	// if errr != nil {
	// 	http.Error(w, errr.Error(), http.StatusInternalServerError)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func (env *Env) handleWeather(w http.ResponseWriter, r *http.Request) {
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
	// 	db, err := models.NewDB("postgresql://postgres@localhost:5432/redesuao?password=worker")
	db, err := models.NewMongoClient()
	if err != nil {
		panic(err)
	}

	env := &Env{db: db}
	// defer db.Close()

	r := chi.NewRouter()
	s := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// r.Get("/api/get", env.handleGet)
	r.Get("/api/datos/{type}", env.handleGet)
	r.Post("/api/insert", env.handlePost)
	// r.Post("/restserver/app/ws", env.handlePost)
	// r.Get("/restserver/app/weather", env.handleWeather)
	s.ListenAndServe()
	fmt.Println("Successfully connected!")
}
