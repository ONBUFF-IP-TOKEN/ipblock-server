package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/constant"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
	"github.com/labstack/echo"
)

func PostRegisterItem(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewRegisterItem()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	// token 생성
	txHash, err := token.GetToken().Tokens[token.Token_nft].Nft_CreateERC721(params.WalletAddr, params.Thumbnail)
	if err != nil {
		resp.SetResult(constant.Result_TokenERC721CreateError)
	} else {
		params.CreateHash = txHash
		if itemId, err := model.GetDB().InsertItem(params); err != nil {
			resp.SetResult(constant.Result_DBError)
		} else {
			resp.Success()
			resp.Value = context.RegisterItemResponse{
				ItemId: itemId,
				TxHash: txHash,
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func DeleteUnregisterItem(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewUnregisterItem()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	// token id 추출
	item, err := model.GetDB().GetItem(params.ItemId)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		// token burn 시도
		txHash, err := token.GetToken().Tokens[token.Token_nft].Nft_Burn(item.TokenId)
		if err != nil {
			resp.SetResult(constant.Result_TokenERC721BurnError)
		} else {
			// db 삭제 없이 리턴하고 후에 콜백 성공하면 삭제한다.
			resp.Success()
			resp.Value = context.RegisterItemResponse{
				ItemId: params.ItemId,
				TxHash: txHash,
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetItemList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewGetItemList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	items, totalCount, err := model.GetDB().GetItemList(params)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: params.PageOffset,
			PageSize:   int64(len(items)),
			TotalSize:  totalCount,
		}
		resp.Value = context.GetItemListResponse{
			PageInfo: pageInfo,

			ItemInfos: items,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func PostPurchaseItem(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewPostPurchaseItem()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	//item, err := model.GetDB().GetItem(params.ItemId)
	itemInfo, err := model.GetDB().GetItem(params.ItemId)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	}
	// 구매 tx hash 검사
	token.GetToken().Tokens[token.Token_onit].CheckTransferReceipt(params, itemInfo)
	resp.Success()
	resp.Value = context.PostPurchaseItemResponse{
		ItemId: params.ItemId,
	}

	return c.JSON(http.StatusOK, resp)
}
