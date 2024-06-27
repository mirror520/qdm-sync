package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/mirror520/qdm-sync/orders"

	sync "github.com/mirror520/qdm-sync"
)

type OrderRepository interface {
	orders.Repository
	DropDatabase(name string) error
}

type orderRepository struct {
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOrderRepository(cfg sync.Persistence) (orders.Repository, error) {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &orderRepository{
		ctx:    ctx,
		cancel: cancel,
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Address))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database)
	if err := db.CreateCollection(ctx, "orders"); err != nil {
		cmdErr, ok := err.(mongo.CommandError)
		if !ok || cmdErr.Code != 48 {
			return nil, err
		}
	}

	repo.db = db

	return repo, nil
}

func (repo *orderRepository) Store(orders []orders.Order) error {
	coll := repo.db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	docs := make([]any, len(orders))
	for i, o := range orders {
		docs[i] = bson.M(o)
	}

	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *orderRepository) StoreCustomerGroups(groups []*orders.CustomerGroup) error {
	coll := repo.db.Collection("customer_groups")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	docs := make([]any, len(groups))
	for i, g := range groups {
		docs[i] = g
	}

	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *orderRepository) Disconnected() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return repo.db.Client().Disconnect(ctx)
}
