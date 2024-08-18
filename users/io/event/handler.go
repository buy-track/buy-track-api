package event

import (
	"encoding/json"
	"fmt"
	"my-stocks/common/aggregator"
	"my-stocks/common/broker"
	"my-stocks/common/queues"
	"my-stocks/domains"
	"my-stocks/users/app"
)

type UserHandler struct {
	userService app.UserService
}

func (u UserHandler) Create(message *broker.Message, publisher broker.Publisher) {
	ser := message.Data.Data.(map[string]interface{})
	jsonString, _ := json.Marshal(ser)
	input := new(domains.User)
	json.Unmarshal(jsonString, input)

	register, err := u.userService.Register(app.CreateUserDto{
		Password: input.Password,
		Email:    input.Email,
		Name:     input.Name,
	})
	if err != nil {
		_ = publisher.Push(queues.UserQueuesCreateERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to create user : %v", err)),
		})
		return
	}
	_ = publisher.Push(queues.UserQueuesCreateSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          register,
	})
}

func NewUserHandler(userService app.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type ProviderHandler struct {
	providerService app.ProviderTokenService
}

func NewProviderHandler(providerService app.ProviderTokenService) *ProviderHandler {
	return &ProviderHandler{providerService: providerService}
}

func (p ProviderHandler) CreateProviderToken(message *broker.Message, publisher broker.Publisher) {
	input := message.Data.Data.(*domains.ProviderToken)
	err := p.providerService.AddProviderToken(input.UserId, input.ProviderId, input.Provider)
	if err != nil {
		_ = publisher.Push(queues.UserQueuesStoreProviderTokenERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to store provider token : %v", err)),
		})
		return
	}
	_ = publisher.Push(queues.UserQueuesStoreProviderTokenSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          true,
	})
}

func (p ProviderHandler) RemoveAllProviderToken(message *broker.Message, publisher broker.Publisher) {
	input := message.Data.Data.(string)
	success := p.providerService.DeleteAllProviderToken(input)
	if !success {
		_ = publisher.Push(queues.UserQueuesRemoveAllProviderTokenERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to remove all provider token")),
		})
		return
	}
	_ = publisher.Push(queues.UserQueuesRemoveAllProviderTokenSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          true,
	})
}
