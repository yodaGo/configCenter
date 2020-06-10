package configCenter

type ConfigCenter interface {
	SubscribeAndQuery(topic string, callback func(topic string, response interface{})) error
}
