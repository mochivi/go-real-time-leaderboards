package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/mochivi/go-real-time-leaderboards/config"
)

var ErrNotFound = errors.New("key not found")

type RedisService interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string, any) error
	JSONSet(context.Context, string, string, any, time.Duration) error
	JSONGet(context.Context, string, string, any) error
}

// redisService is the concrete redis implementation
type redisService struct {
	client *redis.Client
}

func NewRedisService(redisConfig config.RedisConfig) *redisService {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr(),
		Password: redisConfig.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("redis connect err: %v", err)
	}

	return &redisService{client: client}
}

func (r *redisService) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if exp < 0 {
		return errors.New("expiration time cannot be negative")
	}

	serializedValue, err := serializeValue(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value for key %s: %w", key, err)
	}

	status := r.client.Set(ctx, key, serializedValue, exp)
	if err := status.Err(); err != nil {
		return fmt.Errorf("failed redis SET: %w", err)
	}

	return nil
}

// Target should be a reference to the type for data to be deserialized into
func (r *redisService) Get(ctx context.Context, key string, target any) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := deserializeValue(result, target); err != nil {
		return fmt.Errorf("failed to deserialize value for key %s: %w", key, err)
	}

	return nil
}

// Optionally, use the JSON.SET redis command
func (r *redisService) JSONSet(ctx context.Context, key, path string, value any, exp time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value to JSON: %w", err)
	}

	status := r.client.Do(ctx, "JSON.SET", key, path, jsonValue)
	if err := status.Err(); err != nil {
		return fmt.Errorf("failed redis JSON.SET for key %s: %w", key, err)
	}

	if exp > 0 {
		if err := r.client.Expire(ctx, key, exp).Err(); err != nil {
			return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
		}
	}

	return nil
}

// Optionally, use the JSON.GET redis command
func (r *redisService) JSONGet(ctx context.Context, key, path string, target any) error {
	result, err := r.client.Do(ctx, "JSON.GET", key, path).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s not found", key)
	} else if err != nil {
		return fmt.Errorf("failed redis JSON.GET for key %s: %w", key, err)
	}

	jsonData, ok := result.(string)
	if !ok {
		return fmt.Errorf("unexpected type for JSON.GET result: %T", result)
	}

	if err := json.Unmarshal([]byte(jsonData), target); err != nil {
		return fmt.Errorf("failed to deserialize JSON for key %s: %w", key, err)
	}

	return nil
}

func serializeValue(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}

// Users of this function must pass in a reference to the target value
// for data to be deserialized into
func deserializeValue(data string, target any) error {
	switch v := target.(type) {
	case *string:
		*v = data
		return nil
	case *[]byte:
		*v = []byte(data)
		return nil
	default:
		return json.Unmarshal([]byte(data), target)
	}
}
