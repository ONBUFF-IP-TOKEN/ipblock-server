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
		sqlQuery = fmt.Sprintf("SELECT * FROM auc_bids LEFT JOIN auc_products on auc_bids.product_id = auc_products.product_id "+
			"WHERE bid_attendee_wallet_address='%v' ORDER BY bid_ts DESC LIMIT %v,%v", walletaddr, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
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
	var nftId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	meBids := make([]context_auc.MeBid, 0)
	for rows.Next() {
		bid := context_auc.MeBid{}
		if err := rows.Scan(&bid.Bid.Id, &bid.Bid.AucId, &bid.Bid.ProductId, &bid.Bid.BidState, &bid.Bid.BidTs, &bid.Bid.BidAttendeeWalletAddr, &bid.Bid.BidAmount, &bidWinnerTxHash, &bid.Bid.BidWinnerState,
			&bid.Bid.DepositAmount, &bid.Bid.DepositTxHash, &bid.Bid.DepositState, &bid.Bid.TokenType,

			&bid.ProductInfo.Id, &title, &bid.ProductInfo.CreateTs, &desc,
			&bid.ProductInfo.OwnerNickName, &bid.ProductInfo.OwnerWalletAddr, &bid.ProductInfo.CreatorNickName, &bid.ProductInfo.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &bid.ProductInfo.NftState,
			&prices, &content, &bid.ProductInfo.IpOwnerShip, &media); err != nil {
			log.Error(err)
		} else {
			aTitle := context_auc.Localization{}
			json.Unmarshal([]byte(title.String), &aTitle)
			bid.ProductInfo.Title = aTitle

			aDesc := context_auc.Localization{}
			json.Unmarshal([]byte(desc.String), &aDesc)
			bid.ProductInfo.Desc = aDesc

			bid.ProductInfo.NftContract = nftContract.String
			bid.ProductInfo.NftId = nftId.Int64
			bid.ProductInfo.NftCreateTxHash = nftCreateHash.String
			bid.ProductInfo.NftUri = nftUri.String

			//prices 변환
			aPrices := []context_auc.ProductPrice{}
			json.Unmarshal([]byte(prices.String), &aPrices)
			bid.ProductInfo.Prices = aPrices

			//content 변환
			aContent := context_auc.Content{}
			json.Unmarshal([]byte(content.String), &aContent)
			bid.ProductInfo.Content = aContent

			//media 변환
			aMedia := context_auc.MediaInfo{}
			json.Unmarshal([]byte(media.String), &aMedia)
			bid.ProductInfo.Media = aMedia

			bid.Bid.BidWinnerTxHash = bidWinnerTxHash.String
			meBids = append(meBids, bid)
		}
	}

	totalCount, err := o.GetTotalAucMeBidSize(walletaddr, winner)

	return meBids, totalCount, err
}

func (o *DB) GetTotalAucMeBidSize(walletAddr string, winner bool) (int64, error) {
	var sqlQuery string
	var dataCount int64
	if !winner {
		sqlQuery = fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_attendee_wallet_Address = '%v'", walletAddr)
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
