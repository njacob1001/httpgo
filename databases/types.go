package databases

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

// DB database
type DB struct {
	*sql.DB
}

// MDB struct
type MDB struct {
	*mongo.Client
}

// Product struct
type Product struct {
	ID     string `json:"id"`
	Precio int64  `json:"precio"`
	Name   string `json:"name"`
}

// Client struct
type Client struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Cash     int64  `json:"cash"`
}

// ClientResponse struct
type ClientResponse struct {
	Username string `json:"username"`
	Cash     int64  `json:"cash"`
}

// ProductAPIResponse struct
type ProductAPIResponse struct {
	Ok      bool      `json:"ok"`
	Cessage string    `json:"message"`
	Data    []Product `json:"data"`
}
