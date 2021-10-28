package entity

import (
	"github.com/gofiber/fiber/v2"
	"github.com/weblfe/queue_mgr/utils"
)

type (
	// QueueOptions 队列选项
	QueueOptions struct {
		Name         string `json:"name,omitempty" yaml:"name"`         // 队列名
		Type         string `json:"type,omitempty" yaml:"type"`         // 队列类型
		Handler      string `json:"handler" yaml:"handler"`             // 处理器
		ConnUrl      string `json:"conn_url" yaml:"conn_url"`           // amqp,mqtt 链接url
		Queue        string `json:"queue" yaml:"queue"`                 // 队列名
		Topic        string `json:"topic" yaml:"topic"`                 // 主题名
		ErrorHandler string `json:"error_handler" yaml:"error_handler"` // 异常处理器
	}

	QueueParams struct {
		AppID string `form:"appid" json:"appid,omitempty"`
		// 队列状态 0:未启动消费,1:消费中,2:idle(空闲),3:暂停消费
		Status uint `form:"status" json:"status"`
		// 类型 amqp,mqtt,redis,native
		Type string `form:"type" json:"type"`
		// 队列名
		Name string `form:"name" json:"name"`
		// 消费者数量
		ConsumerMaxNum uint `form:"consumer_max_num" json:"consumer_max_num"`
		// 队列配置json
		Properties string `form:"properties" json:"properties"`
		// 备注说明
		Comment string `form:"comment" json:"comment"`
	}

	ConsumerParams struct {
		AppID string `form:"appid" json:"appid,omitempty"`
		// 队列状态 0:未启动消费,1:消费中,2:idle(空闲),3:暂停消费
		Status uint `form:"status" json:"status"`
		// 类型 fastcgi,api,grpc,shell,ws,mysql
		Type string `form:"type" json:"type"`
		// 消费器名
		Name string `form:"name" json:"name"`
		// 消费器配置json
		Properties string `form:"properties" json:"properties"`
		// 备注说明
		Comment string `form:"comment" json:"comment"`
	}

	StateParams struct {
		AppID string `form:"appid" json:"appid,omitempty"`
		// 队列状态 0:未启动消费,1:消费中,2:idle(空闲),3:暂停消费
		Status uint `form:"status" json:"status"`
		// 消费器名
		Name string `form:"queue" json:"queue"`
	}

	BindParams struct {
		AppID string `form:"appid" json:"appid"`
		// 绑定的消费队列
		Queue string `form:"queue" json:"queue"`
		// 消费器
		Consumer string `form:"consumer" json:"consumer"`
		// 状态  1:绑定,2:解绑
		Status uint `form:"status" json:"status"`
		// 消费器配置 json
		Properties string `form:"properties" json:"properties"`
	}

	QueryParams struct {
		AppID    string `form:"appid" json:"appid"`
		Page     uint   `form:"page" json:"page,default=1"`
		Count    uint   `form:"count" json:"count,default=10"`
		Status   uint   `form:"status" json:"status,omitempty"`
		Queue    string `form:"queue" json:"queue,omitempty"`
		Consumer string `form:"consumer" json:"consumer,omitempty"`
	}
)

func NewQueueOption() *QueueOptions {
	var options = new(QueueOptions)
	return options
}

func DecoderOptions(bytes []byte) *QueueOptions {
	if len(bytes) <= 0 {
		return nil
	}
	var opts = NewQueueOption()
	if err := utils.JsonDecode(bytes, opts); err != nil {
		return nil
	}
	return opts
}

func (options *QueueOptions) GetTopic() string {
	return options.Topic
}

func (options *QueueOptions) GetQueue() string {
	return options.Queue
}

func (options *QueueOptions) GetType() string {
	return options.Type
}

func (params *QueueParams) Decode(data []byte) error {
	if err := utils.JsonDecode(data, params); err == nil {
		params.load()
		return nil
	}
	return utils.JsonDecode(data, params)
}

func (params *QueueParams) load() {
	if params.AppID == "" {
		params.AppID = utils.GetEnvVal("APP_ID")
	}
}

