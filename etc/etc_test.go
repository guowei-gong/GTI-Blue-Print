package etc_test

import (
	"github.com/gti-blue-print/config/etc"
	"testing"
)

func Test_Get(t *testing.T) {
	// v := etc.Get("c.redis.addrs.1A", "192.168.0.1:3308").String()
	v := etc.Get("dev.locate.redis.addrs", "read-only").Strings()
	t.Log(v)
}
