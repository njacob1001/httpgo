package databases

import (
	"fmt"
	"log"
)

// CreateUser function
func (psql *DB) CreateUser(username string, password string, cash int64) error {
	_, err := psql.Exec(`
		INSERT INTO clients (username, password, cash)  VALUES (
			$1, $2, $3
		)
	`, username, password, cash)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// UpdateUserCash function
func (psql *DB) UpdateUserCash(username string, cash int64) error {
	_, err := psql.Exec("UPDATE clients SET cash=$2 WHERE username=$1", username, cash)
	if err != nil {
		return err
	}
	return nil
}

// AddCash function
func (psql *DB) AddCash(username string, cash int64) error {
	query := fmt.Sprintf("SELECT username, cash FROM domains WHERE username='%s'", username)
	row := psql.QueryRow(query)
	client := new(Client)
	err := row.Scan(&client.Username, &client.Cash)
	if err != nil {
		return err
	}
	_, err2 := psql.Exec("UPDATE clients SET cash=$1 WHERE username=$2", client.Cash+cash, username)
	if err2 != nil {
		return err
	}
	return nil
}

// FindUserByID func
func (psql *DB) FindUserByID(username string, password string) (*ClientResponse, error) {
	query := fmt.Sprintf("SELECT username, cash FROM clients WHERE username='%s' AND password='%s'", username, password)
	row := psql.QueryRow(query)

	client := new(ClientResponse)
	err := row.Scan(&client.Username, &client.Cash)
	if err != nil {
		return client, err
	}
	return client, nil
}

// GetUser func
func (psql *DB) GetUser(username string) (*ClientResponse, error) {
	query := fmt.Sprintf("SELECT username, cash FROM clients WHERE username='%s'", username)
	row := psql.QueryRow(query)

	client := new(ClientResponse)
	err := row.Scan(&client.Username, &client.Cash)
	if err != nil {
		return client, err
	}
	return client, nil
}
