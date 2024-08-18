package broker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"my-stocks/common/aggregator"
	"my-stocks/common/queues"
)

type Error struct {
	Queue   queues.QueueName
	Message string
}

func NewError(queue queues.QueueName, message string) *Error {
	return &Error{Queue: queue, Message: message}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) String() string {
	return e.Message
}

type Message struct {
	ID    string
	Group queues.QueueNameGroup
	Queue queues.QueueName
	Data  *aggregator.Correlation
}

type Consumer interface {
	ConsumeWithHandler(queue queues.QueueName, group queues.QueueNameGroup, handler ConsumerOnly) error
	Consume(queue queues.QueueName, group queues.QueueNameGroup, sender chan<- *Message) error
	Ack(queue queues.QueueName, group queues.QueueNameGroup, messageId string) error
}

type Publisher interface {
	Push(queue queues.QueueName, data *aggregator.Correlation) error
}

type MessageBroker interface {
	Consumer
	Publisher
}

type ConsumerOnly func(message *Message)

type Handler func(message *Message, pusher Publisher)

type EventDriven struct {
	queues       map[queues.QueueNameGroup]map[queues.QueueName]Handler
	dataReceiver chan *Message
	publisher    Publisher
	consumer     Consumer
	ctx          context.Context
	graceFull    <-chan bool
}

func New(messageBroker MessageBroker, ctx context.Context, graceFull <-chan bool) *EventDriven {
	return NewEventDriven(messageBroker, messageBroker, ctx, graceFull)
}

func NewEventDriven(publisher Publisher, consumer Consumer, ctx context.Context, graceFull <-chan bool) *EventDriven {
	dataReceiver := make(chan *Message)
	q := make(map[queues.QueueNameGroup]map[queues.QueueName]Handler)
	return &EventDriven{dataReceiver: dataReceiver, publisher: publisher, consumer: consumer, ctx: ctx, graceFull: graceFull, queues: q}
}

func (e *EventDriven) AddRoute(queue queues.QueueName, group queues.QueueNameGroup, handler Handler) {
	if _, ok := e.queues[group]; !ok {
		e.queues[group] = make(map[queues.QueueName]Handler)
	}
	e.queues[group][queue] = handler
}

func (e *EventDriven) Start() {
	e.initRoutes()
}

func (e *EventDriven) initRoutes() {
	for group, handlers := range e.queues {
		log.Printf("Start group %v", group)
		for queue, _ := range handlers {
			log.Printf("Start queue consumer %v", queue)
			go e.startConsume(queue, group)
		}
	}
	go e.startHandlers()
}

func (e *EventDriven) startConsume(queue queues.QueueName, group queues.QueueNameGroup) {
	err := e.consumer.Consume(queue, group, e.dataReceiver)
	if err != nil {
		log.Fatalf("error during consume %v", err)
	}
}

func (e *EventDriven) startHandlers() {
	for {
		select {
		case message := <-e.dataReceiver:
			if found, ok := e.queues[message.Group][message.Queue]; ok {
				found(message, e.publisher)
				err := e.consumer.Ack(message.Queue, message.Group, message.ID)
				if err != nil {
					log.Printf("error during acknowledge! %v", err)
				}
			} else {
				log.Printf("There is no handler for group: %v queue: %v", message.Group, message.Queue)
			}
		case closed := <-e.graceFull:
			if closed {
				close(e.dataReceiver)
				return
			}
		}
	}
}
