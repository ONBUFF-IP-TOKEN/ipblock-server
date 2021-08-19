package model

import (
	"fmt"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

type ProductListCache struct {
	PageInfo    *context_auc.PageInfoResponse `json:"page_info"`
	ProductList *[]context_auc.ProductInfo    `json:"list"`
}

// product set
func (o *DB) CacheSetProduct(product *context_auc.ProductInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	ckey := genCacheKeyProduct()
	return o.Cache.HSet(ckey, strconv.FormatInt(product.Id, 10), product)
}

func (o *DB) CacheDelProduct(prductId int64) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	ckey := genCacheKeyProduct()
	return o.Cache.HDel(ckey, strconv.FormatInt(prductId, 10))
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

func genCacheKeyProduct() string {
	return config.GetInstance().DBPrefix + ":AUCTION:PRODUCT"
}

func genCacheKeyByAucProduct(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}
