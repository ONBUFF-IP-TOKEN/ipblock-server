package context_auc

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type CardInfo struct {
	BackgroundColor string `json:"bg_color"`     // 0xffffff
	BorderColor     string `json:"border_color"` // 0xffcc00
	CardGrade       string `json:"grade"`        // level 4 start
	Tier            string `json:"tier"`         // bronze, silver, gold, platinum ....
}

// content info
type Content struct {
	CardInfo CardInfo `json:"card_info"`

	BackgroundColor string `json:"bg_color"` // 0xffffff
}

// media info
type MediaInfo struct {
	MediaOriginal     string `json:"media_origin"`      // url
	MediaOriginalType string `json:"media_origin_type"` // video/mp4, image/png....

	MainImg         string `json:"main_img"`
	MainImgThumnail string `json:"main_img_thumnail"`
	MainImgType     string `json:"main_img_type"` // image/gif, image/png

	SubImg         []string `json:"sub_img"`
	SubImgThumnail []string `json:"sub_img_thumnail"`
	SubImgType     string   `json:"sub_img_type"` // image/png

	CertifiImg         string `json:"certifi_img"`
	CertifiIMgThumnail string `json:"certifi_img_thumnail"`
	CertifiImgType     string `json:"certifi_img_type"` // image/png

	Links  []string `json:"links"`
	Videos []string `json:"videos"`
}

// product
type ProductInfo struct {
	Id       int64        `json:"product_id"`
	Title    Localization `json:"title"`
	CreateTs int64        `json:"create_ts"`
	Desc     Localization `json:"desc"`

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

	Content     Content `json:"content"`
	IpOwnerShip string  `json:"ip_ownership"`

	Media MediaInfo `json:"media"`
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
	if len(o.OwnerNickName) == 0 || len(o.OwnerWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireOwnerInfo)
	}
	if len(o.CreatorNickName) == 0 || len(o.CreatorWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireCreatorInfo)
	}
	if len(o.Prices) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequirePriceInfo)
	}
	if len(o.IpOwnerShip) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireIPOwnerShip)
	}
	return nil
}

////////////////////////////////////////////////

// product 업데이트
type UpdateProduct struct {
	ProductInfo
}

func NewUpdateProduct() *UpdateProduct {
	return new(UpdateProduct)
}

func (o *UpdateProduct) CheckValidate() *base.BaseResponse {
	if len(o.Title.En) == 0 || len(o.Title.Ko) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_Requiredtitle)
	}
	if len(o.Desc.En) == 0 || len(o.Desc.Ko) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Product_RequireDescription)
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
