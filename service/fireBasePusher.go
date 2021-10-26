package service

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/utils"
)

type fireBasePusherImpl struct {
	name   string
	client *messaging.Client
	admin  *fireBaseAdmin
}

func NewFireBasePusher(name ...string) *fireBasePusherImpl {
	name = append(name, "")
	var impl = new(fireBasePusherImpl)
	impl.name = name[0]
	return impl
}

// SendAll 批量发送
func (pusher *fireBasePusherImpl) SendAll(messages []interface{}) (*messaging.BatchResponse, error) {
	if len(messages) <= 0 {
		return nil, errors.New("message empty error")
	}
	var agent, err = pusher.AcquireAgent()
	if err != nil {
		return nil, err
	}
	var all = pusher.transform(messages)
	if len(all) <= 0 {
		return nil, errors.New("message types error")
	}
	if utils.GetEnvBool("APP_DEBUG") {
		var content = utils.JsonEncode(messages).String()
		logger.WithField("message", content).Infoln("push all params")
	}
	return agent.SendAll(context.Background(), all)
}

// SendBatch 批量发送
func (pusher *fireBasePusherImpl) SendBatch(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	if len(messages) <= 0 {
		return nil, errors.New("message empty error")
	}
	var agent, err = pusher.AcquireAgent()
	if err != nil {
		return nil, err
	}
	if utils.GetEnvBool("APP_DEBUG") {
		var content = utils.JsonEncode(messages).String()
		logger.WithField("messages", content).Infoln("pushAll.params")
	}
	return agent.SendAll(context.Background(), messages)
}

// SendAllDryRun  批量验证发送可行性
func (pusher *fireBasePusherImpl) SendAllDryRun(messages []interface{}) (*messaging.BatchResponse, error) {
	if len(messages) <= 0 {
		return nil, errors.New("message empty error")
	}
	var agent, err = pusher.AcquireAgent()
	if err != nil {
		return nil, err
	}
	var all = pusher.transform(messages)
	if len(all) <= 0 {
		return nil, errors.New("message types error")
	}
	return agent.SendAllDryRun(context.Background(), all)
}

// Send 单个发送
func (pusher *fireBasePusherImpl) Send(message *messaging.Message) (string, error) {
	if message == nil {
		return "", errors.New("nil message params")
	}
	var agent, err = pusher.AcquireAgent()
	if err != nil {
		return "", err
	}
	if utils.GetEnvBool("APP_DEBUG") {
		bytes, _ := message.MarshalJSON()
		logger.WithField("message", string(bytes)).Infoln("push params")
	}
	return agent.Send(context.Background(), message)
}

// VerifyToken 验证token 有效性
func (pusher *fireBasePusherImpl) VerifyToken(idToken string) (*auth.Token, error) {
	var admin, err = pusher.AcquireAdmin()
	if err != nil {
		return nil, err
	}
	return admin.VerifyIDToken(idToken)
}

// 消息格式转换器
func (pusher *fireBasePusherImpl) transform(arr []interface{}) []*messaging.Message {
	var msgArr []*messaging.Message
	for _, v := range arr {
		switch v.(type) {
		case *messaging.Message:
			msgArr = append(msgArr, v.(*messaging.Message))
		case []byte:
			var msg = pusher.decodeMsg(v.([]byte))
			if msg != nil {
				msgArr = append(msgArr, msg)
			}
		case string:
			var msg = pusher.decodeMsg([]byte(v.(string)))
			if msg != nil {
				msgArr = append(msgArr, msg)
			}
		case map[string]interface{}:
			var msg = pusher.decodeMsgMap(v.(map[string]interface{}))
			if msg != nil {
				msgArr = append(msgArr, msg)
			}
		case fmt.Stringer:
			var msg = pusher.decodeMsg([]byte(v.(fmt.Stringer).String()))
			if msg != nil {
				msgArr = append(msgArr, msg)
			}
		case messaging.Message:
			var msg = v.(messaging.Message)
			msgArr = append(msgArr, &msg)
		}
	}
	return msgArr
}

// AcquireAgent 获取客户端
func (pusher *fireBasePusherImpl) AcquireAgent() (*messaging.Client, error) {
	if pusher.client != nil {
		return pusher.client, nil
	}
	var app, err = pusher.AcquireAdmin()
	if err != nil {
		return nil, err
	}
	agent, err := app.GetPusher()
	if err == nil {
		pusher.client = agent
	}
	return agent, err
}

// AcquireAdmin 获取admin App
func (pusher *fireBasePusherImpl) AcquireAdmin() (*fireBaseAdmin, error) {
	var err error
	if pusher.admin != nil {
		return pusher.admin, err
	}
	pusher.admin, err = pusher.createAdminApp()
	return pusher.admin, err
}

// 创建firebase App
func (pusher *fireBasePusherImpl) createAdminApp() (*fireBaseAdmin, error) {
	var name = pusher.name
	if name == "" {
		return NewFireBaseAdmin(pusher.GetCfgByAppName()), nil
	}
	return NewFireBaseAdmin(pusher.GetCfgByAppName(name)), nil
}

func (pusher *fireBasePusherImpl) GetCfgByAppName(name ...string) entity.KvMap {
	if len(name) <= 0 || name[0] == "" {
		return entity.KvMap{
			FirebaseCfgKey:     utils.GetEnvVal(FirebaseCfgEnvKey),
			FirebaseCredCfgKey: utils.GetEnvVal(FirebaseCredentialsCfgEnvKey),
		}
	}
	var prefix = name[0]
	return entity.KvMap{
		FirebaseCfgKey:     utils.GetEnvVal(pusher.prefixEnvKey(prefix, FirebaseCfgEnvKey)),
		FirebaseCredCfgKey: utils.GetEnvVal(pusher.prefixEnvKey(prefix, FirebaseCredentialsCfgEnvKey)),
	}
}

// 字节序转 消息体
func (pusher *fireBasePusherImpl) decodeMsg(data []byte) *messaging.Message {
	var msg = new(messaging.Message)
	if len(data) <= 0 {
		return nil
	}
	if err := utils.JsonDecode(data, msg); err != nil {
		return nil
	}
	return msg
}

func (pusher *fireBasePusherImpl) prefixEnvKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", prefix, key)
}

// map 转 消息体
func (pusher *fireBasePusherImpl) decodeMsgMap(data map[string]interface{}) *messaging.Message {
	var msg = new(messaging.Message)
	if len(data) <= 0 {
		return nil
	}
	var (
		jsonData = utils.JsonEncode(data)
		bytes    = jsonData.Bytes()
	)
	if len(bytes) < 0 || jsonData.HasErr() {
		return nil
	}
	if err := utils.JsonDecode(bytes, msg); err != nil {
		return nil
	}
	return msg
}
