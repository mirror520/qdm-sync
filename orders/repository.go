package orders

type Repository interface {
	Store(orders []Order) error
	StoreCustomers(customers []Customer) error
	StoreCustomerGroups(groups []CustomerGroup) error
	Disconnected() error
}
