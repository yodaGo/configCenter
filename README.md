# golang connect mqtt as configCenter

帮助:
httpBroker:="http://127.0.0.1:5000/xxx/topics/"
mqtt_broker:= "tcp://127.0.0.1:1822"
cfg:=configCenter.NewConfigWithBroker(httpBroker,mqtt_broker)
cfg.SubscribeAndQuery("cfg.UserRouterAssign",callback)
msg:=cfg.Message.([]byte)
fmt.Println("===========>",string(msg))
select {

}