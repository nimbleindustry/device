package common

import (
	"fmt"
	"testing"
	"time"
)

func sleepMillis(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func trace(args ...interface{}) {
	fmt.Println(args)
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	channel := "TestSendMessage"
	var message string
	go func() {
		trace("waiting for message")
		message = WaitForBusMessage(channel).(string)
		return
	}()
	sleepMillis(10)
	SendBusMessage(channel, "bar")
	sleepMillis(10)
	if message != "bar" {
		fmt.Println("expected message to be bar")
		t.Fail()
	}
}

func TestSendMessageMultipleChannels(t *testing.T) {
	t.Parallel()
	channel1 := "channel1"
	channel2 := "channel2"
	var message string
	go func() {
		trace("waiting for message on channel1")
		message = WaitForBusMessage(channel1).(string)
		return
	}()
	go func() {
		trace("waiting for message on channel2")
		<-BusChannel(channel2)
		t.Fail()
		return
	}()
	sleepMillis(10)
	SendBusMessage(channel1, "bar")
	sleepMillis(10)
	if message != "bar" {
		fmt.Println("expected message to be bar")
		t.Fail()
	}

}

func TestSendMessageChannel(t *testing.T) {
	t.Parallel()
	channel := "TestSendMessageChannelfoo"
	var message string
	go func() {
		trace("waiting for message")
		for {
			//message = WaitOnBusChannel(channel).(string)
			message = WaitForBusMessage(channel).(string)
			break
		}
		return
	}()
	sleepMillis(10)
	SendBusMessage(channel, "barbar")
	sleepMillis(10)
	if message != "barbar" {
		fmt.Println("expected message to be barbar")
		t.Fail()
	}
}

func TestSendMultipleMessages(t *testing.T) {
	t.Parallel()
	channel := "TestSendMultipleMessages"
	count := 0
	go func() {
		trace("waiting for message")
		for {
			val := WaitForBusMessage(channel).(string)
			if val != "one" && val != "two" && val != "three" {
				fmt.Println("expected string one|two|three, got", val)
				t.Fail()
			}
			count++
			if count == 3 {
				break
			}
		}
		return
	}()
	sleepMillis(20)
	SendBusMessage(channel, "one")
	sleepMillis(20)
	SendBusMessage(channel, "two")
	sleepMillis(20)
	SendBusMessage(channel, "three")
	sleepMillis(20)
	if count != 3 {
		fmt.Println("expected to receive 3 messages, got", count)
		t.Fail()
	}
}

func TestSendToMultipleReceivers(t *testing.T) {
	t.Parallel()
	channel := "TestSendToMultipleReceivers"
	count := 0

	for i := 0; i < 3; i++ {
		go func() {
			fmt.Println("waiting for message")
			for {
				val := WaitForBusMessage(channel).(string)
				if val != "foo" {
					fmt.Println("expected string 'foo', got", val)
				}
				count++
				if count == 3 {
					break
				}
			}
			return
		}()
	}
	sleepMillis(20)
	SendBusMessage(channel, "foo")
	sleepMillis(20)
	if count != 3 {
		fmt.Println("expected to receive 3 messages, got", count)
		t.Fail()
	}
}

func TestWaitForBusMessage(t *testing.T) {
	t.Parallel()
	channel := "TestWaitForBusMessage"
	satisfied := false

	go func() {
		for {
			WaitForBusMessage(channel)
			satisfied = true
		}
	}()

	sleepMillis(1000)
	if satisfied {
		fmt.Println("expected false value for variable 'satisfied'")
		t.Fail()
	}
}

// this doesn't work, bcast doesn't store messages if no listeners are there
func NOTestSendBeforeWait(t *testing.T) {
	channel := "sendBeforeWait"

	SendBusMessage(channel, true)
	val := false
	go func() {
		for {
			select {
			case received := <-BusChannel(channel):
				trace("got value", received)
				val = received.(bool)
				break
			}
		}
	}()
	sleepMillis(20)
	if !val {
		fmt.Println("Expected to receive true value")
		t.Fail()
	}
	SendBusMessage(channel, true)
	sleepMillis(20)
}
