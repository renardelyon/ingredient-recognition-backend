package domain

import "fmt"

type Ingredient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// Validate checks if the ingredient has valid data.
func (i *Ingredient) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("ingredient name cannot be empty")
	}
	if i.Quantity <= 0 {
		return fmt.Errorf("ingredient quantity must be greater than zero")
	}
	if i.Unit == "" {
		return fmt.Errorf("ingredient unit cannot be empty")
	}
	return nil
}
