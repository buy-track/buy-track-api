package event

import (
	log "github.com/sirupsen/logrus"
	"my-stocks/common/aggregator"
	"my-stocks/common/broker"
)

type Aggregator struct {
	manager       *aggregator.CorrelationManager
	messageBroker broker.MessageBroker
}

func (a Aggregator) convertError(correlation *aggregator.RespondCorrelation) *broker.Error {
	return correlation.Data.(*broker.Error)
}

func (a Aggregator) errorHandler(message *broker.Message) {
	a.manager.Respond(message.Data.CorrelationId, &aggregator.RespondCorrelation{
		CorrelationId: message.Data.CorrelationId,
		Status:        false,
		Data:          message.Data.Data,
	})
	err := a.messageBroker.Ack(message.Queue, message.Group, message.ID)
	if err != nil {
		log.Printf("error during acknowledge")
	}
}

func (a Aggregator) succeedHandler(message *broker.Message) {
	a.manager.Respond(message.Data.CorrelationId, &aggregator.RespondCorrelation{
		CorrelationId: message.ID,
		Status:        true,
		Data:          message.Data.Data,
	})
	err := a.messageBroker.Ack(message.Queue, message.Group, message.ID)
	if err != nil {
		log.Printf("error during acknowledge")
	}
}
