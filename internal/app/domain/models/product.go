package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string
	Price       float64   `gorm:"not null"`
	Stock       int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (p *Product) NewProduct(name, description string, price float64, stock int) {
	p.ID = uuid.New()
	p.Description = description
	p.Name = name
	p.Price = price
	p.Stock = stock
}

func (p *Product) UpdateProduct(name, description *string, price *float64, stock *int) {
	p.Description = func() string {
		if (description != nil && *description == "") || description == nil {
			return p.Description
		}
		return *description
	}()
	p.Name = func() string {
		if (name != nil && *name == "") || name == nil {
			return p.Name
		}
		return *name
	}()
	p.Price = func() float64 {
		if (price != nil && *price == 0) || price == nil {
			return p.Price
		}
		return *price
	}()
	p.Stock = func() int {
		if (stock != nil && *stock == 0) || stock == nil {
			return p.Stock
		}
		return *stock
	}()
}
