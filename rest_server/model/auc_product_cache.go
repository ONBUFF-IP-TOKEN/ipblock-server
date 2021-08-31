package model

import (
	"fmt"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

const product_list_key = "PRODUCT-LIST"
const product_item_key = "PRODUCT-ITEM"

type ProductListCache struct {
	PageInfo    *context_auc.PageInfoResponse `json:"page_info"`
	ProductList *[]context_auc.ProductInfo    `json:"list"`
}

// product set
func (o *DB) CacheSetProduct(product *context_auc.ProductInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	ckey := genCacheKeyByAucProduct(product_item_key)
	// value := basedb.Z{
	// 	Score:  float64(product.Id),
	// 	Member: product,
	// }
	// return o.Cache.ZAdd(ckey, value)
	return o.Cache.HSet(ckey, strconv.FormatInt(product.Id, 10), product)
}

func (o *DB) CacheDelProduct(prductId int64) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	ckey := genCacheKeyByAucProduct(product_item_key)
	// _, err := o.Cache.ZRemRangeByScore(ckey, strconv.FormatInt(prductId, 10), strconv.FormatInt(prductId, 10))
	// return err
	return o.Cache.HDel(ckey, strconv.FormatInt(prductId, 10))
}

func (o *DB) DeleteProductList() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucProduct(product_list_key)
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

	cKey := genCacheKeyByAucProduct(product_list_key)
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetProductListCache ", field)
	return o.Cache.HSet(cKey, field, productListCache)
}

func (o *DB) GetProductListCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.ProductInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	productListCache := &ProductListCache{}
	cKey := genCacheKeyByAucProduct(product_list_key)
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, productListCache)

	return productListCache.PageInfo, productListCache.ProductList, err
}

func genCacheKeyByAucProduct(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}
