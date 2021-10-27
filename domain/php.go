package domain

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/utils"
	"github.com/yookoala/gofast"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type PHPFastCgiDomainImpl struct {
	params      entity.KvMap
	ctx         context.Context
	cancel      context.CancelFunc
	connFactory gofast.ConnFactory
	clientPool  *gofast.ClientPool
	caller      gofast.Handler
	logger      *logrus.Logger
	timeout     time.Duration
	addr        string
	typeClass   entity.FastCgiType
}

var (
	defaultTimeout      = 3 * time.Minute
	defaultFastcgiRoot  = "/var/www/html"
	defaultFastcgiIndex = "/var/www/html/index.php"
)

const (
	ParamFastcgiRoot    = "root"
	ParamFastcgiPass    = "fastcgi_pass"
	ParamFastcgiFile    = "fastcgi_file"
	ParamFastcgiName    = "fastcgi_stream"
	ParamFastcgiLogFile = "fastcgi_log"
	defaultNetwork      = "tcp"
	PHPFastCGIType      = entity.FastCgiType("PHP-FastCGI")
)

func NewPHPFastCgiDomain() *PHPFastCgiDomainImpl {
	var domain = new(PHPFastCgiDomainImpl)
	domain.init()
	return domain
}

func (domain *PHPFastCgiDomainImpl) init() {
	domain.params = entity.KvMap{}
	domain.typeClass = PHPFastCGIType
	domain.ctx, domain.cancel = context.WithTimeout(context.Background(), defaultTimeout)
}

func (domain *PHPFastCgiDomainImpl) Parsed() bool {
	if domain.caller != nil && domain.addr != "" {
		return true
	}
	return false
}

func (domain *PHPFastCgiDomainImpl) Parse(properties []byte) error {
	if len(properties) <= 0 {
		return errors.New("empty properties")
	}
	if err := utils.JsonDecode(properties, &domain.params); err != nil {
		return err
	}
	var addr = domain.params.GetStr(ParamFastcgiPass)
	if addr == "" {
		return errors.New("miss param: " + ParamFastcgiPass)
	}
	var (
		root     = domain.params.GetStr(ParamFastcgiRoot, defaultFastcgiRoot)
		endpoint = domain.params.GetStr(ParamFastcgiFile, defaultFastcgiIndex)
	)
	if d := domain.params.GetDuration("fastcgi_timeout", 0); d > 0 {
		domain.SetTimeout(d)
	}
	if root == "" {
		if endpoint == "" {
			return errors.New("miss param: " + ParamFastcgiRoot)
		}
		// 1. conn addr
		domain.connFactory = gofast.SimpleConnFactory(defaultNetwork, addr)
		// domain.clientPool = gofast.NewClientPool(gofast.SimpleClientFactory(domain.connFactory), 3, domain.timeout)
		// 2. root file
		domain.caller = gofast.NewHandler(
			gofast.NewFileEndpoint(endpoint)(gofast.BasicSession),
			gofast.SimpleClientFactory(domain.connFactory),
		)
	} else {
		// 1. conn addr
		domain.connFactory = gofast.SimpleConnFactory(defaultNetwork, addr)
		// domain.clientPool = gofast.NewClientPool(gofast.SimpleClientFactory(domain.connFactory), 3, domain.timeout)
		// 2. root path dir
		domain.caller = gofast.NewHandler(
			gofast.NewPHPFS(root)(gofast.BasicSession),
			gofast.SimpleClientFactory(domain.connFactory),
		)
	}
	domain.caller.SetLogger(domain.getLogger())
	// 设置超时时长

	return nil
}

func (domain *PHPFastCgiDomainImpl) reset() *PHPFastCgiDomainImpl {
	domain.ctx, domain.cancel = context.WithTimeout(context.Background(), defaultTimeout)
	return domain
}

func (domain *PHPFastCgiDomainImpl) Register(pool *sync.Pool) {
	if pool == nil {
		return
	}
	pool.Put(domain.reset())
}

func (domain *PHPFastCgiDomainImpl) Cancel() {
	domain.cancel()
}

func (domain *PHPFastCgiDomainImpl) Type() string {
	if domain.typeClass == "" {
		return PHPFastCGIType.String()
	}
	return domain.typeClass.String()
}

func (domain *PHPFastCgiDomainImpl) Call(ctx *fiber.Ctx) error {
	if !domain.Parsed() {
		return errors.New("not parsed fastcgi handler")
	}
	var res, req, err = utils.CreateHttpNative(ctx)
	if err != nil {
		return err
	}
	domain.caller.ServeHTTP(res, req)
	return err
}

func (domain *PHPFastCgiDomainImpl) Proxy(res http.ResponseWriter, req *http.Request) error {
	if !domain.Parsed() {
		return errors.New("not parsed fastcgi handler")
	}
	domain.caller.ServeHTTP(res, req)
	return nil
}

func (domain *PHPFastCgiDomainImpl) getLogger() *log.Logger {
	if domain.logger != nil {
		return log.New(domain.logger.Out, domain.getName(), log.LstdFlags)
	}
	if file := domain.params.GetStr(ParamFastcgiLogFile, ""); file != "" {
		if fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm); err == nil {
			return log.New(fd, domain.getName(), log.LstdFlags)
		}
	}
	return log.Default()
}

func (domain *PHPFastCgiDomainImpl) getName() string {
	if domain.params != nil {
		return domain.params.GetStr(ParamFastcgiName, domain.addr)
	}
	return domain.addr
}

func (domain *PHPFastCgiDomainImpl) SetTimeout(duration time.Duration) {
	domain.timeout = duration
	domain.ctx, domain.cancel = context.WithTimeout(context.Background(), duration)
}
