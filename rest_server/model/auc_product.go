package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucProduct(product *context_auc.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_products (title, create_ts, description, media_original, media_original_type, media_thumnail, media_thumnail_type, " +
		"links, videos, owner_nickname, owner_wallet_address, creator_nickname, creator_wallet_address," +
		"prices, content ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

	title, _ := json.Marshal(product.Title)
	desc, _ := json.Marshal(product.Desc)
	links, _ := json.Marshal(product.Links)
	videos, _ := json.Marshal(product.Videos)
	prices, _ := json.Marshal(product.Prices)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, string(title), product.CreateTs, string(desc),
		product.MediaOriginal, product.MediaOriginalType, product.MediaThumnail, product.MediaThumnailType,
		string(links), string(videos),
		product.OwnerNickName, product.OwnerWalletAddr, product.CreatorNickName, product.CreatorWalletAddr,
		string(prices), product.Content)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("insert id:", insertId)
	return insertId, nil
}

func (o *DB) UpdateAucProduct(product *context_auc.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_products set title=?, description=?, " +
		"media_original=?, media_original_type=?, media_thumnail=?, media_thumnail_type=?, " +
		"links=?, videos=?, " +
		"owner_nickname=?, owner_wallet_address=?, creator_nickname=?, creator_wallet_address=?, " +
		"prices=?, content=? WHERE product_id=?")

	title, _ := json.Marshal(product.Title)
	desc, _ := json.Marshal(product.Desc)
	links, _ := json.Marshal(product.Links)
	videos, _ := json.Marshal(product.Videos)
	prices, _ := json.Marshal(product.Prices)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, string(title), string(desc),
		product.MediaOriginal, product.MediaOriginalType, product.MediaThumnail, product.MediaThumnailType,
		string(links), string(videos),
		product.OwnerNickName, product.OwnerWalletAddr, product.CreatorNickName, product.CreatorWalletAddr,
		string(prices), product.Content, product.Id)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	Id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("update id:", Id)
	return Id, nil
}

func (o *DB) UpdateAucProductNft(product *context_auc.ProductInfo, nftContract, creatHash, uri string) (int64, error) {
	sqlQuery := "UPDATE auc_products set nft_contract=?, nft_create_txhash=?, nft_uri=? WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, nftContract, creatHash, uri, product.Id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return cnt, nil
}

func (o *DB) DeleteAucProduct(productId int64) (bool, error) {
	sqlQuery := "DELETE FROM auc_products WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, productId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (o *DB) UpdateAucProductNftTokenId(createHash string, tokenId int64) (int64, error) {
	sqlQuery := "UPDATE auc_products set nft_id=?, nft_state=? WHERE nft_create_txhash=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, tokenId, 1, createHash)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return cnt, nil
}

func (o *DB) GetAucProductList(pageInfo *context_auc.ProductList) ([]context_auc.ProductInfo, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_products ORDER BY product_id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var title, desc, links, videos, prices, content sql.NullString
	var nftId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString
	products := make([]context_auc.ProductInfo, 0)
	for rows.Next() {
		product := context_auc.ProductInfo{}
		if err := rows.Scan(&product.Id, &title, &product.CreateTs, &desc,
			&product.MediaOriginal, &product.MediaOriginalType, &product.MediaThumnail, &product.MediaThumnailType,
			&links, &videos,
			&product.OwnerNickName, &product.OwnerWalletAddr, &product.CreatorNickName, &product.CreatorWalletAddr,
			&nftContract, &nftId, &nftCreateHash, &nftUri, &product.NftState,
			&prices, &content); err != nil {
			log.Error(err)
		}

		aTitle := context_auc.Localization{}
		json.Unmarshal([]byte(title.String), &aTitle)
		product.Title = aTitle

		aDesc := context_auc.Localization{}
		json.Unmarshal([]byte(desc.String), &aDesc)
		product.Desc = aDesc

		aLinks := []context_auc.Urls{}
		json.Unmarshal([]byte(links.String), &aLinks)
		product.Links = aLinks

		aVideos := []context_auc.Urls{}
		json.Unmarshal([]byte(videos.String), &aVideos)
		product.Videos = aVideos

		product.NftContract = nftContract.String
		product.NftId = nftId.Int64
		product.NftCreateTxHash = nftContract.String
		product.NftUri = nftUri.String

		//prices 변환
		aPrices := []context_auc.ProductPrice{}
		json.Unmarshal([]byte(prices.String), &aPrices)
		product.Prices = aPrices

		product.Content = content.String

		products = append(products, product)
	}

	totalCount, err := o.GetTotalAucProductSize()

	return products, totalCount, err
}

func (o *DB) GetTotalAucProductSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM auc_products")
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