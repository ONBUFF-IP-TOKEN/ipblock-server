package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucAuction(auction *context_auc.AucAuctionRegister) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_auctions (bid_start_amount, bid_cur_amount, bid_unit, bid_deposit, " +
		"auc_start_ts, auc_end_ts, auc_state, auc_round," +
		"create_ts, active_state, product_id, recommand" +
		") VALUES (?,?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auction.BidStartAmount, auction.BidCurAmount, auction.BidUnit, auction.BidDeposit,
		auction.AucStartTs, auction.AucEndTs, auction.AucState, auction.AucRound,
		auction.CreateTs, auction.ActiveState, auction.ProductId, auction.Recommand)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("InsertAucAuction id:", insertId)
	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	return insertId, nil
}

func (o *DB) UpdateAucAuction(auction *context_auc.AucAuctionUpdate) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_auctions set bid_start_amount=?, bid_cur_amount=?, bid_unit=?, bid_deposit=?," +
		"auc_start_ts=?, auc_end_ts=?, auc_state=?, auc_round=?, " +
		"create_ts=?, active_state=?, product_id=?, recommand=? WHERE auc_id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auction.BidStartAmount, auction.BidCurAmount, auction.BidUnit, auction.BidDeposit,
		auction.AucStartTs, auction.AucEndTs, auction.AucState, auction.AucRound,
		auction.CreateTs, auction.ActiveState, auction.ProductId, auction.Recommand, auction.Id)

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
	return id, nil
}

func (o *DB) GetAucAuctionList(pageInfo *context_auc.AuctionList) ([]context_auc.AucAuction, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_auctions LEFT JOIN auc_products on auc_auctions.product_id = auc_products.product_id "+
		"ORDER BY auc_id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var title, desc, prices, content, media sql.NullString
	var nftId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	auctions := make([]context_auc.AucAuction, 0)
	for rows.Next() {
		auction := context_auc.AucAuction{}
		product := context_auc.ProductInfo{}
		if err := rows.Scan(&auction.Id, &auction.BidStartAmount, &auction.BidCurAmount, &auction.BidUnit, &auction.BidDeposit,
			&auction.AucStartTs, &auction.AucEndTs, &auction.AucState, &auction.AucRound,
			&auction.CreateTs, &auction.ActiveState, &auction.ProductId, &auction.Recommand,

			&product.Id, &title, &product.CreateTs, &desc,
			&product.OwnerNickName, &product.OwnerWalletAddr, &product.CreatorNickName, &product.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &product.NftState,
			&prices, &content, &product.IpOwnerShip, &media); err != nil {
			log.Error(err)
		}

		aTitle := context_auc.Localization{}
		json.Unmarshal([]byte(title.String), &aTitle)
		product.Title = aTitle

		aDesc := context_auc.Localization{}
		json.Unmarshal([]byte(desc.String), &aDesc)
		product.Desc = aDesc

		product.NftContract = nftContract.String
		product.NftId = nftId.Int64
		product.NftCreateTxHash = nftCreateHash.String
		product.NftUri = nftUri.String

		//prices 변환
		aPrices := []context_auc.ProductPrice{}
		json.Unmarshal([]byte(prices.String), &aPrices)
		product.Prices = aPrices

		//content 변환
		aContent := context_auc.Content{}
		json.Unmarshal([]byte(content.String), &aContent)
		product.Content = aContent

		//media 변환
		aMedia := context_auc.MediaInfo{}
		json.Unmarshal([]byte(media.String), &aMedia)
		product.Media = aMedia

		auction.ProductInfo = product
		auctions = append(auctions, auction)
	}

	totalCount, err := o.GetTotalAucAuctionSize()

	return auctions, totalCount, err
}

func (o *DB) GetAucAuction(aucId int64) (*context_auc.AucAuction, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_auctions LEFT JOIN auc_products on auc_auctions.product_id = auc_products.product_id WHERE auc_id=%v", aucId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var title, desc, prices, content, media sql.NullString
	var nftId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	auction := &context_auc.AucAuction{}
	product := context_auc.ProductInfo{}
	for rows.Next() {
		if err := rows.Scan(&auction.Id, &auction.BidStartAmount, &auction.BidCurAmount, &auction.BidUnit, &auction.BidDeposit,
			&auction.AucStartTs, &auction.AucEndTs, &auction.AucState, &auction.AucRound,
			&auction.CreateTs, &auction.ActiveState, &auction.ProductId, &auction.Recommand,

			&product.Id, &title, &product.CreateTs, &desc,
			&product.OwnerNickName, &product.OwnerWalletAddr, &product.CreatorNickName, &product.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &product.NftState,
			&prices, &content, &product.IpOwnerShip, &media); err != nil {
			log.Error(err)
		}

		aTitle := context_auc.Localization{}
		json.Unmarshal([]byte(title.String), &aTitle)
		product.Title = aTitle

		aDesc := context_auc.Localization{}
		json.Unmarshal([]byte(desc.String), &aDesc)
		product.Desc = aDesc

		product.NftContract = nftContract.String
		product.NftId = nftId.Int64
		product.NftCreateTxHash = nftCreateHash.String
		product.NftUri = nftUri.String

		//prices 변환
		aPrices := []context_auc.ProductPrice{}
		json.Unmarshal([]byte(prices.String), &aPrices)
		product.Prices = aPrices

		//content 변환
		aContent := context_auc.Content{}
		json.Unmarshal([]byte(content.String), &aContent)
		product.Content = aContent

		//media 변환
		aMedia := context_auc.MediaInfo{}
		json.Unmarshal([]byte(media.String), &aMedia)
		product.Media = aMedia

		auction.ProductInfo = product
	}

	return auction, err
}

func (o *DB) DeleteAucAuction(auctionId int64) (bool, error) {
	sqlQuery := "DELETE FROM auc_auctions WHERE auc_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auctionId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	// auction list cache 전체 삭제
	o.DeleteAuctionList()

	return true, nil
}

func (o *DB) GetTotalAucAuctionSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM auc_auctions")
	var count int64
	if err != nil {
		log.Error(err)
		return count, err
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}

	return count, nil
}
