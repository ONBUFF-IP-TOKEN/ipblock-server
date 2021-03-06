package model

import (
	"fmt"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

const auction_list_key = "AUCTION-LIST"
const auction_item_key = "AUCTION-ITEM"

type AuctionListCache struct {
	PageInfo    *context_auc.PageInfoResponse `json:"page_info"`
	AuctionList *[]context_auc.AucAuction     `json:"list"`
}

func genCacheKeyByAucAuction(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}

// 단일 경매 정보 delete
func (o *DB) DeleteAuctionCache(aucId int64) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucAuction(auction_item_key)
	field := strconv.FormatInt(aucId, 10)
	return o.Cache.HDel(cKey, field)
}

// 단일 경매 정보 key 전체 삭제
func (o *DB) DeleteAuctionCacheAll() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucAuction(auction_item_key)
	return o.Cache.Del(cKey)
}

// 단일 경매 정보 set
func (o *DB) SetAuctionCache(data *context_auc.AucAuction) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	strAucId := strconv.FormatInt(data.Id, 10)
	cKey := genCacheKeyByAucAuction(auction_item_key)
	return o.Cache.HSet(cKey, strAucId, data)
}

// 단일 경매 정보 get
func (o *DB) GetAuctionCache(aucId int64) (*context_auc.AucAuction, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	auction := &context_auc.AucAuction{}
	cKey := genCacheKeyByAucAuction(auction_item_key)
	field := strconv.FormatInt(aucId, 10)

	err := o.Cache.HGet(cKey, field, auction)

	return auction, err
}

// 경매 리스트 delete
func (o *DB) DeleteAuctionList() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucAuction(auction_list_key)
	return o.Cache.Del(cKey)
}

// 경매 리스트 set
func (o *DB) SetAuctionListCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.AucAuction) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	auctionListCache := &AuctionListCache{
		PageInfo:    pageInfo,
		AuctionList: data,
	}

	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetAuctionListCache ", field)
	return o.Cache.HSet(cKey, field, auctionListCache)
}

// 경매 리스트 set (경매 상태)
func (o *DB) SetAuctionListByAucStateCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.AucAuction, aucState int64) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	auctionListCache := &AuctionListCache{
		PageInfo:    pageInfo,
		AuctionList: data,
	}

	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v-%v", aucState, reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetAuctionListByAucStateCache ", field)
	return o.Cache.HSet(cKey, field, auctionListCache)
}

// 경매 리스트 set (추천 경매)
func (o *DB) SetAuctionListByRecommandCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.AucAuction) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	auctionListCache := &AuctionListCache{
		PageInfo:    pageInfo,
		AuctionList: data,
	}

	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v-%v", "r", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetAuctionListByRecommandCache ", field)
	return o.Cache.HSet(cKey, field, auctionListCache)
}

// 경매 리스트 get
func (o *DB) GetAuctionListCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.AucAuction, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	auctionListCache := &AuctionListCache{}
	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, auctionListCache)

	return auctionListCache.PageInfo, auctionListCache.AuctionList, err
}

// 경매 리스트 get (경매 상태)
func (o *DB) GetAuctionListByAucStateCache(pageInfo *context_auc.PageInfo, aucState int64) (*context_auc.PageInfoResponse, *[]context_auc.AucAuction, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	auctionListCache := &AuctionListCache{}
	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v-%v", aucState, pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, auctionListCache)

	return auctionListCache.PageInfo, auctionListCache.AuctionList, err
}

// 경매 리스트 get (추천 경매)
func (o *DB) GetAuctionListByRecommandCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.AucAuction, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	auctionListCache := &AuctionListCache{}
	cKey := genCacheKeyByAucAuction(auction_list_key)
	field := fmt.Sprintf("%v-%v-%v", "r", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, auctionListCache)

	return auctionListCache.PageInfo, auctionListCache.AuctionList, err
}
