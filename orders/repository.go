package orders

type Repository interface {
	Store(orders []Order) error
	StoreCustomerGroups(groups []CustomerGroup) error
	Disconnected() error
}
