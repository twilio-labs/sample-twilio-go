package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	// import postgres driver pq
	_ "github.com/lib/pq"
	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

// initializeDB initializes the database connection
func InitializeDB() (*DB, error) {
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// DB is a wrapper around the database/sql.DB type
type DB struct {
	*sql.DB
}

// NewDB creates a new DB instance
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetCustomerByID returns a customer by id
func (db *DB) GetCustomerByID(ctx context.Context, id string) (*configuration.Customer, error) {
	customer := &configuration.Customer{}
	err := db.QueryRowContext(ctx, "SELECT * FROM customers WHERE id = $1", id).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.PhoneNumber, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// GetCustomerByPhoneNumber
func GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*configuration.Customer, error) {
	db, err := InitializeDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	customer, err := db.GetCustomerByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// GetCustomerByEmail returns a customer by email
func (db *DB) GetCustomerByEmail(ctx context.Context, email string) (*configuration.Customer, error) {
	customer := &configuration.Customer{}
	err := db.QueryRowContext(ctx, "SELECT * FROM customers WHERE email = $1", email).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.PhoneNumber, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// InitalizeDB and create a new customer
func CreateNewCustomer(ctx context.Context, customer *configuration.Customer) error {
	db, err := InitializeDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.CreateCustomer(ctx, customer)
	if err != nil {
		return err
	}
	return nil
}

// CreateCustomer creates a new customer
func (db *DB) CreateCustomer(ctx context.Context, customer *configuration.Customer) error {
	// Add Validation
	// Check if the Customer already exists in the databse
	check, err := db.QueryContext(ctx, "SELECT * FROM customers WHERE email = $1 AND phone_number = $2", customer.Email, customer.PhoneNumber)
	if err != nil {
		return err
	}
	fmt.Print(check)
	// If the customer already exists, return an error
	defer check.Close()
	if check.Next() {
		return fmt.Errorf("Customer already exists")
	}
	_, err = db.ExecContext(ctx, "INSERT INTO customers (id, first_name, last_name, phone_number, email, created_at) VALUES ($1, $2, $3, $4, $5, $6)", customer.ID, customer.FirstName, customer.LastName, customer.PhoneNumber, customer.Email, customer.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetCustomers returns all customers
func (db *DB) GetCustomers(ctx context.Context) ([]*configuration.Customer, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []*configuration.Customer{}
	for rows.Next() {
		customer := &configuration.Customer{}
		err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.CreatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

// GetCustomerByPhoneNumber returns a customer by phone number
func (db *DB) GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*configuration.Customer, error) {
	customer := &configuration.Customer{}
	err := db.QueryRowContext(ctx, "SELECT * FROM customers WHERE phone_number = $1", phoneNumber).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.PhoneNumber, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
