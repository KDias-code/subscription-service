package model

import (
	"github.com/google/uuid"
	"time"
)

type Subscriptions struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ServiceName string    `json:"serviceName" db:"service_name"`
	Price       uint64    `json:"price" db:"price"`
	UserID      uuid.UUID `json:"userID" db:"user_id"`
	StartDate   string    `json:"startDate" db:"start_date"`
	EndDate     string    `json:"endDate" db:"end_date"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type SubscriptionsSum struct {
	ID          *uuid.UUID `json:"id" db:"id"`
	ServiceName *string    `json:"serviceName" db:"service_name"`
	StartDate   string     `json:"startDate" db:"start_date"`
	EndDate     string     `json:"endDate" json:"end_date"`
}
