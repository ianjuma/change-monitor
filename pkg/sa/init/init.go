package init

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/ianjuma/change-monitor/pkg/sa"
)

func init() {
	log.SetLevel(sa.SetLogLevel(sa.LogLevel))
	log.Info("initialising rdb client")
	sa.Rdb = redis.NewClient(&redis.Options{
		Addr: sa.RdbHostPort,
		DB:   0,
	})
	if sa.Rdb == nil {
		panic("redis failed to init")
	}
}
