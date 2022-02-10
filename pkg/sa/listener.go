package sa

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	RedisTimeout    = 5 * time.Second
	ListenerTimeout = 11 * time.Second // used to avoid timeout on listener - for-select timeout
)

func listenForInsertEvent(ctx context.Context, l *pq.Listener, eventChan chan<- *pq.Notification) {
	const op = "product.listenForInsertEvent"

	for {
		select {
		case event := <-l.Notify:
			log.WithFields(log.Fields{
				"op":    op,
				"event": event,
			}).Debug("received notification")
			eventChan <- event

		case <-time.After(ListenerTimeout): // avoid a timeout on listen, yield after timeout
			go func() {
				err := l.Ping()
				if err != nil {
					log.WithFields(log.Fields{
						"op":  op,
						"err": err,
					}).Error("could not ping the database")
				}
			}()

		case <-ctx.Done():
			log.WithFields(log.Fields{
				"op": op,
			}).Trace("exiting")
			return
		}
	}
}

func listenForActiveEvent(ctx context.Context,
	l *pq.Listener,
	eventChan chan<- *pq.Notification) {
	const op = "product.listenForActiveEvent"

	for {
		select {
		case event := <-l.Notify:
			log.WithFields(log.Fields{
				"op": op,
			}).Info("received event")
			eventChan <- event

		case <-time.After(ListenerTimeout):
			go func() {
				err := l.Ping()
				if err != nil {
					log.WithFields(log.Fields{
						"op":  op,
						"err": err,
					}).Error("could not ping the database")
				}
			}()

		case <-ctx.Done():
			log.WithFields(log.Fields{
				"op": op,
			}).Trace("exiting")
			return
		}
	}
}

func processProductEvent(p *pq.Notification) {
	const op = "processProductEvent"
	if p == nil {
		return
	}

	var product ProductEvent
	if err := json.Unmarshal([]byte(p.Extra), &product); err != nil {
		log.WithFields(log.Fields{
			"op":  op,
			"err": err,
		}).Error("failed to Unmarshall event")
		return
	}

	switch p.Channel {
	case activeEvent:
		if !activeChanged(product) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), RedisTimeout)
		defer cancel()
		log.Info(p.Extra, activeQueueName)
		Rdb.LPush(ctx, activeQueueName, p.Extra) //panic

	case insertEvent:
		ctx, cancel := context.WithTimeout(context.Background(), RedisTimeout)
		defer cancel()
		log.Info(p.Extra, insertQueueName)
		Rdb.LPush(ctx, insertQueueName, p.Extra) //panic
	}
}
