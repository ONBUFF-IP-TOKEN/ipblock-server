package commonapi_auc

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func GetAucBidListMe(pageInfo *context_auc.MeBidList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	meBids, totalCount, err := model.GetDB().GetAucBidListMe(pageInfo, pageInfo.WalletAddr, false)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: pageInfo.PageOffset,
			PageSize:   int64(len(meBids)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.MeBidListResponse{
			PageInfo: pageInfo,
			MeBids:   meBids,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 내 낙찰 내역
func GetAucBidWinnerListMe(pageInfo *context_auc.MeBidList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	meBids, totalCount, err := model.GetDB().GetAucBidListMe(pageInfo, pageInfo.WalletAddr, true)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: pageInfo.PageOffset,
			PageSize:   int64(len(meBids)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.MeBidListResponse{
			PageInfo: pageInfo,
			MeBids:   meBids,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetAucNftListMe(pageInfo *context_auc.MeBidList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	meBids, totalCount, err := model.GetDB().GetAucBidListMe(pageInfo, pageInfo.WalletAddr, true)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: pageInfo.PageOffset,
			PageSize:   int64(len(meBids)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.MeBidListResponse{
			PageInfo: pageInfo,
			MeBids:   meBids,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
