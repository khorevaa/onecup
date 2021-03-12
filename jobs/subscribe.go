package jobs

import (
	"fmt"
)

type Subscribe struct {
	events       chan Event
	eventsFilter []EventType
}

func (s Subscribe) emmitEvent(event Event) {

	fmt.Printf("emit event <%v>\n", event)

}

func NewSubscriber(eventsFilter ...EventType) *Subscribe {
	events := make(chan Event)

	return &Subscribe{
		events:       events,
		eventsFilter: eventsFilter,
	}

}
