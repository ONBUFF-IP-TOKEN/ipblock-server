package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

type ProductListCache struct {
	PageInfo    *context_auc.PageInfoResponse `json:"page_info"`
	ProductList *[]context_auc.ProductInfo    `json:"list"`
}

func (o *DB) DeleteProductList() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucProduct("PRODUCT-LIST")
	return o.Cache.Del(cKey)
}

func (o *DB) SetProductListCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.ProductInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	productListCache := &ProductListCache{
		PageInfo:    pageInfo,
		ProductList: data,
	}

	cKey := genCacheKeyByAucProduct("PRODUCT-LIST")
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetProductListCache ", field)
	return o.Cache.HSet(cKey, field, productListCache)
}

func (o *DB) GetProductListCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.ProductInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	productListCache := &ProductListCache{}
	cKey := genCacheKeyByAucProduct("PRODUCT-LIST")
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, productListCache)

	return productListCache.PageInfo, productListCache.ProductList, err
}

func genCacheKeyByAucProduct(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}
