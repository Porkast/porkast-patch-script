package cache

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_initRedisClient(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = gctx.New()
			err error
		)
		redisClient := initRedisClient(ctx)
		if redisClient == nil {
			t.Fatal("init redis client failed")
		}
		_, err = redisClient.Do(ctx, "SET", "Test", "test_value")
		if err != nil {
			t.Fatal("init redis client failed : \n", err)
		}
	})
}

func Test_initRedisConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ctx = gctx.New()
		redisConfig := initRedisConfig(ctx)
		if redisConfig.Address == "" {
			t.Fatal("redis config address is empty")
		}

		if redisConfig.Pass == "" {
			t.Fatal("redis config password is empty")
		}
	})
}
