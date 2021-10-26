package service

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/utils"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fireBaseAdmin struct {
	config         entity.KvMap
	client         *firebase.App
	ctx            context.Context
	cancel         context.CancelFunc
	firebaseConfig *firebase.Config
}

const (
	FirebaseCfgEnvKey            = "FIREBASE_CONFIG"
	FirebaseCredentialsCfgEnvKey = "CREDENTIALS_CONFIG"
	FirebaseCfgKey               = "firebase_config"
	FirebaseCredCfgKey           = "credentials_config"
)

// NewFireBaseAdmin 创建 firebase Admin App
func NewFireBaseAdmin(cfg ...entity.KvMap) *fireBaseAdmin {
	var admin = new(fireBaseAdmin)
	if len(cfg) > 0 {
		admin.config = cfg[0]
	} else {
		admin.config = entity.KvMap{
			FirebaseCfgKey:     utils.GetEnvVal(FirebaseCfgEnvKey),
			FirebaseCredCfgKey: utils.GetEnvVal(FirebaseCredentialsCfgEnvKey),
		}
	}
	admin.ctx, admin.cancel = context.WithCancel(context.Background())
	return admin
}

func (admin *fireBaseAdmin) AcquireAgent() (*firebase.App, error) {
	var err error
	if admin == nil {
		return nil, errors.New("nil Admin object")
	}
	if admin.client != nil {
		return admin.client, nil
	}
	admin.client, err = firebase.NewApp(admin.ctx, admin.configure(), admin.options())
	if err != nil {
		return nil, err
	}
	return admin.client, nil
}

func (admin *fireBaseAdmin) VerifyIDToken(token string) (*auth.Token, error) {
	var client, err = admin.AcquireAgent()
	if err != nil {
		return nil, err
	}
	var authClient, err2 = client.Auth(admin.ctx)
	if err2 != nil {
		return nil, err2
	}
	return authClient.VerifyIDToken(admin.ctx, token)
}

func (admin *fireBaseAdmin) configure() *firebase.Config {
	if admin.firebaseConfig != nil {
		return admin.firebaseConfig
	}
	if file, ok := admin.configKv(FirebaseCfgKey); ok && file != "" && file != nil {
		switch file.(type) {
		case string:
			admin.firebaseConfig, _ = admin.parse(file.(string))
		case *firebase.Config:
			admin.firebaseConfig = file.(*firebase.Config)
		case firebase.Config:
			var cfg = file.(firebase.Config)
			admin.firebaseConfig = &cfg
		case map[string]interface{}:
			var m = file.(map[string]interface{})
			if len(m) > 0 {
				var (
					fcg   = &firebase.Config{}
					bytes = utils.JsonEncode(m).Bytes()
					err2  = utils.JsonDecode(bytes, fcg)
				)
				if err2 == nil {
					admin.firebaseConfig = fcg
				}
			}
		}
	}
	return admin.firebaseConfig
}

func (admin *fireBaseAdmin) parse(confFileName string) (*firebase.Config, error) {
	var fbc = &firebase.Config{}
	if confFileName == "" {
		return nil, nil
	}
	var dat []byte
	if confFileName[0] == byte('{') {
		dat = []byte(confFileName)
	} else {
		var err error
		if dat, err = ioutil.ReadFile(confFileName); err != nil {
			return nil, err
		}
	}
	if err := utils.JsonDecode(dat, fbc); err != nil {
		return nil, err
	}
	// Some special handling necessary for db auth overrides
	var m map[string]interface{}
	if err := utils.JsonDecode(dat, &m); err != nil {
		return nil, err
	}
	if ao, ok := m["databaseAuthVariableOverride"]; ok && ao == nil {
		// Auth overrides are explicitly set to null
		var nullMap map[string]interface{}
		fbc.AuthOverride = &nullMap
	}
	return fbc, nil
}

func (admin *fireBaseAdmin) configKv(key string) (interface{}, bool) {
	var v = admin.config.Get(key, nil)
	if v == nil {
		return v, false
	}
	return v, true
}

func (admin *fireBaseAdmin) options() option.ClientOption {
	var file, ok = admin.configKv(FirebaseCredCfgKey)
	if !ok || file == "" || file == nil {
		return nil
	}
	switch file.(type) {
	case option.ClientOption:
		return file.(option.ClientOption)
	case []option.ClientOption:
		var opts = file.([]option.ClientOption)
		if len(opts) <= 0 {
			return nil
		}
	case string:
		var str = file.(string)
		if json.Valid([]byte(str)) {
			return option.WithCredentialsJSON([]byte(str))
		}
		path:=admin.resolvePath(str)
		if _, err := os.Stat(path); err != nil {
			return nil
		}
		return option.WithCredentialsFile(path)
	case []byte:
		var bytes = file.([]byte)
		if json.Valid(bytes) {
			return option.WithCredentialsJSON(bytes)
		}
		var str = string(bytes)
		if _, err := os.Stat(str); err != nil {
			return nil
		}
		return option.WithCredentialsFile(str)
	case map[string]interface{}:
		var m = file.(map[string]interface{})
		if len(m) < 0 {
			return nil
		}
		var bytes = utils.JsonEncode(m).Bytes()
		if len(bytes) <= 0 {
			return nil
		}
		return option.WithCredentialsJSON(bytes)
	}
	return nil
}

func (admin *fireBaseAdmin) Cancel() {
	if admin != nil && admin.cancel != nil {
		admin.cancel()
	}
	return
}

// GetPusher 获取推送对象
func (admin *fireBaseAdmin) GetPusher() (*messaging.Client, error) {
	var app, err = admin.AcquireAgent()
	if err != nil {
		return nil, err
	}
	return app.Messaging(context.Background())
}

func (admin *fireBaseAdmin) resolvePath(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	dir, _ := os.Getwd()
	if dir != "" {
		if strings.HasPrefix(file, "/") {
			file = dir + file
		} else {
			file = dir + "/" + file
		}
	}
	if path, err := filepath.Abs(file); err == nil {
		return path
	}
	return file
}
