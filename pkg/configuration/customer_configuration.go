package configuration

import ksuid "github.com/segmentio/ksuid"

type Customer struct {
	PhoneNumber string
	Name   string
	Email       string
	ID 		string
	CreatedAt   string
}

func CustomerToMap(customer *Customer) map[string]interface{} {
	return map[string]interface{}{
		"phoneNumber": customer.PhoneNumber,
		"name": customer.Name,
		"email": customer.Email,
		"id": customer.ID,
		"createdAt": customer.CreatedAt,
	}
}

// kuid is the unique identifier for a customer
type Kuid string

// create function to generate a kuid string for the id of the customer
func GenerateUserID() string {
	return (ksuid.New().String())
}
