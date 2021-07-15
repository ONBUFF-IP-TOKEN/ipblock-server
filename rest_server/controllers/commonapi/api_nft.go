package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func GetNftList(NftList *context.NftList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	nfts, totalCount, err := model.GetDB().GetNftList(NftList)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: NftList.PageOffset,
			PageSize:   int64(len(nfts)),
			TotalSize:  totalCount,
		}
		resp.Value = context.NftListResponse{
			PageInfo: pageInfo,
			Nfts:     nfts,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
