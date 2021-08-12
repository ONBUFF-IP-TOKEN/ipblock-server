package commonapi_auc

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

func PostAucProductRegister(product *context_auc.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_product table에 저장
	product.CreateTs = datetime.GetTS2MilliSec()
	if id, err := model.GetDB().InsertAucProduct(product); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		product.Id = id

		//2. nft 토큰 생성
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_CreatNftByAut,
			Data:        product,
			Callback:    make(chan interface{}),
		}
		*resp = commonapi.GetTokenProc(data)
		resp.Value = product
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostAucProductUpdate(product *context_auc.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_product table에 업데이트
	if _, err := model.GetDB().UpdateAucProduct(product); err != nil {
		log.Error("UpdateAucProduct error : ", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Value = product
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteAucProductRemove(product *context_auc.RemoveProduct, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteAucProduct(product.Id); err != nil {
		log.Error("DeleteAucProductRemove :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetAucProductList(productList *context_auc.ProductList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)

	//redis exist check
	if pageInfo, products, err := model.GetDB().GetProductListCache(&productList.PageInfo); err == nil {
		resp.Success()
		resp.Value = context_auc.ProductListResponse{
			PageInfo: *pageInfo,
			Products: *products,
		}
	} else {
		// cache 에 없다면 db에서 직접 로드
		products, totalCount, err := model.GetDB().GetAucProductList(productList)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: productList.PageOffset,
				PageSize:   int64(len(products)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.ProductListResponse{
				PageInfo: pageInfo,
				Products: products,
			}
			model.GetDB().SetProductListCache(&productList.PageInfo, &pageInfo, &products)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
