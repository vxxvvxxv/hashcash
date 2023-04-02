package db

type DB interface {
	GetDataFromDB(id int) (string, bool)
	GetRandomDataFromDB() string
	// FillTestData fills DB with test data from local file
	FillTestData() error
}
