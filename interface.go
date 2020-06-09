package golang_config_center_connector

type ConfigCenter interface {
	SubscribeAndQuery(topic string, callback func(topic string, response interface{})) error
}
