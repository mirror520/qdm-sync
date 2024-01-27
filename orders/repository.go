package orders

type Repository interface {
	Store(orders []Order) error
	Disconnected() error
}
