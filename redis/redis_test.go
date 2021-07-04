package redis

import (
	"testing"

	"github.com/changyenh/sessions"
	"github.com/changyenh/sessions/tester"
)

const redisTestServer = "localhost:6379"

var newRedisStore = func(_ *testing.T) sessions.Store {
	store, err := NewStore(10, redisTestServer, "")
	if err != nil {
		panic(err)
	}
	return store
}

func TestRedis_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newRedisStore)
}

func TestRedis_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newRedisStore)
}

func TestRedis_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newRedisStore)
}

func TestRedis_SessionClear(t *testing.T) {
	tester.Clear(t, newRedisStore)
}

func TestRedis_SessionOptions(t *testing.T) {
	tester.Options(t, newRedisStore)
}

func TestRedis_SessionMany(t *testing.T) {
	tester.Many(t, newRedisStore)
}

func TestGetRedisStore(t *testing.T) {
	t.Run("unmatched type", func(t *testing.T) {
		type store struct{ Store }
		err, redisStore := GetRedisStore(store{})
		if err == nil || redisStore != nil {
			t.Fail()
		}
	})
}
