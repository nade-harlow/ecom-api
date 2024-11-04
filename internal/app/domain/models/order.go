package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

const (
	OrderStatus_Pending   OrderStatus = "Pending"
	OrderStatus_Canceled  OrderStatus = "Canceled"
	OrderStatus_Fulfilled OrderStatus = "Fulfilled"
)

type Order struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID   `gorm:"type:uuid;index;not null" json:"-"`
	User      User        `gorm:"foreignKey:UserID" json:"-"`
	Status    OrderStatus `gorm:"type:varchar(50);default:Pending;not null" json:"status"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	Items     []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;index;not null" json:"-"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"-"`
	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"not null" json:"price"`
}

func (o *Order) NewOrder() {
	o.ID = uuid.New()
	o.Status = OrderStatus_Pending
	o.Items = func() []OrderItem {
		for i, item := range o.Items {
			o.Items[i].ID = uuid.New()
			o.Items[i].OrderID = o.ID
			o.Items[i].ProductID = item.ProductID
			o.Items[i].Quantity = item.Quantity
			o.Items[i].Price = item.Price
		}
		return o.Items
	}()
}
