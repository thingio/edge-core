package talk

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/talk"
	"time"
)

type TMsgHandler func(path string, payload []byte) ([]byte, error)

type TClient struct {
	clientId    string
	client      mqtt.Client
	qos         byte
	timeout     time.Duration
	callQueue   map[string]chan *talk.TMessage
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
	log.Infof("mqttclient '%s' connected to %s", opts.ClientID(), opts.Servers()[0].String())
	return
}

func onConnectionLost(c mqtt.Client, err error) {
	opts := c.OptionsReader()
	log.WithError(err).Errorf("mqttclient '%s' connection to %s has lost, will try to reconnect", opts.ClientID(), opts.Servers()[0].String())
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
		callQueue:   make(map[string]chan *talk.TMessage),
		tchanPrefix: tchanPrefix,
	}
}

func (t *TClient) ClientId() string {
	return t.clientId
}

func (t *TClient) SetReqHandler(h TMsgHandler) {
	t.reqHandler = h
}

func (t *TClient) buildChatTopic(service string, serviceFunction string) string {
	return t.tchanPrefix + fmt.Sprintf("/%s/%s", service, serviceFunction)
}

func (t *TClient) Connect() error {
	tk := t.client.Connect()
	tk.Wait()
	if err := tk.Error(); err != nil {
		return err
	}

	channelTopic := t.buildChatTopic(t.clientId, "#")
	log.Infof("start to sub %s", channelTopic)

	tk = t.client.Subscribe(channelTopic, t.qos, t.callHandler)
	tk.Wait()
	if err := tk.Error(); err != nil {
		return err
	}
	return nil
}

func (t *TClient) callHandler(cli mqtt.Client, msg mqtt.Message) {
	tmsg := &talk.TMessage{}
	err := proto.Unmarshal(msg.Payload(), tmsg)
	if err != nil {
		log.Errorf("fail to unmarshal payload, ignore msg: %+v", msg)
		return
	}

	switch tmsg.Method {
	case talk.MethodREQ:
		log.Infof("receive req, tmsg: %+v", tmsg)
		req := tmsg.Payload
		if t.clientId != tmsg.ServiceId() {
			log.Errorf("shouldn't happen, client:%s != receiver:%s, ignore msg: %+v", t.clientId, tmsg.ServiceId(), msg)
			return
		}

		function := tmsg.ServiceFunction()
		if function == "" {
			log.Errorf("invalid req key %s, ignore msg: %+v", tmsg.Key, msg)
			return
		}

		reply := &talk.TMessage{Id: tmsg.Id, Key: t.buildChatTopic(tmsg.Sender, function)}
		rsp, err := t.reqHandler(function, req)
		if err != nil {
			reply.Method = talk.MethodERR
			reply.Payload = []byte(err.Error())
		} else {
			reply.Method = talk.MethodRSP
			reply.Payload = rsp
		}
		if _, err = t.Send(reply); err != nil {
			log.Errorf("invalid req key %s, ignore msg: %+v", tmsg.Key, msg)
		}
		return

	case talk.MethodRSP, talk.MethodERR:
		log.Infof("receive rsp, tmsg: %+v", tmsg)
		ch, ok := t.callQueue[tmsg.Id]
		if !ok {
			log.Errorf("fail to find corresponding MethodREQ msg, ignore msg: %+v", msg)
			return
		}
		ch <- tmsg
		return
	}

	log.Errorf("shouldn't happen, invalid method %s, ignore msg: %+v", tmsg.Method, msg)
}

func (t *TClient) Send(tmsg *talk.TMessage) (*talk.TMessage, error) {

	if tmsg.Method == talk.MethodREQ {
		rspCh := make(chan *talk.TMessage)
		defer close(rspCh)

		t.callQueue[tmsg.Id] = rspCh
		defer delete(t.callQueue, tmsg.Id)

		err := t.send(tmsg)
		if err != nil {
			return nil, err
		}

		select {
		case res := <-rspCh:
			return res, nil
		case <-time.After(t.timeout):
			log.Errorf("sending tmsg timeout: {%+v}", tmsg)
			return nil, errors.Timeoutf("tmessage %+v", tmsg)
		}
	}

	return nil, t.send(tmsg)
}

func (t *TClient) send(tmsg *talk.TMessage) error {

	log.Infof("sending tmsg: {%+v}", tmsg)
	tmsg.Sender = t.clientId

	data, err := proto.Marshal(tmsg)
	if err != nil {
		log.Errorf("marshal tmsg failed: {%+v}", tmsg)
		return err
	}

	tk := t.client.Publish(tmsg.Key, t.qos, false, data)
	tk.Wait()
	if err := tk.Error(); err != nil {
		log.Errorf("send tmsg failed: %+v", tmsg)
		return err
	}

	return nil
}

func (t *TClient) Watch(key string) (chan *talk.TMessage, error) {
	log.Infof("start to watch %s", key)
	push := make(chan *talk.TMessage)
	handler := func(cli mqtt.Client, msg mqtt.Message) {
		tmsg := &talk.TMessage{}
		err := proto.Unmarshal(msg.Payload(), tmsg)
		if err != nil {
			log.Errorf("fail to unmarshal payload, ignore msg: %+v", msg)
			return
		}

		if tmsg.Method != talk.MethodPUB {
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
