package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rbcervilla/redisstore/v8"
	"strconv"

	"github.com/changyenh/sessions"
)

type Store interface {
	sessions.Store
}

// size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(size int, address, password string) (Store, error) {
	ctx := context.Background()
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              []string{address},
		DB:                 0,
		Password:           password,
		PoolSize:           size,
	})
	s, err := redisstore.NewRedisStore(ctx, rdb)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithDB - like NewStore but accepts `DB` parameter to select
// redis DB instead of using the default one ("0")
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithDB
func NewStoreWithDB(size int, address, password, DB string) (Store, error) {
	ctx := context.Background()
	db, err := strconv.Atoi(DB)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              []string{address},
		DB:                 db,
		Password:           password,
		PoolSize:           size,
	})
	s, err := redisstore.NewRedisStore(ctx, rdb)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithPool instantiates a RediStore with a *redis.Pool passed in.
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithPool
func NewStoreWithRedis(redis *redis.UniversalClient) (Store, error) {
	ctx := context.Background()
	s, err := redisstore.NewRedisStore(ctx, *redis)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

type store struct {
	*redisstore.RedisStore
}

// GetRedisStore get the actual woking store.
// Ref: https://godoc.org/github.com/boj/redistore#RediStore
func GetRedisStore(s Store) (err error, redisStore *redisstore.RedisStore) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return
	}

	redisStore = realStore.RedisStore
	return
}

// SetKeyPrefix sets the key prefix in the redis database.
func SetKeyPrefix(s Store, prefix string) error {
	err, redisStore := GetRedisStore(s)
	if err != nil {
		return err
	}

	redisStore.KeyPrefix(prefix)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RedisStore.Options(*options.ToGorillaOptions())
}
