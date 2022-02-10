package sa

import (
	log "github.com/sirupsen/logrus"
)

// monitorBuffer checks the channel size and warns if it exceeds the threshold
func monitorBuffer(length int, capacity int) {
	const op = "main.monitorBuffer"
	log.WithFields(log.Fields{
		"op":   op,
		"size": length,
		"cap":  capacity,
	}).Trace("ProductEventChannel buffer size")

	if (length/capacity)*100 > 30 {
		log.WithFields(log.Fields{
			"op":  op,
			"len": length,
			"cap": capacity,
		}).Warn("channel buffer is too large")
	}
}

func activeChanged(product ProductEvent) bool {
	return product.Old.Active != product.New.Active
}

func SetLogLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}
