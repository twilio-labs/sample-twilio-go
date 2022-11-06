package db

import (
	"context"
	"testing"

	// import postgres driver pq
	"github.com/bmizerany/assert"
	_ "github.com/lib/pq"
)

// stub a postgres database
//
// docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
// create mock data in postgres
//
// docker exec -it postgres psql -U postgres
// create database postgres;
// \c postgres
// create table customers (id varchar(255), name varchar(255), email varchar(255), phone_number varchar(255), created_at varchar(255));
// insert into customers values ('1', 'John Doe', '+1555555555', '2019-01-01 00:00:00');

// test functions in db.go
func TestInitializeDB(t *testing.T) {
	// Act
	result, error := InitializeDB()
	if error != nil {
		t.Error(error)
	}

	// Assert
	assert.Equal(t, true, result != nil)
}

func TestNewDB(t *testing.T) {
	// Arrange
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	// Act
	result, error := NewDB(dsn)
	if error != nil {
		t.Error(error)
	}

	// Assert
	assert.Equal(t, true, result != nil)
}

func TestClose(t *testing.T) {
	// Arrange
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, error := NewDB(dsn)
	if error != nil {
		t.Error(error)
	}

	// Act
	result := db.Close()

	// Assert
	assert.Equal(t, nil, result)
}

func TestGetCustomerByID(t *testing.T) {
	// Arrange
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, error := NewDB(dsn)
	if error != nil {
		t.Error(error)
	}
	id := "1"

	// Act
	result, error := db.GetCustomerByID(context.Background(), id)
	if error != nil {
		t.Error(error)
	}

	// Assert
	assert.Equal(t, true, result != nil)
}

func TestGetCustomerByPhoneNumber(t *testing.T) {
	// Arrange
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, error := NewDB(dsn)
	if error != nil {
		t.Error(error)
	}
	phoneNumber := "+15555555555"

	// Act
	// mock context
	ctx := context.Background()
	result, error := db.GetCustomerByPhoneNumber(ctx,phoneNumber)
	if error != nil {
		t.Error(error)
	}

	// Assert
	assert.Equal(t, true, result != nil)
}
