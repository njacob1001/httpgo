package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/njacob1001/server/databases"
	"github.com/njacob1001/server/routes"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "rysm"
	password = "worker"
	dbname   = "uaostore"
)

// Env estruct
type Env struct {
	Mongo databases.MongoInterface
	Psql  databases.PostgresInterface
}

func (env *Env) loginUser(w http.ResponseWriter, r *http.Request) {
	var Body routes.UserVerification
	if err := json.NewDecoder(r.Body).Decode(&Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	client, err := env.Psql.FindUserByID(Body.Username, Body.Password)

	if err != nil {
		resp := routes.UserAuthResponse{
			Ok:      false,
			Message: "Usuario o contraseña incorrectos",
		}
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	js, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func (env *Env) createUser(w http.ResponseWriter, r *http.Request) {
	var Body databases.Client
	if err := json.NewDecoder(r.Body).Decode(&Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err := env.Psql.CreateUser(Body.Username, Body.Password, Body.Cash)
	if err != nil {
		resp := routes.UserAuthResponse{
			Ok:      false,
			Message: "No se pudo crear el usuario",
		}
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	resp := routes.UserAuthResponse{
		Ok: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

func main() {
	pdb, err := databases.NewDB("postgresql://rysm@localhost:5432/uaostore?password=worker&sslmode=disable")
	if err != nil {
		panic(err)
	}

	db, err := databases.NewMongoClient()
	if err != nil {
		panic(err)
	}

	env := &Env{Mongo: db, Psql: pdb}
	// defer db.Close()

	r := chi.NewRouter()
	s := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	r.Post("/api/user/login", env.loginUser)
	r.Post("/api/user/create", env.createUser)
	// r.Get("/api/get", env.handleGet)
	// r.Get("/api/datos/{type}", env.handleGet)
	// r.Post("/api/insert", env.handlePost)
	// r.Post("/restserver/app/ws", env.handlePost)
	// r.Get("/restserver/app/weather", databases.env.algo)
	fmt.Println("Successfully connected!")
	s.ListenAndServe()
}
