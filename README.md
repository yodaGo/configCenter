# golang connect mqtt as configCenter

httpBroker:="http://127.0.0.1:5000/xxx/topics/"
mqtt_broker:= "tcp://127.0.0.1:1822"
cfg:=configCenter.NewConfigWithBroker(httpBroker,mqtt_broker)
cfg.SubscribeAndQuery("testTopicName"",callback)
msg:=cfg.Message.([]byte)
fmt.Println("===========>",string(msg))
select {

}


func callback(topic string,msg  interface{}) {
	response := msg.([]byte)
 	fmt.Println(topic,"=================>",string(response))
}