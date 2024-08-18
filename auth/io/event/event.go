package event

import (
	"github.com/golobby/container/v3"
	"my-stocks/auth/app"
	"my-stocks/common/broker"
	"my-stocks/common/queues"
)

var driven *broker.EventDriven

func AddRoutes(eventDriven *broker.EventDriven, ctr container.Container) {
	driven = eventDriven
	var authService app.AuthService
	_ = ctr.Resolve(&authService)

	authHandler := NewAuthHandler(authService)

	driven.AddRoute(queues.AuthQueuesGenerateToken, queues.AuthQueueGroup, authHandler.CreateToken)
	driven.AddRoute(queues.AuthQueuesRemoveToken, queues.AuthQueueGroup, authHandler.RemoveToken)
	driven.AddRoute(queues.AuthQueuesRemoveAllToken, queues.AuthQueueGroup, authHandler.RemoveAllToken)
}

func Start() {
	driven.Start()
}
