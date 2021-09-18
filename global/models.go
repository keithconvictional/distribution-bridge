package global

import "time"

const (
	DomainOrders = "orders"
	DomainProducts = "products"
	DomainGeneral = "general"
)

type State struct {
	LastInventorySync time.Time
	LastProductSync time.Time
	LastOrderSync time.Time
}
