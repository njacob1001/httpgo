package databases

// MongoInterface interface
type MongoInterface interface {
	GetDataByID(articuloID string) (*Product, error)
	CalculateSum(articles []string) (int64, error)
	// InsertData(temp string, humed string) error
}

// PostgresInterface interface
type PostgresInterface interface {
	CreateUser(username string, password string, cash int64) error
	UpdateUserCash(username string, cash int64) error
	AddCash(username string, cash int64) error
	FindUserByID(username string, password string) (*ClientResponse, error)
	GetUser(username string) (*ClientResponse, error)
}
