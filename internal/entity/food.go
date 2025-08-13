package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

type Food struct {
	Id           uuid.UUID  `json:"id"`
	RestaurantId string     `json:"restaurant_id"`
	CategoryId   uuid.UUID  `json:"category_id,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Price        float64    `json:"price"`
	Images       string     `json:"images"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (Food) TableName() string {
	return "foods"
}

type FoodUpdateReq struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Status       *string `json:"status,omitempty"`
	RestaurantId *string `json:"restaurant_id,omitempty"`
	CategoryId   *string `json:"category_id,omitempty"`
	Image        *string `json:"image,omitempty"`

	Id uuid.UUID `json:"-"`
}

func (FoodUpdateReq) TableName() string {
	return Food{}.TableName()
}

func (c FoodUpdateReq) Validate() error {
	if c.Id == uuid.Nil {
		return ErrFoodInvalidID
	}

	if c.Name == "" {
		return errors.New("name cannot be empty")
	}

	if c.Status != nil {
		status := *c.Status
		validStatuses := map[string]bool{
			string(StatusActive):   true,
			string(StatusInactive): true,
			string(StatusDeleted):  true,
		}
		if !validStatuses[status] {
			return fmt.Errorf("%w: %s", ErrFoodInvalidStatus, status)
		}
	}

	if c.RestaurantId != nil {
		if *c.RestaurantId == "" {
			return ErrFoodInvalidRestaurant
		}
	}

	if c.CategoryId != nil {
		if *c.CategoryId == "" {
			return ErrFoodInvalidCategory
		}
	}

	return nil
}
