package queues

type QueueNameGroup string

const (
	defaultGroup   QueueNameGroup = "default"
	AuthQueueGroup QueueNameGroup = "auth-consumer-group"
	UserQueueGroup QueueNameGroup = "user-consumer-group"
)

type QueueName string

// AuthQueue names
const (
	defaultQueue             QueueName = "default"
	AuthQueuesGenerateToken  QueueName = "auth_generate_token"
	AuthQueuesRemoveToken    QueueName = "auth_remove_token"
	AuthQueuesRemoveAllToken QueueName = "auth_remove_all_token"
)

// AuthQueue Success const
const (
	AuthQueuesGenerateTokenSUCCESS  QueueName = "auth_generate_token_success"
	AuthQueuesRemoveTokenSUCCESS    QueueName = "auth_remove_token_success"
	AuthQueuesRemoveAllTokenSUCCESS QueueName = "auth_remove_all_token_success"
)

// AuthQueue Error const
const (
	AuthQueuesGenerateTokenERROR  QueueName = "auth_generate_token_error"
	AuthQueuesRemoveTokenERROR    QueueName = "auth_remove_token_error"
	AuthQueuesRemoveAllTokenERROR QueueName = "auth_remove_all_token_error"
)

// UserQueue const
const (
	UserQueuesCreate                 QueueName = "user_create"
	UserQueuesStoreProviderToken     QueueName = "user_store_provider_token"
	UserQueuesRemoveAllProviderToken QueueName = "user_remove_all_provider_token"
)

// UserQueue Success const
const (
	UserQueuesCreateSUCCESS                 QueueName = "user_create_success"
	UserQueuesStoreProviderTokenSUCCESS     QueueName = "user_store_provider_token_success"
	UserQueuesRemoveAllProviderTokenSUCCESS QueueName = "user_remove_all_provider_token_success"
)

// UserQueue Error const
const (
	UserQueuesCreateERROR                 QueueName = "user_create_error"
	UserQueuesStoreProviderTokenERROR     QueueName = "user_store_provider_token_error"
	UserQueuesRemoveAllProviderTokenERROR QueueName = "user_remove_all_provider_token_error"
)
