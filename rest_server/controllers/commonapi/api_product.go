package commonapi

import (
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

func PostRegisterProduct(product *context.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. product table에 저장
	product.CreateTs = datetime.GetTS2MilliSec()
	product.SetStateRegistering()
	if product.QuantityRemaining == 0 {
		product.QuantityRemaining = product.QuantityTotal
	}
	if id, err := model.GetDB().InsertProduct(product); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		product.Id = id

		//2. product_nft table에 전체 개수만큼 저장
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_CreateNft,
			Data:        product,
			Callback:    make(chan interface{}),
		}
		*resp = GetTokenProc(data)
		resp.Value = product
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostUpdateProductState(product *context.ProductUpdateState, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if _, err := model.GetDB().UpdateProductState(product); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetTokenProc(data *basenet.CommandData) base.BaseResponse {
	if ch, exist := context.GetChanInstance().Get(context.TokenChannel); exist {
		ch.(chan *basenet.CommandData) <- data
	}

	ticker := time.NewTicker(90 * time.Second)

	resp := base.BaseResponse{}
	select {
	case callback := <-data.Callback:
		ticker.Stop()
		msg, ok := callback.(*base.BaseResponse)
		if ok {
			resp = *msg
		}
	case <-ticker.C:
		ticker.Stop()
		resp = base.BaseResponseInternalServerError()
	}

	return resp
}
