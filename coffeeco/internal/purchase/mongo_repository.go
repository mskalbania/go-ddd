package purchase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoRepository struct {
	purchases *mongo.Collection
}

func NewMongoRepository(ctx context.Context, conn string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}
	purchases := client.Database("coffeeco").Collection("purchases")
	return &MongoRepository{purchases: purchases}, nil
}

func (m *MongoRepository) Store(ctx context.Context, purchase Purchase) error {
	mongoP := toMongoPurchase(purchase)
	if _, err := m.purchases.InsertOne(ctx, mongoP); err != nil {
		return fmt.Errorf("failed to store purchase: %w", err)
	}
	return nil
}

func (m *MongoRepository) Ping(ctx context.Context) error {
	if _, err := m.purchases.EstimatedDocumentCount(ctx); err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

type mongoPurchase struct {
	ID             uuid.UUID      `bson:"id"`
	Store          mongoStore     `bson:"store"`
	Products       []mongoProduct `bson:"products_purchased"`
	Total          int64          `bson:"purchase_total"`
	PaymentMeans   string         `bson:"payment_means"`
	TimeOfPurchase time.Time      `bson:"created_at"`
	CardToken      *string        `bson:"card_token"`
}

type mongoProduct struct {
	ItemName  string `bson:"item_name"`
	BasePrice int64  `bson:"base_price"`
}

type mongoStore struct {
	ID       uuid.UUID `bson:"id"`
	Location string    `bson:"location"`
}

func toMongoPurchase(purchase Purchase) mongoPurchase {
	var mongoProducts []mongoProduct
	for _, p := range purchase.Products {
		mongoProducts = append(mongoProducts, mongoProduct{
			ItemName:  p.ItemName,
			BasePrice: p.BasePrice.Amount(),
		})
	}
	return mongoPurchase{
		ID: purchase.id,
		Store: mongoStore{
			ID:       purchase.Store.ID,
			Location: purchase.Store.Location,
		},
		Products:       mongoProducts,
		Total:          purchase.total.Amount(),
		PaymentMeans:   string(purchase.PaymentMeans),
		TimeOfPurchase: purchase.timeOfPurchase,
		CardToken:      purchase.CardToken,
	}
}
