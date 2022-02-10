package sa

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/ianjuma/change-monitor/sa/init"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

/*
Package listen is a self-contained Go program which uses the LISTEN / NOTIFY
mechanism to avoid polling the database while waiting for change events to arrive.
*/

const (
	activeEvent = "product_active_change_event"
	insertEvent = "product_insert_event"

	activeQueueName = "product:active-change"
	insertQueueName = "product:insert"
	timeout         = 11
	maxChannelSize  = 33 // big enough channel capacity to avoid blocking and hold events when redis is down
)

var (
	Rdb *redis.Client
)

// func init() {
// 	log.SetLevel(SetLogLevel(LogLevel))
// 	log.Info("initialising rdb client")
// 	Rdb = redis.NewClient(&redis.Options{
// 		Addr: RdbHostPort,
// 		DB:   0,
// 	})
// 	if Rdb == nil {
// 		panic("redis failed to init")
// 	}
// }

func Init() {
	//  wait for postgres and dependencies to be ready
	time.Sleep(10 * time.Second)
	const op = "main.entrypoint"
	log.WithFields(log.Fields{
		"op": op,
	}).Info("starting")

	var conninfo = fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		pdbUser,
		pdbPass,
		pdbName,
		pdbHost,
		"5432",
	)
	if doTrigger {
		createTriggers(conninfo)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		log.WithFields(log.Fields{
			"op":    op,
			"event": ev,
		}).Trace("listener connection change notification - 1/3 is bad")

		if err != nil {
			log.WithFields(log.Fields{
				"op":  op,
				"err": err,
			}).Error("reportProblem")
		}
	}
	minReconn := 10 * time.Millisecond
	maxReconn := timeout * time.Second
	log.WithFields(log.Fields{
		"minRecon ms": minReconn,
		"maxRecon  s": maxReconn,
	}).Info("conn values")

	activeListener := pq.NewListener(conninfo, minReconn, maxReconn, reportProblem)
	err := activeListener.Listen(activeEvent)
	if err != nil {
		log.WithFields(log.Fields{
			"op":  op,
			"err": err,
		}).Fatal("could not set active listener")
		return
	}

	insertListener := pq.NewListener(conninfo, minReconn, maxReconn, reportProblem) // list of listeners, from config file
	err = insertListener.Listen(insertEvent)
	if err != nil {
		log.WithFields(log.Fields{
			"op":  op,
			"err": err,
		}).Fatal("could not set insert listener")
		return
	}
	log.WithFields(log.Fields{
		"op": op,
	}).Info("entering main loop")

	quitChan := make(chan os.Signal, 1) // processing pending events
	signal.Notify(quitChan, syscall.SIGTERM, os.Interrupt)
	productEventChannel := make(chan *pq.Notification, maxChannelSize)

	ctx, cancel := context.WithCancel(context.Background())
	go listenForActiveEvent(ctx, activeListener, productEventChannel)
	go listenForInsertEvent(ctx, insertListener, productEventChannel)

	for {
		select {
		case product := <-productEventChannel:
			go processProductEvent(product)

		case <-time.After(time.Minute * 2):
			monitorBuffer(len(productEventChannel), cap(productEventChannel))

		// handle os exit events
		case <-quitChan:
			log.WithFields(log.Fields{
				"op": op,
			}).Info("received exit signal")

			// un-listen and close, un-listen alone might still send events
			_ = activeListener.Unlisten(activeEvent)
			_ = insertListener.Unlisten(insertEvent)
			_ = activeListener.Close()
			_ = insertListener.Close()

			// exit all listener goroutines
			cancel()

			// wait for pending for-select cancel events to be handled
			// before closing channel
			time.Sleep(time.Second * 3)

			log.WithFields(log.Fields{
				"op":   op,
				"size": len(productEventChannel),
			}).Info("processing remaining buffered events")

			// close channel after un-listen, so we can process any remaining events
			close(productEventChannel)

			// pull and wait for all work on the event channel
			// notice we have closed the channel so the range over channel loop terminates
			var wg sync.WaitGroup
			for notification := range productEventChannel { // todo: remove, not needed in an unbuffered channel
				wg.Add(1)
				go func(n *pq.Notification) {
					defer wg.Done()
					processProductEvent(n)
				}(notification)
			}
			wg.Wait()

			close(quitChan)
		}
	}
}
