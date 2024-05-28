package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
	"yadro/internal/core/comics"
)

/*func TestHead(t *testing.T) {
	// Create a temporary file for the SQLite database
	tempFile, err := os.CreateTemp("", "testdb_*.sqlite")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Mock the SQL operations
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Replace the sql.Open with the mock DB
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	// Mock the migrate instance
	migrateMock := &MockMigrate{}
	migrateNewWithDatabaseInstance = func(sourceURL string, databaseName string, instance database.Driver) (*migrate.Migrate, error) {
		return migrateMock, nil
	}

	// Set expectations
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS comics").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO comics").WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the Head function
	comicsMap := map[int]comics.Write{
		1: {Tscript: []string{"tst1", "tst1"}, Img: "c1"},
		2: {Tscript: []string{"tst2", "tst2"}, Img: "c2"},
	}
	indexMap := map[string][]string{
		"index1": {"Comic 1", "Comic 2"},
	}

	repository.Head(tempFile.Name(), comicsMap, indexMap)

	// Ensure all expectations are met
	require.NoError(t, mock.ExpectationsWereMet())
}

type MockMigrate struct{}

func (m *MockMigrate) Up() error {
	return nil
}

func (m *MockMigrate) Close() error {
	return nil
}

// Mock functions to override actual implementations
var (
	sqlOpen                        = sql.Open
	migrateNewWithDatabaseInstance = migrate.NewWithDatabaseInstance
)*/

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	comicsMap := map[int]comics.Write{
		1: {Tscript: []string{"keyword1", "keyword2"}, Img: "http://example.com/image1.jpg"},
		2: {Tscript: []string{"keyword3", "keyword4"}, Img: "http://example.com/image2.jpg"},
	}

	indexMap := map[string][]string{
		"keyword1": {"1"},
		"keyword2": {"1", "2"},
	}

	// Expected queries and their parameters
	mock.ExpectExec(`INSERT INTO database\(id, Keywords, Url\) VALUES \(\?, \?, \?\)`).
		WithArgs(1, "keyword1,keyword2", "http://example.com/image1.jpg").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO database\(id, Keywords, Url\) VALUES \(\?, \?, \?\)`).
		WithArgs(2, "keyword3,keyword4", "http://example.com/image2.jpg").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO index_table\(Keywords, Numbers\) VALUES \(\?, \?\)`).
		WithArgs("keyword1", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO index_table\(Keywords, Numbers\) VALUES \(\?, \?\)`).
		WithArgs("keyword2", "1,2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Add(db, comicsMap, indexMap)

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchRecords(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	comicsData := []struct {
		id       int
		keywords string
		url      string
	}{
		{1, "keyword1,keyword2", "http://example.com/image1.jpg"},
		{2, "keyword3,keyword4", "http://example.com/image2.jpg"},
	}

	indexData := []struct {
		keyword string
		numbers string
	}{
		{"keyword1", "1,2"},
		{"keyword2", "3"},
	}

	// Expect database query for comics
	mock.ExpectQuery("SELECT \\* FROM database").
		WillReturnRows(sqlmock.NewRows([]string{"id", "keywords", "url"}).
			AddRow(comicsData[0].id, comicsData[0].keywords, comicsData[0].url).
			AddRow(comicsData[1].id, comicsData[1].keywords, comicsData[1].url))

	// Expect database query for index
	mock.ExpectQuery("SELECT \\* FROM index_table").
		WillReturnRows(sqlmock.NewRows([]string{"keyword", "numbers"}).
			AddRow(indexData[0].keyword, indexData[0].numbers).
			AddRow(indexData[1].keyword, indexData[1].numbers))

	comicsMap, indexMap := FetchRecords(db)

	// Check if fetched comics map matches the expected map
	expectedComicsMap := map[int]comics.Write{
		1: {Tscript: []string{"keyword1", "keyword2"}, Img: "http://example.com/image1.jpg"},
		2: {Tscript: []string{"keyword3", "keyword4"}, Img: "http://example.com/image2.jpg"},
	}

	if !reflect.DeepEqual(comicsMap, expectedComicsMap) {
		t.Errorf("fetched comics map doesn't match expected: got %v, want %v", comicsMap, expectedComicsMap)
	}

	// Check if fetched index map matches the expected map
	expectedIndexMap := map[string][]int{
		"keyword1": {1, 2},
		"keyword2": {3},
	}
	if !reflect.DeepEqual(indexMap, expectedIndexMap) {
		t.Errorf("fetched index map doesn't match expected: got %v, want %v", indexMap, expectedIndexMap)
	}

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
