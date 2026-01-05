package database

import (
	"context"
	"encoding/json"
	"time"
)

func SetCache(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RDB.Set(context.Background(), key, data, expiration).Err()
}

func GetCache(key string, dest interface{}) error {
	data, err := RDB.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func DeleteCache(key string) error {
	return RDB.Del(context.Background(), key).Err()
}

func DeleteCachePattern(pattern string) error {
	ctx := context.Background()
	iter := RDB.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := RDB.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
