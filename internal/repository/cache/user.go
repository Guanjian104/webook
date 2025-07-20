package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Guanjian104/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

// type UserCache interface {
// 	Get(ctx context.Context, uid int64) (domain.User, error)
// 	Set(ctx context.Context, du domain.User) error
// }

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *UserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
