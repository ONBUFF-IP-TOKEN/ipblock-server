package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

// 경매 전체 업데이트
func (o *DB) UpdateAucAuction(auction *context_auc.AucAuctionUpdate) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_auctions set bid_start_amount=?, bid_cur_amount=?, bid_unit=?, bid_deposit=?," +
		"auc_start_ts=?, auc_end_ts=?, auc_state=?, auc_round=?, " +
		"create_ts=?, active_state=?, product_id=?, recommand=?, token_type=?, price=? WHERE auc_id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auction.BidStartAmount, auction.BidCurAmount, auction.BidUnit, auction.BidDeposit,
		auction.AucStartTs, auction.AucEndTs, auction.AucState, auction.AucRound,
		auction.CreateTs, auction.ActiveState, auction.ProductId, auction.Recommand, auction.TokenType, auction.Price, auction.Id)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("UpdateAucAuction id:", id)

	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	o.DeleteAuctionCache(auction.Id)
	return id, nil
}

// 최고가 업데이트
func (o *DB) UpdateAucAuctionBestBid(auctionId int64, curAmount float64) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_auctions set bid_cur_amount=? WHERE auc_id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, curAmount, auctionId)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("UpdateAucAuctionBestBid id:", id)

	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	o.DeleteAuctionCache(auctionId)
	return id, nil
}

// 경매 종료 정보 업데이트
func (o *DB) UpdateAucAuctionAucState(auctionId int64, aucState context_auc.Auction_auc_state, refreshCache bool) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_auctions set auc_state=? WHERE auc_id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, aucState, auctionId)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}

	log.Debug("UpdateAucAuctionAucState id:", id)

	if refreshCache {
		// auction list cache 전체 삭제
		o.DeleteAuctionList()
		o.DeleteAuctionCache(auctionId)
	}

	return id, nil
}
