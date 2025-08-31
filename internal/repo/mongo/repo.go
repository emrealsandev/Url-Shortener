package mongo

import (
	"context"
	"errors"
	"time"

	"url-shortener/internal/short"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepo struct {
	coll *mongo.Collection
}

func NewURLRepo(coll *mongo.Collection) *URLRepo { return &URLRepo{coll: coll} }

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
