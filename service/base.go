package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/config"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/repo"
	"github.com/weblfe/queue_mgr/utils"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// 服务
type baseServiceImpl struct {
	serviceID   string
	addr        string
	serviceType entity.ServiceType
	desc        string
	host        string
	constructor sync.Once
}

var (
	NoError = errors.New("")
	logger  *logrus.Logger
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47"

func (service *baseServiceImpl) Type() entity.ServiceType {
	return service.serviceType
}

func (service *baseServiceImpl) setType(t entity.ServiceType) {
	service.serviceType = t
}

func (service *baseServiceImpl) initService(initialize func()) {
	service.constructor.Do(initialize)
}

func (service *baseServiceImpl) ServiceID() string {
	return service.serviceID
}

func (service *baseServiceImpl) setServiceID(id string) {
	if service.serviceID != "" || id == "" {
		return
	}
	service.serviceID = id
}

func (service *baseServiceImpl) Addr() string {
	return service.addr
}

func (service *baseServiceImpl) setAddr(addr string) {
	if service.addr != "" || addr == "" {
		return
	}
	service.addr = addr
}

func (service *baseServiceImpl) GetClient(method ...string) *fiber.Agent {
	method = append(method, "")
	var serv = fiber.AcquireAgent()
	if method[0] != "" {
		addr := service.GetServAddr(method[0])
		serv.Request().SetRequestURI(addr)
	}
	serv.Add("Version", "0.1.0")
	serv.Add("Platform", "robot_gift")
	serv.Add("Channel", "robot")
	serv.Add("RequestId", GetRequestID(method[0]))
	serv.UserAgent(userAgent)
	return serv
}

// GetLogger 获取日志
func (service *baseServiceImpl)GetLogger() *logrus.Logger {
	if logger == nil {
		Boot()
	}
	return logger
}

func (service *baseServiceImpl) GetServAddr(method string) string {
	var (
		serv       = service.GetServConfigure()
		entrypoint = serv.GetEntryPoint()
	)
	if strings.HasSuffix(entrypoint, "/") {
		return serv.GetEntryPoint() + "?service=" + method
	}
	return serv.GetEntryPoint() + "/?service=" + method
}

func (service *baseServiceImpl) GetServConfigure() config.Service {
	return config.GetAppConfig().GetService(service.serviceID)
}

func (service *baseServiceImpl) Desc() string {
	return service.desc
}

func (service *baseServiceImpl) setDesc(desc string) {
	if service.desc != "" || desc == "" {
		return
	}
	service.desc = desc
}

func (service *baseServiceImpl) initAddr() {
	if service.addr == "" {
		service.addr = fmt.Sprintf("%s://%s/%s", service.Type().Desc(), service.getHostIp(), service.ServiceID())
	}
}

func (service *baseServiceImpl) PostForm(client *fiber.Agent, args *fiber.Args) error {
	var (
		req = client.Request()
	)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType(fiber.MIMEApplicationForm)
	if args != nil {
		if utils.GetEnvBool("APP_DEBUG") {
			logger.Infoln(fmt.Sprintf("url: %s, body: %s", req.URI(), string(args.QueryString())))
		}
		req.SetBodyRaw(args.QueryString())
	}
	defer fiber.ReleaseArgs(args)
	if err := client.Parse(); err != nil {
		return err
	}
	return nil
}

func (service *baseServiceImpl) GetUrl(client *fiber.Agent, args *fiber.Args) (statusCode int, body []byte, err error) {
	var (
		queryUrlValues     = url.Values{}
		req                = client.Request()
		baseUrl            = req.RequestURI()
		dst                = new([]byte)
		baseUrlObj, errUrl = url.ParseRequestURI(string(baseUrl))
	)
	if errUrl != nil {
		return 0, nil, errUrl
	}
	// url 本身 query
	if baseUrlObj.RawQuery != "" {
		if v, err2 := url.ParseQuery(baseUrlObj.RawQuery); err2 == nil {
			queryUrlValues = v
		}
	}
	if args != nil {
		defer fiber.ReleaseArgs(args)
		var (
			argsStr       = string(args.QueryString())
			argsValues, _ = url.ParseQuery(argsStr)
		)
		if len(argsValues) > 0 {
			for k, v := range argsValues {
				queryUrlValues.Set(k, v[0])
			}
		}
		baseUrlObj.RawQuery = queryUrlValues.Encode()
		if utils.GetEnvBool("APP_DEBUG") {
			logger.Infoln(fmt.Sprintf("request url: %s", baseUrlObj.String()))
		}

	}
	if err2 := client.Parse(); err2 != nil {
		return 0, nil, err2
	}
	return client.Get(*dst, baseUrlObj.String())
}

func (service *baseServiceImpl) getHostIp() string {
	if service.host != "" {
		return service.host
	}
	var ip = os.Getenv("EXPORT_HOST_IP")
	if ip == "" {
		if ip, err := getRealIp(); err == nil && ip != "" {
			service.host = ip
			return ip
		}
		if ip, err := getLocalIp(); err == nil && ip != "" {
			service.host = ip
			return ip
		}
	}
	if ip != "" {
		service.host = ip
	} else {
		ip = "127.0.0.1"
		service.host = ip
	}
	return ip
}

func getRealIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func getLocalIp() (ip string, err error) {
	ip = ""
	addrArr, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrArr {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	if ip == "" {
		return "", errors.New("not ip address")
	}
	return
}

// GetRequestID 生成请求ID
func GetRequestID(str ...string) string {
	str = append(str, "")
	var (
		md   = md5.New()
		data = fmt.Sprintf("%d:%d:$%s", time.Now().UnixNano(), rand.Int63(), str[0])
	)
	md.Write([]byte(data))
	return fmt.Sprintf("%x", md.Sum(nil))
}

func Boot() *logrus.Logger {
	if logger == nil {
		logger = repo.GetLogger("service")
	}
	return logger
}
