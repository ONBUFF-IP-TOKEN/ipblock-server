package test

import (
	"os"
	"testing"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
)

type testData struct {
	intURIBase   string
	intMinAPIVer string
	intMaxAPIVer string
	extURIBase   string
	extMinAPIVer string
	extMaxAPIVer string
	conf         *config.ServerConfig
}

var gTestData testData

func setupTest() {
	conf := config.GetInstance("./etc/conf/config.local.yml")
	if conf == nil {
		log.Error("not load config file")
	}
	gTestData = testData{
		intURIBase:   "http://localhost:10360",
		intMinAPIVer: "m1.0",
		intMaxAPIVer: "m1.0",
		extURIBase:   "http://localhost:20284",
		extMinAPIVer: "v1.0",
		extMaxAPIVer: "v1.0",
		conf:         conf}
}

func TestMain(m *testing.M) {
	setupTest()
	log.Info("TestMain start")

	ret := m.Run()

	os.Exit(ret)
}

func TestAuth(m *testing.T) {
	log.Debug("TestAuth")
	// 	iAuth, err := auth.NewIAuth(&gTestData.conf.Auth)
	// 	if err != nil {
	// 		log.Error(err)
	// 		return
	// 	}

	// 	vd, err := iAuth.EncryptLoginVD("0x9Ec7EDE9204E17dfa34e1d381ED5f49A0D578e96")
	// 	log.Debug("make vd : ", vd)

	// 	jwtToken, _, err := iAuth.EncryptJwt("0x9Ec7EDE9204E17dfa34e1d381ED5f49A0D578e96")
	// 	if err != nil {
	// 		log.Error("make error : ", err)
	// 	}
	// 	log.Debug("make token : ", jwtToken)
}
