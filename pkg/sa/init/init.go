package init

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/ianjuma/change-monitor/pkg/sa"
)

func init() {

	log.Info("initialising rdb client")
	sa.Rdb = redis.NewClient(&redis.Options{
		Addr: sa.RdbHostPort,
		DB:   0,
	})
	log.SetLevel(sa.SetLogLevel(sa.LogLevel))
}
