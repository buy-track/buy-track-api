package event

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"my-stocks/api-gateway/command"
	"my-stocks/common/aggregator"
	"my-stocks/common/broker"
	"my-stocks/common/queues"
	"my-stocks/domains"
)

const (
	AuthErrorGroup   queues.QueueNameGroup = "auth-error-group-aggregator"
	AuthSuccessGroup queues.QueueNameGroup = "auth-succeed-group-aggregator"
)

type AuthCommand struct {
	Aggregator
}

func NewAuthCommand(messageBroker broker.MessageBroker) command.AuthCommander {
	a := AuthCommand{
		Aggregator: Aggregator{
			manager: aggregator.NewCorrelationManager(), messageBroker: messageBroker,
		},
	}
	go a.consumeErrors()
	go a.consumerSuccesses()

	return a
}

func (a AuthCommand) Login(userId string) (*domains.Token, error) {
	data, receiver := a.manager.Add(userId)
	err := a.messageBroker.Push(queues.AuthQueuesGenerateToken, data)
	if err != nil {
		return nil, err
	}
	if result := <-receiver; result.Status {
		jsonStr, _ := json.Marshal(result.Data)
		var token domains.Token
		_ = json.Unmarshal(jsonStr, &token)
		return &token, nil
	} else {
		return nil, a.convertError(result)
	}

}

func (a AuthCommand) Logout(token string) error {
	data, receiver := a.manager.Add(token)
	err := a.messageBroker.Push(queues.AuthQueuesRemoveToken, data)
	if err != nil {
		return err
	}

	if result := <-receiver; result.Status {
		return nil
	} else {
		return a.convertError(result)
	}
}

func (a AuthCommand) RevokeAllTokens(userId string) error {
	data, receiver := a.manager.Add(userId)
	err := a.messageBroker.Push(queues.AuthQueuesRemoveAllToken, data)
	if err != nil {
		return err
	}

	if result := <-receiver; result.Status {
		return nil
	} else {
		return a.convertError(result)
	}
}

func (a AuthCommand) consumeErrors() {
	err := a.messageBroker.ConsumeWithHandler(queues.AuthQueuesGenerateTokenERROR, AuthErrorGroup, a.errorHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = a.messageBroker.ConsumeWithHandler(queues.AuthQueuesRemoveTokenERROR, AuthErrorGroup, a.errorHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = a.messageBroker.ConsumeWithHandler(queues.AuthQueuesRemoveAllTokenERROR, AuthErrorGroup, a.errorHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
}

func (a AuthCommand) consumerSuccesses() {
	err := a.messageBroker.ConsumeWithHandler(queues.AuthQueuesGenerateTokenSUCCESS, AuthSuccessGroup, a.succeedHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = a.messageBroker.ConsumeWithHandler(queues.AuthQueuesRemoveTokenSUCCESS, AuthSuccessGroup, a.succeedHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = a.messageBroker.ConsumeWithHandler(queues.AuthQueuesRemoveAllTokenSUCCESS, AuthSuccessGroup, a.succeedHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
}
