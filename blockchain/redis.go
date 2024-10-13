package blockchain

import (
	"github.com/go-redis/redis"
)

// RedisStorage can only be used for debug,
// because Redis do not support ACID
type RedisStorage struct {
	client *redis.Client
}

var _ Storage = &RedisStorage{}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (r *RedisStorage) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisStorage) GetBlock(key string) (*Block, error) {
	serializedBlock, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	return str2Block(serializedBlock)
}

func (r *RedisStorage) Set(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisStorage) SetBlock(key string, value *Block) error {
	serializedBlock, err := block2Str(value)
	if err != nil {
		return err
	}
	return r.client.Set(key, serializedBlock, 0).Err()
}

func (r *RedisStorage) Delete(key string) error {
	return r.client.Del(key).Err()
}

func (r *RedisStorage) KeyNotFound(err error) bool {
	return err.Error() == redis.Nil.Error()
}
