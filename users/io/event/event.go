package event

import (
	"github.com/golobby/container/v3"
	"my-stocks/common/broker"
	"my-stocks/common/queues"
	"my-stocks/users/app"
)

var driven *broker.EventDriven

func AddRoutes(eventDriven *broker.EventDriven, ctr container.Container) {
	driven = eventDriven

	var userService app.UserService
	_ = ctr.Resolve(&userService)
	userHandler := NewUserHandler(userService)
	driven.AddRoute(queues.UserQueuesCreate, queues.UserQueueGroup, userHandler.Create)

	var providerService app.ProviderTokenService
	_ = ctr.Resolve(&providerService)
	providerHandler := NewProviderHandler(providerService)
	driven.AddRoute(queues.UserQueuesStoreProviderToken, queues.UserQueueGroup, providerHandler.CreateProviderToken)
	driven.AddRoute(queues.UserQueuesRemoveAllProviderToken, queues.UserQueueGroup, providerHandler.RemoveAllProviderToken)

}

func Start() {
	driven.Start()
}
