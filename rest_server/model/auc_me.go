package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) GetAucBidListMe(pageInfo *context_auc.MeBidList, walletaddr string, winner bool) ([]context_auc.MeBid, int64, error) {
	var sqlQuery string
	if !winner {
		if pageInfo.AucId == -1 {
			sqlQuery = fmt.Sprintf("SELECT * FROM auc_bids LEFT JOIN auc_products on auc_bids.product_id = auc_products.product_id "+
				"WHERE bid_attendee_wallet_address='%v' ORDER BY bid_ts DESC LIMIT %v,%v", walletaddr, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
		} else {
			sqlQuery = fmt.Sprintf("SELECT * FROM auc_bids LEFT JOIN auc_products on auc_bids.product_id = auc_products.product_id "+
				"WHERE bid_attendee_wallet_address='%v' AND auc_id=%v ORDER BY bid_ts DESC LIMIT %v,%v", walletaddr, pageInfo.AucId, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
		}

	} else {
		sqlQuery = fmt.Sprintf("SELECT * FROM auc_bids LEFT JOIN auc_products on auc_bids.product_id = auc_products.product_id "+
			"WHERE bid_attendee_wallet_address='%v' AND bid_state=%v ORDER BY bid_ts DESC LIMIT %v,%v", walletaddr, context_auc.Bid_state_success, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	}

	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var bidWinnerTxHash sql.NullString
	var title, desc, prices, content, media sql.NullString
	var nftId, productId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	meBids := make([]context_auc.MeBid, 0)
	for rows.Next() {
		bid := context_auc.MeBid{}
		if err := rows.Scan(&bid.Bid.Id, &bid.Bid.AucId, &bid.Bid.ProductId, &bid.Bid.BidState, &bid.Bid.BidTs, &bid.Bid.BidAttendeeWalletAddr, &bid.Bid.BidAmount, &bidWinnerTxHash, &bid.Bid.BidWinnerState,
			&bid.Bid.TokenType,

			&productId, &bid.ProductInfo.SNo, &title, &bid.ProductInfo.CreateTs, &desc,
			&bid.ProductInfo.OwnerNickName, &bid.ProductInfo.OwnerWalletAddr, &bid.ProductInfo.CreatorNickName, &bid.ProductInfo.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &bid.ProductInfo.NftState,
			&prices, &content,
			&bid.ProductInfo.CardInfo.BackgroundColor, &bid.ProductInfo.CardInfo.BorderColor, &bid.ProductInfo.CardInfo.CardGrade, &bid.ProductInfo.CardInfo.Tier,
			&bid.ProductInfo.Company.IpOwnerShip, &bid.ProductInfo.Company.IpOwnerShipLogoUrl, &bid.ProductInfo.Company.IpCategory,
			&media); err != nil {
			log.Error(err)
		} else {

			aTitle := context_auc.Localization{}
			json.Unmarshal([]byte(title.String), &aTitle)
			bid.ProductInfo.Title = aTitle

			aDesc := context_auc.Localization{}
			json.Unmarshal([]byte(desc.String), &aDesc)
			bid.ProductInfo.Desc = aDesc

			bid.ProductInfo.Id = productId.Int64
			bid.ProductInfo.NftContract = nftContract.String
			bid.ProductInfo.NftId = nftId.Int64
			bid.ProductInfo.NftCreateTxHash = nftCreateHash.String
			bid.ProductInfo.NftUri = nftUri.String

			//prices ??????
			aPrices := []context_auc.ProductPrice{}
			json.Unmarshal([]byte(prices.String), &aPrices)
			bid.ProductInfo.Prices = aPrices

			//content ??????
			aContent := context_auc.Content{}
			json.Unmarshal([]byte(content.String), &aContent)
			bid.ProductInfo.Content = aContent

			//media ??????
			aMedia := context_auc.MediaInfo{}
			json.Unmarshal([]byte(media.String), &aMedia)
			bid.ProductInfo.Media = aMedia

			bid.Bid.BidWinnerTxHash = bidWinnerTxHash.String
			meBids = append(meBids, bid)
		}
	}

	totalCount, err := o.GetTotalAucMeBidSize(walletaddr, winner, pageInfo.AucId)

	return meBids, totalCount, err
}

func (o *DB) GetTotalAucMeBidSize(walletAddr string, winner bool, aucId int64) (int64, error) {
	var sqlQuery string
	var dataCount int64
	if !winner {
		if aucId == -1 {
			sqlQuery = fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_attendee_wallet_Address = '%v'", walletAddr)
		} else {
			sqlQuery = fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_attendee_wallet_Address = '%v' AND auc_id=%v", walletAddr, aucId)
		}

	} else {
		sqlQuery = fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_attendee_wallet_Address = '%v' AND bid_state=%v",
			walletAddr, context_auc.Bid_state_success)
	}

	err := o.Mysql.QueryRow(sqlQuery, &dataCount)
	if err != nil {
		log.Error(err)
		return dataCount, err
	}

	return dataCount, nil
}

func (o *DB) GetAucBidNftListMe(pageInfo *context_auc.MeBidList, walletaddr string) ([]context_auc.MeBid, int64, error) {
	var sqlQuery = fmt.Sprintf("SELECT * FROM auc_bids LEFT JOIN auc_products on auc_bids.product_id = auc_products.product_id "+
		"WHERE bid_attendee_wallet_address='%v' AND bid_winner_state=%v ORDER BY bid_ts DESC LIMIT %v,%v", walletaddr, context_auc.Bid_winner_state_submit_complete, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)

	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var bidWinnerTxHash sql.NullString
	var title, desc, prices, content, media sql.NullString
	var nftId, productId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	meBids := make([]context_auc.MeBid, 0)
	for rows.Next() {
		bid := context_auc.MeBid{}
		if err := rows.Scan(&bid.Bid.Id, &bid.Bid.AucId, &bid.Bid.ProductId, &bid.Bid.BidState, &bid.Bid.BidTs, &bid.Bid.BidAttendeeWalletAddr, &bid.Bid.BidAmount, &bidWinnerTxHash, &bid.Bid.BidWinnerState,
			&bid.Bid.TokenType,

			&productId, &bid.ProductInfo.SNo, &title, &bid.ProductInfo.CreateTs, &desc,
			&bid.ProductInfo.OwnerNickName, &bid.ProductInfo.OwnerWalletAddr, &bid.ProductInfo.CreatorNickName, &bid.ProductInfo.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &bid.ProductInfo.NftState,
			&prices, &content,
			&bid.ProductInfo.CardInfo.BackgroundColor, &bid.ProductInfo.CardInfo.BorderColor, &bid.ProductInfo.CardInfo.CardGrade, &bid.ProductInfo.CardInfo.Tier,
			&bid.ProductInfo.Company.IpOwnerShip, &bid.ProductInfo.Company.IpOwnerShipLogoUrl, &bid.ProductInfo.Company.IpCategory,
			&media); err != nil {
			log.Error(err)
		} else {

			aTitle := context_auc.Localization{}
			json.Unmarshal([]byte(title.String), &aTitle)
			bid.ProductInfo.Title = aTitle

			aDesc := context_auc.Localization{}
			json.Unmarshal([]byte(desc.String), &aDesc)
			bid.ProductInfo.Desc = aDesc

			bid.ProductInfo.Id = productId.Int64
			bid.ProductInfo.NftContract = nftContract.String
			bid.ProductInfo.NftId = nftId.Int64
			bid.ProductInfo.NftCreateTxHash = nftCreateHash.String
			bid.ProductInfo.NftUri = nftUri.String

			//prices ??????
			aPrices := []context_auc.ProductPrice{}
			json.Unmarshal([]byte(prices.String), &aPrices)
			bid.ProductInfo.Prices = aPrices

			//content ??????
			aContent := context_auc.Content{}
			json.Unmarshal([]byte(content.String), &aContent)
			bid.ProductInfo.Content = aContent

			//media ??????
			aMedia := context_auc.MediaInfo{}
			json.Unmarshal([]byte(media.String), &aMedia)
			bid.ProductInfo.Media = aMedia

			bid.Bid.BidWinnerTxHash = bidWinnerTxHash.String
			meBids = append(meBids, bid)
		}
	}

	totalCount, err := o.GetTotalAucMeBidNftSize(walletaddr, pageInfo.AucId)

	return meBids, totalCount, err
}

func (o *DB) GetTotalAucMeBidNftSize(walletAddr string, aucId int64) (int64, error) {
	var sqlQuery string
	var dataCount int64

	sqlQuery = fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_attendee_wallet_Address = '%v' AND bid_winner_state=%v",
		walletAddr, context_auc.Bid_winner_state_submit_complete)

	err := o.Mysql.QueryRow(sqlQuery, &dataCount)
	if err != nil {
		log.Error(err)
		return dataCount, err
	}

	return dataCount, nil
}
