package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Cache struct {
	client *redis.Client
}

func NewCache() (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Cache{client: client}, nil
}

func (c *Cache) SetUserState(userID int64, userState UserState) error {
	ctx := context.Background()

	err := c.client.Set(ctx, strconv.FormatInt(userID, 10), userState.String(), 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *Cache) GetUserState(userID int64) *UserState {
	ctx := context.Background()

	strCmd := c.client.Get(ctx, strconv.FormatInt(userID, 10))
	if strCmd.Err() == redis.Nil {
		return nil
	}
	return NewUserStateFromString(strCmd.Val())
}
