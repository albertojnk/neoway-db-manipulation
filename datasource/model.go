package datasource

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ClientInfo represents a table in the database that is related to a client
type ClientInfo struct {
	gorm.Model
	CPF                string
	Private            bool
	Incomplete         bool
	LastPurchaseDate   time.Time
	AverageBudget      float64
	LastPurchaseBudget float64
	MostFrequentStore  string
	LastPurchaseStore  string
}
