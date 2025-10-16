package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/emrealsandev/Url-Shortener/internal/repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLRepo struct {
	urlCollection      *mongo.Collection
	seqCollection      *mongo.Collection
	settingsCollection *mongo.Collection
}

func NewURLRepo(db *mongo.Database) *URLRepo {
	return &URLRepo{
		urlCollection:      db.Collection(repo.COLLECTION_URLS),
		seqCollection:      db.Collection(repo.COLLECTION_SEQUENCE),
		settingsCollection: db.Collection(repo.COLLECTION_SETTINGS),
	}
}

func (r *URLRepo) Insert(u repo.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.urlCollection.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("duplicate") // service ErrConflict'a çeviriyor
	}
	return err
}

func (r *URLRepo) GetByCode(code string) (*repo.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var out repo.URL
	err := r.urlCollection.FindOne(ctx, bson.M{"code": code}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &out, err
}

func (r *URLRepo) FindOneAndUpdate(ctx context.Context) (uint64, error) {
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After).
		SetHint(bson.D{{Key: "_id", Value: 1}}) // opsiyonel

	var out struct {
		Seq int64 `bson:"seq"`
	}

	err := r.seqCollection.FindOneAndUpdate(
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

func (r *URLRepo) GetCodeByUrl(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var out repo.URL
	err := r.urlCollection.FindOne(ctx, bson.M{"target": url}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	if out.Code == "" {
		return "", nil
	}

	return out.Code, nil
}

func (r *URLRepo) GetAllSettings() (*repo.Settings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var settings repo.Settings
	// TEK BİR sorgu ile tüm ayarları al.
	err := r.settingsCollection.FindOne(ctx, bson.M{}).Decode(&settings)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Hiç ayar yoksa boş struct dön, bu bir hata değil.
			return &repo.Settings{}, nil
		}
		// Gerçek bir hata varsa onu dön.
		return nil, err
	}
	return &settings, nil
}
