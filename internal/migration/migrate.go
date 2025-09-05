package migration

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UrlsColl     = "urls"
	SequenceColl = "sequence"
	IdxCodeV1    = "uniq_code_v1"
	IdxAliasV1   = "uniq_custom_alias_v1"
	IdxExpireV1  = "ttl_expire_v1"
)

type Migrator struct {
	DB *mongo.Database
	// buraya istersek logger, options vs. ekleyebiliriz
}

func New(db *mongo.Database) *Migrator {
	return &Migrator{DB: db}
}

func (m *Migrator) RunAll(ctx context.Context) error {
	if err := m.ensureCollection(ctx, UrlsColl); err != nil {
		return fmt.Errorf("ensure collection urls: %w", err)
	}
	urlCollection := m.DB.Collection(UrlsColl)

	if err := m.ensureValidator(ctx, urlCollection); err != nil {
		return fmt.Errorf("ensure validator: %w", err)
	}
	if err := m.ensureIndexes(ctx, urlCollection); err != nil {
		return fmt.Errorf("ensure indexes: %w", err)
	}

	if err := m.ensureCollection(ctx, SequenceColl); err != nil {
		return fmt.Errorf("ensure collection sequence: %w", err)
	}

	sequenceCollection := m.DB.Collection(SequenceColl)

	if err := m.ensureValidator(ctx, sequenceCollection); err != nil {
		return fmt.Errorf("ensure validator: %w", err)
	}

	return nil
}

func (m *Migrator) ensureCollection(ctx context.Context, name string) error {
	names, err := m.DB.ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if len(names) > 0 {
		return nil
	}
	return m.DB.CreateCollection(ctx, name)
}

func (m *Migrator) ensureValidator(ctx context.Context, coll *mongo.Collection) error {
	schema := bson.M{}
	if coll.Name() == UrlsColl {
		schema = bson.M{
			"bsonType":             "object",
			"required":             bson.A{"code", "target", "created_at", "disabled"},
			"additionalProperties": false,
			"properties": bson.M{
				"_id":          bson.M{"bsonType": "objectId"},
				"code":         bson.M{"bsonType": "string"},
				"target":       bson.M{"bsonType": "string"},
				"created_at":   bson.M{"bsonType": "date"},
				"disabled":     bson.M{"bsonType": "bool"},
				"expires_at":   bson.M{"bsonType": bson.A{"date", "null"}},
				"custom_alias": bson.M{"bsonType": "string"},
			},
		}

	} else if coll.Name() == SequenceColl {
		schema = bson.M{
			"bsonType":             "object",
			"required":             bson.A{"seq"},
			"additionalProperties": false,
			"properties": bson.M{
				"_id": bson.M{"bsonType": "string"},
				"seq": bson.M{"bsonType": bson.A{"long", "int"}},
			},
		}
	} else {
		return fmt.Errorf("collection %s not found", coll.Name())
	}

	cmd := bson.D{
		{Key: "collMod", Value: coll.Name()},
		{Key: "validator", Value: bson.M{"$jsonSchema": schema}},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return coll.Database().RunCommand(ctx, cmd).Err()
}

func (m *Migrator) ensureIndexes(ctx context.Context, coll *mongo.Collection) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetName(IdxCodeV1).SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "custom_alias", Value: 1}},
			Options: options.Index().
				SetName(IdxAliasV1).
				SetUnique(true).
				SetPartialFilterExpression(bson.M{"custom_alias": bson.M{"$type": "string"}}),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetName(IdxExpireV1).SetExpireAfterSeconds(0),
		},
	}

	existing := map[string]struct{}{}
	cur, err := coll.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var x bson.M
		if err := cur.Decode(&x); err != nil {
			return err
		}
		if name, _ := x["name"].(string); name != "" {
			existing[name] = struct{}{}
		}
	}

	for _, idx := range indexes {
		name := *idx.Options.Name
		if _, ok := existing[name]; ok {
			continue // zaten var
		}
		if _, err := coll.Indexes().CreateOne(ctx, idx); err != nil {
			return fmt.Errorf("create index %s: %w", name, err)
		}
	}
	return nil
}
