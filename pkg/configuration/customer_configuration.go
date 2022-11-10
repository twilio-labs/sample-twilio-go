package configuration

import "github.com/google/uuid"

type Customer struct {
	PhoneNumber string `json:"email`
	FirstName   string `json:"firstName`
	LastName    string `json:"lastName`
	Email       string `json:"email`
	ID          string `json:"id`
	CreatedAt   string `json:"createdAt`
}

func CustomerToMap(customer *Customer) map[string]interface{} {
	return map[string]interface{}{
		"phoneNumber": customer.PhoneNumber,
		"firstName":   customer.FirstName,
		"lastName":    customer.LastName,
		"email":       customer.Email,
		"id":          customer.ID,
		"createdAt":   customer.CreatedAt,
	}
}

// create function to generate a kuid string for the id of the customer
func GenerateUserID() string {
	return uuid.New().String()
}
