package event

import (
	"fmt"
	"my-stocks/auth/app"
	"my-stocks/common/aggregator"
	"my-stocks/common/broker"
	"my-stocks/common/queues"
)

type AuthHandler struct {
	authService app.AuthService
}

func NewAuthHandler(authService app.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a AuthHandler) CreateToken(message *broker.Message, p broker.Publisher) {
	input := message.Data.Data.(string)
	token, err := a.authService.GenerateAccessToken(input)
	if err != nil {
		_ = p.Push(queues.AuthQueuesGenerateTokenERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to generate token : %v", err)),
		})
		return
	}
	_ = p.Push(queues.AuthQueuesGenerateTokenSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          token,
	})
}

func (a AuthHandler) RemoveToken(message *broker.Message, p broker.Publisher) {
	input := message.Data.Data.(string)
	err := a.authService.DeleteAccessToken(input)
	if err != nil {
		_ = p.Push(queues.AuthQueuesRemoveTokenERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to remove token : %v", err)),
		})
		return
	}
	_ = p.Push(queues.AuthQueuesRemoveTokenSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          true,
	})
}

func (a AuthHandler) RemoveAllToken(message *broker.Message, p broker.Publisher) {
	input := message.Data.Data.(string)
	err := a.authService.DeleteAllAccessToken(input)
	if err != nil {
		_ = p.Push(queues.AuthQueuesRemoveAllTokenERROR, &aggregator.Correlation{
			CorrelationId: message.Data.CorrelationId,
			Data:          broker.NewError(message.Queue, fmt.Sprintf("failed to remove all tokens : %v", err)),
		})
		return
	}
	_ = p.Push(queues.AuthQueuesRemoveAllTokenSUCCESS, &aggregator.Correlation{
		CorrelationId: message.Data.CorrelationId,
		Data:          true,
	})
}
