package client

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/toolkit"
	"time"
)

type DatahubClient struct {
}

type Method string

const (
	ERR Method = "ERR"
	REQ Method = "REQ"
	RSP Method = "RSP"
	PUB Method = "PUB"
	SUB Method = "SUB"
)

type TMsg struct {
	ID      string
	Sender  string
	Method  Method
	Key     string
	Payload interface{}
	QoS     byte
}

func NewTMsg() *TMsg {
	return &TMsg{ID: toolkit.NewUUID(), QoS: 1}
}

func (t *TMsg) WithMethod(method Method) {
	t.Method = method
}

func (t *TMsg) WithKey(key string) {
	t.Key = key
}

func (t *TMsg) WithPayload(payload interface{}) {
	t.Payload = payload
}

func (t *TMsg) WithQoS(qos byte) {
	t.QoS = qos
}

type TMsgHandler func(*TMsg)

type TClient struct {
	clientID  string
	client    mqtt.Client
	qos       byte
	timeout   time.Duration
	callQueue map[int64]chan *TMsg
}

func parseOpts(config conf.MqttConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.SetConnectTimeout(time.Duration(config.ConnectTimeoutMS) * time.Millisecond)
	opts.AddBroker("tcp://" + config.BrokerAddr)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(60 * time.Second)
	//opts.SetCleanSession(false)
	opts.SetOnConnectHandler(onConnect)
	opts.SetConnectionLostHandler(onConnectionLost)
	return opts
}

func onConnect(c mqtt.Client) {
	opts := c.OptionsReader()
	log.Infof("MQTTClient %s connected to %s", opts.ClientID(), opts.Servers()[0].String())
	return
}

func onConnectionLost(c mqtt.Client, err error) {
	opts := c.OptionsReader()
	log.WithError(err).Errorf("MQTTClient %s connection to %s has lost, will try to reconnect", opts.ClientID(), opts.Servers()[0].String())
	return
}

func NewTClient(c conf.MqttConfig, clientId string) *TClient {
	mqttOpts := parseOpts(c)
	mqttOpts.ClientID = clientId
	return &TClient{
		clientID:  clientId,
		client:    mqtt.NewClient(mqttOpts),
		qos:       c.QoS,
		timeout:   time.Duration(c.RequestTimeoutMS) * time.Millisecond,
		callQueue: make(map[int64]chan *TMsg),
	}
}

const MqttChannelPrefix = "/channel/"

func (t *TClient) Connect() error {
	tk := t.client.Connect()
	tk.Wait()
	if err := tk.Error(); err != nil {
		return err
	}

	channelTopic := MqttChannelPrefix + t.clientID
	tk = t.client.Subscribe(channelTopic, t.qos, t.callHandler)
	tk.Wait()
	if err := tk.Error(); err != nil {
		return err
	}

	return nil
}

func (t *TClient) callHandler(cli mqtt.Client, msg mqtt.Message) {
	tmsg := &TMsg{}
	err := json.Unmarshal(msg.Payload(), tmsg)
	if err != nil {
		log.Errorf("fail to unmarshal data %+v", msg)
	}

	if tmsg.Method == REQ {

	}
}

func (t *TClient) RegHandler(topic string, handler TMsgHandler) chan *TMsg {

}

func (t *TClient) Send(msg *TMsg) chan *TMsg {

}

func (t *TClient) SendAndWait(ctx context.Context, msg *TMsg) (*TMsg, error) {
}

func (t *TClient) Watch(key string) chan *TMsg {
}
