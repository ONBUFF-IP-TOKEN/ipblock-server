package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

// IPBlockServerContext API의 Request Context
type IPBlockServerContext struct {
	*base.BaseContext
	walletAddr string
}

// NewIPBlockServerContext 새로운 IPBlockserver Context 생성
func NewIPBlockServerContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(IPBlockServerContext)
	ctx.BaseContext = baseCtx

	return ctx
}

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
}

func (o *IPBlockServerContext) SetWalletAddr(walletAddr string) {
	o.walletAddr = walletAddr
}

func (o *IPBlockServerContext) WalletAddr() string {
	return o.walletAddr
}