func (params *QueueParams) Parse(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		if err := ctx.BodyParser(params); err != nil {
			if err2 := params.Decode(ctx.Body()); err2 != nil {
				return err2
			}
			return err
		}
	} else {
		var bytes = argsToJsonBytes(ctx.Request().URI().QueryArgs())
		if len(bytes) > 0 {
			if err2 := params.Decode(argsToJsonBytes(ctx.Request().URI().QueryArgs())); err2 != nil {
				return err2
			}
		}
	}
	if params.AppID == "" {
		params.AppID = ctx.Params("appID", utils.GetEnvVal("APP_ID"))
	}
	return nil
}

func (params *StateParams) Decode(data []byte) error {
	if err := utils.JsonDecode(data, params); err == nil {
		params.load()
		return nil
	}
	return utils.JsonDecode(data, params)
}

func (params *StateParams) load() {
	if params.AppID == "" {
		params.AppID = utils.GetEnvVal("APP_ID")
	}
}

func (params *StateParams) Parse(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		if err := ctx.BodyParser(params); err != nil {
			if err2 := params.Decode(ctx.Body()); err2 != nil {
				return err2
			}
			return err
		}
	} else {
		var bytes = argsToJsonBytes(ctx.Request().URI().QueryArgs())
		if len(bytes) > 0 {
			if err2 := params.Decode(argsToJsonBytes(ctx.Request().URI().QueryArgs())); err2 != nil {
				return err2
			}
		}
	}
	if params.AppID == "" {
		params.AppID = ctx.Params("appID", utils.GetEnvVal("APP_ID"))
	}
	return nil
}

func (params *BindParams) Decode(data []byte) error {
	if err := utils.JsonDecode(data, params); err == nil {
		params.load()
		return nil
	}
	return utils.JsonDecode(data, params)
}

func (params *BindParams) load() {
	if params.AppID == "" {
		params.AppID = utils.GetEnvVal("APP_ID")
	}
}

func (params *BindParams) Parse(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		if err := ctx.BodyParser(params); err != nil {
			if err2 := params.Decode(ctx.Body()); err2 != nil {
				return err2
			}
			return err
		}
	} else {
		var bytes = argsToJsonBytes(ctx.Request().URI().QueryArgs())
		if len(bytes) > 0 {
			if err2 := params.Decode(argsToJsonBytes(ctx.Request().URI().QueryArgs())); err2 != nil {
				return err2
			}
		}
	}
	if params.AppID == "" {
		params.AppID = ctx.Params("appID", utils.GetEnvVal("APP_ID"))
	}
	return nil
}

func (params *ConsumerParams) Decode(data []byte) error {
	if err := utils.JsonDecode(data, params); err == nil {
		params.load()
		return nil
	}
	return utils.JsonDecode(data, params)
}

func (params *ConsumerParams) load() {
	if params.AppID == "" {
		params.AppID = utils.GetEnvVal("APP_ID")
	}
}

func (params *ConsumerParams) Parse(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		if err := ctx.BodyParser(params); err != nil {
			if err2 := params.Decode(ctx.Body()); err2 != nil {
				return err2
			}
			return err
		}
	} else {
		var bytes = argsToJsonBytes(ctx.Request().URI().QueryArgs())
		if len(bytes) > 0 {
			if err2 := params.Decode(argsToJsonBytes(ctx.Request().URI().QueryArgs())); err2 != nil {
				return err2
			}
		}
	}
	if params.AppID == "" {
		params.AppID = ctx.Params("appID", utils.GetEnvVal("APP_ID"))
	}
	return nil
}

func (params *QueryParams) Decode(data []byte) error {
	if err := utils.JsonDecode(data, params); err == nil {
		params.load()
		return nil
	}
	return utils.JsonDecode(data, params)
}

func (params *QueryParams) load() {
	if params.AppID == "" {
		params.AppID = utils.GetEnvVal("APP_ID")
	}
}

func (params *QueryParams) Parse(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		if err := ctx.BodyParser(params); err != nil {
			if err2 := params.Decode(ctx.Body()); err2 != nil {
				return err2
			}
			return err
		}
	} else {
		var bytes = argsToJsonBytes(ctx.Request().URI().QueryArgs())
		if len(bytes) > 0 {
			if err2 := params.Decode(argsToJsonBytes(ctx.Request().URI().QueryArgs())); err2 != nil {
				return err2
			}
		}
	}
	if params.AppID == "" {
		params.AppID = ctx.Params("appID", utils.GetEnvVal("APP_ID"))
	}
	return nil
}
