package main

import "fmt"

// Event Bus
//
// An event bus delivers messages to all registered subscribers.
// Each subscriber accumulates the events it receives.
//
// Expected output:
//   received: [order.created order.shipped]

type Subscriber interface {
    OnEvent(event string)
}

type Collector struct {
    received []string
}

func (c *Collector) OnEvent(event string) {
    c.received = append(c.received, event)
}

type Bus struct {
    subs []Subscriber
}

func (b *Bus) Subscribe(s Subscriber) {
    b.subs = append(b.subs, s)
}

func (b *Bus) Publish(event string) {
    for _, s := range b.subs {
        s.OnEvent(event)
    }
}

func main() {
    bus := &Bus{}
    c := Collector{}
    bus.Subscribe(&c)
    bus.Publish("order.created")
    bus.Publish("order.shipped")
    fmt.Println("received:", c.received)
}
