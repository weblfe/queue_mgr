package api

import (
	"github.com/gofiber/fiber/v2"
)

// QueueManagerApi app队列管理服务接口集合
type QueueManagerApi interface {

	// CreateConsumer godoc
	// @Summary 创建队列消费器
	// @Tags QueueMgrServ
	// @Description register live app info
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_name formData string true "app name"
	// @Param app_logo_url formData string false "app logo url"
	// @Param callback_url formData string false "callback url"
	// @Param sync_url formData string false "sync lives api"
	// @Param proxy_url formData string false "proxy url addr"
	// @Param sync_type formData int false "sync type ( 1:post_sync,2:interval_sync )" Enums(1, 2) default(1)
	// @Param interval_duration formData int false "interval sync duration (unit: second)" Enums(0,10,15,30,60,180)
	// @Param interval_fail_times formData int false "interval fail retry times" Enums(0,3,5,10) default(0)
	// @Param extras formData string false "extras data " default({"pro_im_url":"","test_im_url":"","dev_im_url":""})
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /consumer/create [post]
	CreateConsumer(ctx *fiber.Ctx) error

	// State godoc
	// @Summary 查询队列消费器状态
	// @Tags QueueMgrServ
	// @Description query queue consumer state
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "app id"
	// @Param app_name formData string false "app name"
	// @Param app_logo_url formData string false "app logo url"
	// @Param callback_url formData string false "callback url"
	// @Param proxy_url formData string false "proxy url addr"
	// @Param sync_url formData string false "sync lives api"
	// @Param sync_type formData int false "sync type ( 1:post_sync,2:interval_sync )" Enums(1, 2)
	// @Param interval_duration formData int false "interval sync duration (unit: second)" Enums(10,15,30,60,180)
	// @Param interval_fail_times formData int false "interval fail retry times" Enums(0,3,5,10)
	// @Param extras formData string false "extras data" default({"pro_im_url":"","test_im_url":"","dev_im_url":""})
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /state [get]
	State(ctx *fiber.Ctx) error


	// Control godoc
	// @Summary 控制队列消费器状态
	// @Tags QueueMgrServ
	// @Description change queue consumer state
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "app id"
	// @Param app_name formData string false "app name"
	// @Param app_logo_url formData string false "app logo url"
	// @Param callback_url formData string false "callback url"
	// @Param proxy_url formData string false "proxy url addr"
	// @Param sync_url formData string false "sync lives api"
	// @Param sync_type formData int false "sync type ( 1:post_sync,2:interval_sync )" Enums(1, 2)
	// @Param interval_duration formData int false "interval sync duration (unit: second)" Enums(10,15,30,60,180)
	// @Param interval_fail_times formData int false "interval fail retry times" Enums(0,3,5,10)
	// @Param extras formData string false "extras data" default({"pro_im_url":"","test_im_url":"","dev_im_url":""})
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /state/update [post]
	Control(ctx *fiber.Ctx) error

}
