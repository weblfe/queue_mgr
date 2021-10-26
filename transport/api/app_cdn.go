package api

import (
	"github.com/gofiber/fiber/v2"
)

// AppCdnApi app内容分发服务接口集合
type AppCdnApi interface {

	// CreateApp godoc
	// @Summary 创建(注册)app数据
	// @Tags AppCdnServ
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
	// @Router /app/create [post]
	CreateApp(ctx *fiber.Ctx) error

	// UpdateApp godoc
	// @Summary 更新直播应用数据
	// @Tags AppCdnServ
	// @Description update live app info
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
	// @Router /app/update [put]
	UpdateApp(ctx *fiber.Ctx) error

	// SaveLive godoc
	// @Summary 保存开播数据
	// @Tags AppCdnServ
	// @Description save opened live room info
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "应用ID"
	// @Param uid formData string true "开播用户id"
	// @Param stream formData string true "stream name/live room number"
	// @Param title formData string true "直播标题"
	// @Param thumb formData string true "直播封面"
	// @Param user_nicename formData string true "主播昵称"
	// @Param level formData string false "用户等级"
	// @Param anchor_level formData string false "主播等级"
	// @Param avatar formData string false "主播头像"
	// @Param bongo_id formData string false "主播openid"
	// @Param country formData string false "live location country name"
	// @Param province formData string false "live location province name"
	// @Param city formData string false "live location city name"
	// @Param pull formData string false "live video stream url"
	// @Param lng formData string false "直播地理位置经度"
	// @Param lat formData string false "直播地理位置纬度"
	// @Param type formData string true "live type"
	// @Param anyway formData string true "横竖屏: 0表示竖屏, 1表示横屏" Enums(0,1)
	// @Param liveclassid formData string true "直播分类"
	// @Param hotvotes formData int true "直播间热门礼物金豆数"
	// @Param pkuid formData int false "pk用户UID"
	// @Param pkstream formData string false "pk用户stream"
	// @Param isvideo formData int true "是否录屏"
	// @Param ismic formData int true "连麦开关, 0:关, 1:开" Enums(0,1)
	// @Param ishot formData int true "是否热门, 0:否, 1:热门" Enums(0,1)
	// @Param isrecommend formData int true "是否推荐, 0:不推荐, 1:推荐" Enums(0,1)
	// @Param isshop formData int true "是否开启店铺, 0:否, 1:是" Enums(0,1)
	// @Param islive formData int true "直播状态,0:未开始直播,1:直播中,2:暂停" Enums(0,1,2)
	// @Param isoff formData int false "是否断流,0:否,1:是" Enums(0,1)
	// @Param offtime formData int false "断流时间戳"
	// @Param nums formData int true "直播间人数"
	// @Param starttime formData int true "直播开始时间"
	// @Param endtime formData int false "直播结束时间"
	// @Param deviceinfo formData string true "设备信息"
	// @Param score formData int false "排序值/分数"
	// @Param match_tag_hash formData string false "匹配直播标签hash"
	// @Param is_only_audit formData int false "是否仅审核环境使用,0:否,1:是" Enums(0,1)
	// @Param vision_score formData float64 false "视觉评分 (1000~0.001)[0.001-0.999:睡播或卡播 ,1.00-100.000:绿播,100.001~300.99:黄播,400.00~600.99:xx播,>=700非人直播 ]"
	// @Param banker_coin formData int true "庄家余额"
	// @Param game_action formData int false "游戏类型"
	// @Param coin_total formData int false "当前直播收到钻石数"
	// @Param rtc_token formData string false "声网rtc_token(默认有效12小时)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /save/live [post]
	SaveLive(ctx *fiber.Ctx) error

	// SetLiveState godoc
	// @Summary 通知直播开关播状态
	// @Tags AppCdnServ
	// @Description notify live close state
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "应用ID"
	// @Param stream formData string true "live stream|live room number"
	// @Param uid formData int true "关播用户UID"
	// @Param timestamp formData int false "关播时间(通知开播时为0,关闭时必须>0)"
	// @Param starttime formData int true "开播时间"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /live/state [post]
	SetLiveState(ctx *fiber.Ctx) error

	// SaveAnchorState godoc
	// @Summary 保存同步各平台主播签约状态
	// @Tags AppCdnServ
	// @Description  save anchor info (family_id,uid,status)
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "应用ID"
	// @Param family_id formData int true "公会ID"
	// @Param uid formData int true "用户ID"
	// @Param status formData int true "签约状态,1:签约,2:解约" Enums(1,2)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /save/anchor [post]
	SaveAnchorState(ctx *fiber.Ctx) error

	// GetHotLives godoc
	// @Summary 获取热门主播列表
	// @Tags AppCdnServ
	// @Description get popular lives lists
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param AppID header string false "app id"
	// @Param audit_channel header string false "audit channel"
	// @Param Version header string false "app version"
	// @Param Channel header string false "app package channel"
	// @Param Platform header string false "app device os platform"
	// @Param appID path string true "当前应用 appID"
	// @Param ios_version formData string false "ios version/ios版本号" default(1.0.0)
	// @Param page formData integer false "page number/页码" default(1)
	// @Param count formData integer false "page size number/分页量" default(10)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /{appID}/lives/hot [post]
	GetHotLives(ctx *fiber.Ctx) error

	// GetRecommendLives godoc
	// @Summary 获取推荐主播列表
	// @Tags AppCdnServ
	// @Description get recommend lives lists
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param AppID header string false "app id"
	// @Param audit_channel header string false "audit channel"
	// @Param Version header string false "app version"
	// @Param Channel header string false "app package channel"
	// @Param Platform header string false "app device os platform"
	// @Param appID path string true "当前应用 appID"
	// @Param ios_version formData string false "ios version/ios版本号" default(1.0.0)
	// @Param page formData integer false "page number/页码" default(1)
	// @Param count formData integer false "page size number/分页量" default(10)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /{appID}/lives/recommend [post]
	GetRecommendLives(ctx *fiber.Ctx) error

	// GetNewLives godoc
	// @Summary 获取最近开播直播间列表
	// @Tags AppCdnServ
	// @Description get last opened live lists
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param AppID header string false "your app id"
	// @Param audit_channel header string false "audit channel"
	// @Param Version header string false "app version"
	// @Param Channel header string false "app package channel"
	// @Param Platform header string false "app device os platform"
	// @Param appID path string true "当前应用 appID"
	// @Param ios_version formData string false "ios version/ios版本号" default(1.0.0)
	// @Param lng formData string false "经度值"
	// @Param lat formData string false "纬度值"
	// @Param page formData integer false "page number/页码" default(1)
	// @Param count formData integer false "page size number/分页量" default(10)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /{appID}/lives/new [post]
	GetNewLives(ctx *fiber.Ctx) error

	// SaveRecommendLives godoc
	// @Summary 保存推荐直播列表
	// @Tags AppCdnServ
	// @Description save recommend lives info
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param app_id formData string true "应用ID"
	// @Param uid formData string true "开播用户id"
	// @Param stream formData string true "stream name/live room number"
	// @Param title formData string true "直播标题"
	// @Param thumb formData string true "直播封面"
	// @Param user_nicename formData string true "主播昵称"
	// @Param avatar formData string false "主播头像"
	// @Param bongo_id formData string false "主播openid"
	// @Param level formData string false "用户等级"
	// @Param anchor_level formData string false "主播等级"
	// @Param country formData string false "live location country name"
	// @Param province formData string false "live location province name"
	// @Param city formData string false "live location city name"
	// @Param pull formData string false "live video stream url"
	// @Param lng formData string false "直播地理位置经度"
	// @Param lat formData string false "直播地理位置纬度"
	// @Param type formData string true "live type"
	// @Param anyway formData string true "横竖屏: 0表示竖屏, 1表示横屏" Enums(0,1)
	// @Param liveclassid formData string true "直播分类"
	// @Param hotvotes formData int true "直播间热门礼物金豆数"
	// @Param pkuid formData int false "pk用户UID"
	// @Param pkstream formData string false "pk用户stream"
	// @Param isvideo formData int true "是否录屏"
	// @Param ismic formData int true "连麦开关, 0:关, 1:开" Enums(0,1)
	// @Param ishot formData int true "是否热门, 0:否, 1:热门" Enums(0,1)
	// @Param isrecommend formData int true "是否推荐, 0:不推荐, 1:推荐" Enums(0,1)
	// @Param isshop formData int true "是否开启店铺, 0:否, 1:是" Enums(0,1)
	// @Param islive formData int true "直播状态,0:未开始直播,1:直播中,2:暂停" Enums(0,1,2)
	// @Param isoff formData int false "是否断流,0:否,1:是" Enums(0,1)
	// @Param offtime formData int false "断流时间戳"
	// @Param nums formData int true "直播间人数"
	// @Param starttime formData int true "直播开始时间"
	// @Param endtime formData int false "直播结束时间"
	// @Param deviceinfo formData string false "设备信息"
	// @Param score formData int false "排序值/分数"
	// @Param match_tag_hash formData string false "匹配直播标签hash"
	// @Param is_only_audit formData int false "是否仅审核环境使用,0:否,1:是" Enums(0,1)
	// @Param vision_score formData float64 false "视觉评分 (1000~0.001)[0.001-0.999:睡播或卡播 ,1.00-100.000:绿播,100.001~300.99:黄播,400.00~600.99:xx播,>=700非人直播 ]"
	// @Param banker_coin formData int true "庄家余额"
	// @Param game_action formData int false "游戏类型"
	// @Param coin_total formData int false "当前直播收到钻石数"
	// @Param recommend_position formData int true "推荐页码位置"
	// @Param rtc_token formData string true "声网rtc_token(默认有效12小时)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /save/recommends [post]
	SaveRecommendLives(ctx *fiber.Ctx) error
}
