package configCenter

import (
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Broker     string
	HttpBroker string
	clientId   string
	mqttClient MQTT.Client
	locker     *sync.Mutex
	Message    interface{}
	callback   func(topic string, response interface{})
}

var globThisMap sync.Map

func NewConfigCenter() *Config {
	cId := time.Now().Format("20060102150405")
	conn := &Config{clientId: cId}
	globThisMap.Store(cId, conn)
	return conn
}

func NewConfigWithBroker(httpbroker string, tcpbroker string) (*Config, error) {
	cId := time.Now().Format("20060102150405")
	conn := &Config{clientId: cId, Broker: tcpbroker, HttpBroker: httpbroker}
	globThisMap.Store(cId, conn)
	err := conn.connect()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func onSubscribeMessage(client MQTT.Client, message MQTT.Message) {
	optionReader := client.OptionsReader()
	clientId := optionReader.ClientID()
	this, ok := globThisMap.Load(clientId)
	if !ok {
		fmt.Println("not found clientId:", clientId)
		return
	}
	subClient := this.(*Config)
	subClient.callback(message.Topic(), message.Payload())
}

func (cc *Config) SubscribeAndQuery(topic string, callback func(topic string, response interface{})) error {
	cc.callback = callback
	err := cc.initTopic(topic)
	if err != nil {
		fmt.Println("init topic error:", err)
	}

	if token := cc.mqttClient.Subscribe(topic, 0, onSubscribeMessage); token.Wait() && token.Error() != nil {
		log.Println("subscribe topic:", topic, " error:", token.Error())
		return token.Error()
	}

	return nil
}

func (cc *Config) initTopic(topic string) error {
	if len(cc.HttpBroker) != 0 {
		if len(topic) == 0 {
			return errors.New("the topic is empty")
		}

		host := cc.HttpBroker + strings.Split(topic, ".")[1]

		resp, err := http.Get(host)
		if err != nil {
			fmt.Println("连接配置中心错误:", err)
			return err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取配置中心数据错误1:", err)
			return err
		}

		cc.Message = data

		return nil
	}
	return errors.New("no httpBroker")
}

func (cc *Config) connect() error {

	options := MQTT.NewClientOptions()
	options.SetAutoReconnect(true)
	options.SetKeepAlive(10 * time.Second)
	options.AddBroker(cc.Broker)

	options.SetCleanSession(true)
	options.SetClientID(cc.clientId)

	options.SetDefaultPublishHandler(onSubscribeMessage)
	cc.locker = &sync.Mutex{}

	cli := MQTT.NewClient(options)
	if !cli.IsConnected() {
		cc.locker.Lock()
		defer cc.locker.Unlock()
		if token := cli.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		cc.mqttClient = cli
	}
	return nil

}
