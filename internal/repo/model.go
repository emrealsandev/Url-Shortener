package repo

import (
	"time"
)

const COLLECTION_URLS = "urls"
const COLLECTION_SETTINGS = "settings"
const COLLECTION_SEQUENCE = "sequence"

type URL struct {
	Code        string     `bson:"code" json:"code"`
	Target      string     `bson:"target" json:"target"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	ExpiresAt   *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	Disabled    bool       `bson:"disabled" json:"disabled"`
	CustomAlias *string    `bson:"custom_alias,omitempty" json:"custom_alias,omitempty"`
	OwnerID     *int64     `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
}

type Settings struct {
	TtlTime      int8 `bson:"ttl_time" json:"ttl_time"`
	RedisTtlTime int8 `bson:"redis_ttl" json:"redis_ttl"`
}
