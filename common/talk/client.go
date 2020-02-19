package talk

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
	"time"
)

type TMsgHandler func(path string, payload interface{}) (interface{}, error)

type TClient struct {
	clientId    string
	client      mqtt.Client
	qos         byte
	timeout     time.Duration
	callQueue   map[string]chan *TMsg
	reqHandler  TMsgHandler
	tchanPrefix string
}

func parseOpts(config conf.MqttConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.SetConnectTimeout(time.Duration(config.ConnectTimeoutMS) * time.Millisecond)
	opts.AddBroker("tcp://" + config.BrokerAddr)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetCleanSession(false)
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

func NewTClient(c conf.MqttConfig, clientId string, tchanPrefix string) *TClient {
	mqttOpts := parseOpts(c)
	mqttOpts.ClientID = clientId
	return &TClient{
		clientId:    clientId,
		client:      mqtt.NewClient(mqttOpts),
		qos:         c.QoS,
		timeout:     time.Duration(c.RequestTimeoutMS) * time.Millisecond,
		callQueue:   make(map[string]chan *TMsg),
		tchanPrefix: tchanPrefix,
	}
}

func (t *TClient) ClientId() string {
	return t.clientId
}

func (t *TClient) SetReqHandler(h TMsgHandler) {
	t.reqHandler = h
}

func (t *TClient) buildChanTopic(service string, serviceFunction string) string {
	return t.tchanPrefix + fmt.Sprintf("/%s/%s", service, serviceFunction)
}

func (t *TClient) Connect() error {
	tk := t.client.Connect()
	tk.Wait()
	if err := tk.Error(); err != nil {
		return err
	}

	channelTopic := t.buildChanTopic(t.clientId, "#")
	log.Infof("start to sub %s", channelTopic)

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
		log.Errorf("fail to unmarshal payload, ignore msg: %+v", msg)
		return
	}

	if tmsg.Method == REQ {
		log.Infof("receive req, tmsg: %+v", tmsg)

		req := tmsg.Payload
		if t.clientId != tmsg.KeyPart(KPI_Service) {
			log.Errorf("shouldn't happen, client:%s != receiver:%s, ignore msg: %+v", t.clientId, tmsg.KeyPart(KPI_Service), msg)
			return
		}

		function := tmsg.KeyPart(KPI_Function)
		if function == "" {
			log.Errorf("invalid req key %s, ignore msg: %+v", tmsg.Key, msg)
			return
		}

		reply := &TMsg{ID: tmsg.ID, Key: t.buildChanTopic(tmsg.Sender, function)}
		rsp, err := t.reqHandler(function, req)
		if err != nil {
			reply.Method = ERR
			reply.Payload = err.Error()
		} else {
			reply.Method = RSP
			reply.Payload = rsp
		}
		if _, err = t.Send(reply); err != nil {
			log.Errorf("invalid req key %s, ignore msg: %+v", tmsg.Key, msg)
		}
		return
	}

	if tmsg.Method == RSP || tmsg.Method == ERR {
		log.Infof("receive rsp, tmsg: %+v", tmsg)
		ch, ok := t.callQueue[tmsg.ID]
		if !ok {
			log.Errorf("fail to find corresponding REQ msg, ignore msg: %+v", msg)
			return
		}
		ch <- tmsg
		return
	}

	log.Errorf("shouldn't happen, invalid method %s, ignore msg: %+v", tmsg.Method, msg)
}

func (t *TClient) Send(tmsg *TMsg) (chan *TMsg, error) {

	tmsg.Sender = t.clientId

	var rspCh chan *TMsg
	if tmsg.Method == REQ {
		rspCh = make(chan *TMsg)
	}

	t.callQueue[tmsg.ID] = rspCh
	data, err := json.Marshal(tmsg)
	if err != nil {
		return nil, err
	}

	log.Infof("send msg %s: %s", tmsg.Key, data)
	tk := t.client.Publish(tmsg.Key, t.qos, false, data)
	tk.Wait()
	if err := tk.Error(); err != nil {
		return nil, err
	}
	return rspCh, nil

}

func (t *TClient) Watch(key string) (chan *TMsg, error) {
	log.Infof("start to watch %s", key)
	push := make(chan *TMsg)
	handler := func(cli mqtt.Client, msg mqtt.Message) {
		tmsg := &TMsg{}
		err := json.Unmarshal(msg.Payload(), tmsg)
		if err != nil {
			log.Errorf("fail to unmarshal payload, ignore msg: %+v", msg)
			return
		}

		if tmsg.Method != PUB {
			log.Errorf("shouldn't happen, invalid method %s, ignore msg: %+v", tmsg.Method, msg)
		}
		push <- tmsg
	}
	tk := t.client.Subscribe(key, t.qos, handler)
	tk.Wait()
	if err := tk.Error(); err != nil {
		return nil, err
	}
	return push, nil
}
