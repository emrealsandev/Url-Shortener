package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"url-shortener/internal/short"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepo struct {
	coll *mongo.Collection
}

func NewURLRepo(coll *mongo.Collection) *URLRepo {
	return &URLRepo{coll: coll}
}

func (r *URLRepo) Insert(u short.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.coll.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("duplicate") // service ErrConflict'a çeviriyor
	}
	return err
}

func (r *URLRepo) GetByCode(code string) (*short.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var out short.URL
	err := r.coll.FindOne(ctx, bson.M{"code": code}).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &out, err
}

func (r *URLRepo) FindOneAndUpdate(ctx context.Context) (uint64, error) {
	// upsert + returnDocument:After ⇒ atomik artır, yeni değeri döndür
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After).
		SetHint(bson.D{{Key: "_id", Value: 1}}) // opsiyonel

	var out struct {
		Seq int64 `bson:"seq"`
	}

	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "url"},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&out)
	if err != nil {
		return 0, err
	}
	return uint64(out.Seq), nil
}
