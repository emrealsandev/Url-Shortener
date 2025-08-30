package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// İleride kullanacağımız URL koleksiyonu indexleri:
// - uniq_code: code unique
// - uniq_custom_alias: custom_alias varsa unique (partial)
// - ttl_expire: expires_at alanına TTL
func EnsureIndexes(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_code"),
	})
	if err != nil {
		return err
	}

	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "custom_alias", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("uniq_custom_alias").
			SetPartialFilterExpression(bson.M{"custom_alias": bson.M{"$type": "string"}}),
	})
	if err != nil {
		return err
	}

	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_expire"),
	})
	return err
}
