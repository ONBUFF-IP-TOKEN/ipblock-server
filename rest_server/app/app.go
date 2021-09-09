package app

import (
	"fmt"
	"sync"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/cdn/azure"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/externalapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/internalapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/schedule"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

type ServerApp struct {
	base.BaseApp
	conf       *config.ServerConfig
	configFile string

	token    *token.IToken
	auth     *auth.IAuth
	azure    *azure.Azure
	schedule *schedule.AuctionScheduler
}

func (o *ServerApp) Init(configFile string) (err error) {
	o.conf = config.GetInstance(configFile)
	base.AppendReturnCodeText(&resultcode.ResultCodeText)
	context.AppendRequestParameter()

	if err := o.NewDB(o.conf); err != nil {
		return err
	}
	if err := o.InitScheduler(); err != nil {
		return err
	}
	if err := o.InitToken(); err != nil {
		return err
	}
	if auth, err := auth.NewIAuth(&o.conf.Auth); err != nil {
		return err
	} else {
		o.auth = auth
	}
	// cdn init
	if err := o.InitCdn(); err != nil {
		return err
	}

	return err
}

func (o *ServerApp) CleanUp() {
	fmt.Println("CleanUp")
}

func (o *ServerApp) Run(wg *sync.WaitGroup) error {
	return nil
}

func (o *ServerApp) GetConfig() *baseconf.Config {
	return &o.conf.Config
}

func NewApp() (*ServerApp, error) {
	app := &ServerApp{}

	intAPI := internalapi.NewAPI()
	extAPI := externalapi.NewAPI()

	if err := app.BaseApp.Init(app, intAPI, extAPI); err != nil {
		return nil, err
	}

	return app, nil
}

func (o *ServerApp) NewDB(conf *config.ServerConfig) error {
	auth := conf.MysqlDBAuth
	mysqlDB, err := basedb.GetMysql(auth.Host, auth.ID, auth.Password, auth.Database, auth.PoolSize, auth.IdleSize)
	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v, %v, %v",
			err, auth.Host, auth.ID, auth.Password, auth.Database, auth.PoolSize, auth.IdleSize)
		return err
	}

	gCache := basedb.GetCache(&conf.Cache)

	model.SetDB(mysqlDB, gCache)

	return nil
}

func (o *ServerApp) InitToken() error {
	o.token = token.NewTokenManager(&o.conf.Token, &o.conf.Cdn)

	if err := o.token.Init(); err != nil {
		return err
	}
	return nil
}

func (o *ServerApp) InitCdn() error {
	if azure, err := azure.InitAzure(
		o.conf.Cdn.Azure.AzureStorageAccount,
		o.conf.Cdn.Azure.AzureStorageAccessKey,
		o.conf.Cdn.Azure.Domain,
		o.conf.Cdn.Azure.ContainerNft,
		o.conf.Cdn.Azure.ContainerProduct,
	); err != nil {
		return err
	} else {
		o.azure = azure
	}

	return nil
}

func (o *ServerApp) InitScheduler() error {
	o.schedule = schedule.GetScheduler()
	if o.schedule != nil {
		o.schedule.ResetAuctionScheduler()
	}
	return nil
}
