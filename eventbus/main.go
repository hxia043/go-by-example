package main

import (
	"eventbus/eventbus"
	"fmt"
	"sync"
	"time"

	"math/rand"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent
type DataChannelSlice []DataChannel

type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

var eb = &EventBus{
	subscribers: make(map[string]DataChannelSlice),
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.subscribers[topic]; found {
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlice DataChannelSlice) {
			for _, ch := range dataChannelSlice {
				ch <- data
			}
		}(DataEvent{Topic: topic, Data: data}, channels)
	}
	eb.rm.RUnlock()
}

func publisTo(topic string, data string) {
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func printDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func calculator(a int, b int) {
	fmt.Printf("%d\n", a+b)
}

func main() {
	bus := eventbus.New()
	bus.Subscribe("main:calculator", calculator)
	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator)

	/*
		ch1 := make(chan DataEvent)
		ch2 := make(chan DataEvent)
		ch3 := make(chan DataEvent)

		eb.Subscribe("topic1", ch1)
		eb.Subscribe("topic2", ch2)
		eb.Subscribe("topic3", ch3)

		go publisTo("topic1", "Hi topic 1")
		go publisTo("topic2", "Welcome to topic 2")

		for {
			select {
			case d := <-ch1:
				go printDataEvent("ch1", d)
			case d := <-ch2:
				go printDataEvent("ch2", d)
			case d := <-ch3:
				go printDataEvent("ch3", d)
			}
		}
	*/
}
