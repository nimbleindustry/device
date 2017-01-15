package common

import (
	"sync"

	"github.com/nimbleindustry/bcast"
)

type bus struct {
	groupMap      map[string]*bcast.Group
	creationMutex *sync.Mutex
	sendMutex     *sync.Mutex
}

var instance *bus
var once sync.Once

func getInstance() *bus {
	once.Do(func() {
		instance = &bus{}
		instance.creationMutex = &sync.Mutex{}
		instance.sendMutex = &sync.Mutex{}
		instance.groupMap = make(map[string]*bcast.Group, 0)
	})
	return instance
}

func (bus *bus) getBusGroup(topic string) *bcast.Group {
	bus.creationMutex.Lock()
	defer bus.creationMutex.Unlock()
	if group, found := bus.groupMap[topic]; found {
		//fmt.Println("returning found group", topic)
		return group
	}
	group := bcast.NewGroup()
	bus.groupMap[topic] = group
	go group.Broadcast()
	return group
}

// SendBusMessage sends the supplied message into the bus channel topic.
// If no listeners are waiting on the topic, this call is a no-op (as there is no store
// and forward)
func SendBusMessage(topic string, message interface{}) {
	getInstance().sendMutex.Lock()
	defer getInstance().sendMutex.Unlock()
	group := getInstance().getBusGroup(topic)
	//fmt.Println("sending message to group", topic)
	group.Send(message)
}

// WaitForBusMessage blocks until receiving a message on bus channel topic.
// It returns what ever primitive or object was sent.
func WaitForBusMessage(topic string) interface{} {
	return getInstance().getBusGroup(topic).Join().Recv()
}

// BusChannel returns the channel associated with the topic. This is useful
// for 'select' operations in which a loop is blocked on the channel.
func BusChannel(topic string) chan interface{} {
	return getInstance().getBusGroup(topic).Join().In
}
