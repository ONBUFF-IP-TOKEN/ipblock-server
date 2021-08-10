package context_auc

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

// product
type ProductInfo struct {
	Id       int64        `json:"product_id"`
	Title    Localization `json:"title"`
	CreateTs int64        `json:"create_ts"`
	Desc     Localization `json:"desc"`

	MediaOriginal     string `json:"media_original"`
	MediaOriginalType string `json:"media_original_type"`
	MediaThumnail     string `json:"media_thumnail"`
	MediaThumnailType string `json:"media_thumnail_type"`

	Links  []Urls `json:"links"`
	Videos []Urls `json:"videos"`

	OwnerNickName     string `json:"owner_nickname"`
	OwnerWalletAddr   string `json:"owner_wallet_address"`
	CreatorNickName   string `json:"creator_nickname"`
	CreatorWalletAddr string `json:"creator_wallet_address"`

	NftContract     string `json:"nft_contract"`
	NftId           int64  `json:"nft_id"`
	NftCreateTxHash string `json:"nft_create_txhash"`
	NftUri          string `json:"nft_uri"`
	NftState        int64  `json:"nft_state"`

	Prices []ProductPrice `json:"product_prices"`

	Content string `json:"content"`
}

func NewProductInfo() *ProductInfo {
	return new(ProductInfo)
}

func (o *ProductInfo) CheckValidate() *base.BaseResponse {
	if len(o.Title.En) == 0 || len(o.Title.Ko) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_Requiredtitle)
	}
	if len(o.Desc.En) == 0 || len(o.Desc.Ko) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireDescription)
	}
	if len(o.MediaOriginal) == 0 || len(o.MediaOriginalType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireMediaOriginal)
	}
	if len(o.MediaThumnail) == 0 || len(o.MediaThumnailType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireMediaThumnail)
	}
	if len(o.OwnerNickName) == 0 || len(o.OwnerWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireOwnerInfo)
	}
	if len(o.CreatorNickName) == 0 || len(o.CreatorWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireCreatorInfo)
	}
	if len(o.Prices) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequirePriceInfo)
	}
	return nil
}

////////////////////////////////////////////////

// proudct 삭제
type RemoveProduct struct {
	Id int64 `query:"product_id"`
}

func NewRemoveProduct() *RemoveProduct {
	return new(RemoveProduct)
}

func (o *RemoveProduct) CheckValidate() *base.BaseResponse {
	if o.Id == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireProductId)
	}
	return nil
}

////////////////////////////////////////////////

// product list request
type ProductList struct {
	PageInfo
}

func NewProductList() *ProductList {
	return new(ProductList)
}

func (o *ProductList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type ProductListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Products []ProductInfo    `json:"products"`
}

////////////////////////////////////////////////
