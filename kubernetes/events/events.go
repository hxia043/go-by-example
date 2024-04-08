package main

import (
	"fmt"
	"sync"
	"time"

	"events/types"
	"events/watch"
)

const queueLength = int64(1)

type EventBroadcaster interface {
	Eventf(etype, reason, message string)
	StartLogging() watch.Interface
	Stop()
}

type eventBroadcasterImpl struct {
	*watch.Broadcaster
}

func (eventBroadcaster *eventBroadcasterImpl) Stop() {
	eventBroadcaster.Shutdown()
}

func NewEventBroadcaster() EventBroadcaster {
	return &eventBroadcasterImpl{watch.NewBroadcaster(queueLength)}
}

func (eventBroadcaster *eventBroadcasterImpl) Eventf(etype, reason, message string) {
	events := &types.Events{Type: etype, Reason: reason, Message: message}
	eventBroadcaster.Action(events)
}

func (eventBroadcaster *eventBroadcasterImpl) StartLogging() watch.Interface {
	watcher := eventBroadcaster.Watch()
	go func() {
		for watchEvent := range watcher.ResultChan() {
			fmt.Printf("%v\n", watchEvent)
		}
	}()

	// test watcher client stop
	go func() {
		time.Sleep(time.Second * 4)
		watcher.Stop()
	}()

	return watcher
}

func main() {
	eventBroadcast := NewEventBroadcaster()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		eventBroadcast.Eventf("add", "test", "1")
		time.Sleep(time.Second * 2)
		eventBroadcast.Eventf("add", "test", "2")
		time.Sleep(time.Second * 3)
		eventBroadcast.Eventf("add", "test", "3")
	}()

	eventBroadcast.StartLogging()
	wg.Wait()
}
