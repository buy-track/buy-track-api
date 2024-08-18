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
	UserErrorGroup   queues.QueueNameGroup = "user-error-group-aggregator"
	UserSuccessGroup queues.QueueNameGroup = "user-succeed-group-aggregator"
)

type UserCommand struct {
	Aggregator
}

func NewUserCommand(messageBroker broker.MessageBroker) command.UserCommander {
	a := UserCommand{
		Aggregator: Aggregator{
			manager: aggregator.NewCorrelationManager(), messageBroker: messageBroker,
		},
	}
	go a.consumeErrors()
	go a.consumerSuccesses()

	return a
}

func (u UserCommand) Create(data *domains.User) (*domains.User, error) {
	userData, receiver := u.manager.Add(data)
	err := u.messageBroker.Push(queues.UserQueuesCreate, userData)
	if err != nil {
		return nil, err
	}

	if result := <-receiver; result.Status {
		jsonString, _ := json.Marshal(result.Data)
		input := new(domains.User)
		json.Unmarshal(jsonString, input)
		return input, nil
	} else {
		return nil, u.convertError(result)
	}
}

func (u UserCommand) AddProviderToken(userId string, token string, provider domains.Provider) error {
	userData, receiver := u.manager.Add(&domains.ProviderToken{
		UserId:     userId,
		ProviderId: token,
		Provider:   provider,
	})

	err := u.messageBroker.Push(queues.UserQueuesStoreProviderToken, userData)
	if err != nil {
		return err
	}

	if result := <-receiver; result.Status {
		return nil
	} else {
		return u.convertError(result)
	}
}

func (u UserCommand) consumeErrors() {
	err := u.messageBroker.ConsumeWithHandler(queues.UserQueuesCreateERROR, UserErrorGroup, u.errorHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = u.messageBroker.ConsumeWithHandler(queues.UserQueuesStoreProviderTokenERROR, UserErrorGroup, u.errorHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
}

func (u UserCommand) consumerSuccesses() {
	err := u.messageBroker.ConsumeWithHandler(queues.UserQueuesCreateSUCCESS, UserSuccessGroup, u.succeedHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
	err = u.messageBroker.ConsumeWithHandler(queues.UserQueuesStoreProviderTokenSUCCESS, UserSuccessGroup, u.succeedHandler)
	if err != nil {
		log.Printf("error during consume : %v", err)
	}
}
