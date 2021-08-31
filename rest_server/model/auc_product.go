package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucProduct(product *context_auc.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_products (sno, title, create_ts, description, " +
		"owner_nickname, owner_wallet_address, creator_nickname, creator_wallet_address," +
		"prices, content, " +
		"card_bg_color, card_border_color, card_grade, card_tier, " +
		"ip_ownership, ip_ownership_log_url, ip_category, " +
		"media ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

	title, _ := json.Marshal(product.Title)
	desc, _ := json.Marshal(product.Desc)
	prices, _ := json.Marshal(product.Prices)
	content, _ := json.Marshal(product.Content)
	media, _ := json.Marshal(product.Media)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.SNo, string(title), product.CreateTs, string(desc),
		product.OwnerNickName, product.OwnerWalletAddr, product.CreatorNickName, product.CreatorWalletAddr,
		string(prices), content,
		product.CardInfo.BackgroundColor, product.CardInfo.BorderColor, product.CardInfo.CardGrade, product.CardInfo.Tier,
		product.Company.IpOwnerShip, product.Company.IpOwnerShipLogoUrl, product.Company.IpCategory,
		media)

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
	product.Id = insertId
	o.CacheSetProduct(product)
	return insertId, nil
}

// 물품 정보 업데이트
func (o *DB) UpdateAucProduct(product *context_auc.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_products set sno=?, title=?, description=?, " +
		"owner_nickname=?, owner_wallet_address=?, creator_nickname=?, creator_wallet_address=?, " +
		"prices=?, content=?, " +
		"card_bg_color=?, card_border_color=?, card_grade=?, card_tier=?, " +
		"ip_ownership=?, ip_ownership_log_url=?, ip_category=?, " +
		"media=? WHERE product_id=?")

	title, _ := json.Marshal(product.Title)
	desc, _ := json.Marshal(product.Desc)
	prices, _ := json.Marshal(product.Prices)
	content, _ := json.Marshal(product.Content)
	media, _ := json.Marshal(product.Media)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.SNo, string(title), string(desc),
		product.OwnerNickName, product.OwnerWalletAddr, product.CreatorNickName, product.CreatorWalletAddr,
		string(prices), content,
		product.CardInfo.BackgroundColor, product.CardInfo.BorderColor, product.CardInfo.CardGrade, product.CardInfo.Tier,
		product.Company.IpOwnerShip, product.Company.IpOwnerShipLogoUrl, product.Company.IpCategory,
		media, product.Id)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	Id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("UpdateAucProduct id:", Id)

	if Id != 0 {
		// cache 삭제
		o.CacheDelProduct(product.Id)
		// product list cache 전체 삭제
		o.DeleteProductList()
		// auction list cache 전체 삭제
		o.DeleteAuctionList()
		// auction 단일 cache 전체 삭제
		o.DeleteAuctionCacheAll()
	}

	return Id, nil
}

// nft 정보 업데이트
func (o *DB) UpdateAucProductNft(product *context_auc.ProductInfo) (int64, error) {
	sqlQuery := "UPDATE auc_products set nft_contract=?, nft_create_txhash=?, nft_uri=? WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.NftContract, product.NftCreateTxHash, product.NftUri, product.Id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	// product 단일 cache 삭제
	o.CacheDelProduct(product.Id)
	// product list cache 전체 삭제
	o.DeleteProductList()
	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	// auction 단일 cache 전체 삭제
	o.DeleteAuctionCacheAll()

	return cnt, nil
}

// 물품 삭제
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

	// product 단일 cache 삭제
	o.CacheDelProduct(productId)
	// product list cache 전체 삭제
	o.DeleteProductList()
	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	// auction 단일 cache 전체 삭제
	o.DeleteAuctionCacheAll()

	return true, nil
}

// nft 생성 후 생성 정보 업데이트
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

	if newProduct, err := o.GetAucProductByNftCreatHash(createHash); err == nil {
		// product 단일 cache 삭제
		o.CacheDelProduct(newProduct.Id)
		// product list cache 전체 삭제
		o.DeleteProductList()
		// auction list cache 전체 삭제
		o.DeleteAuctionList()
		// auction 단일 cache 전체 삭제
		o.DeleteAuctionCacheAll()
	}

	return cnt, nil
}

func (o *DB) GetAucProductByNftCreatHash(createHash string) (*context_auc.ProductInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_products WHERE nft_create_txhash='%v'", createHash)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	product := &context_auc.ProductInfo{}
	for rows.Next() {
		var err error
		product, err = o.ScanProduct(rows)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	o.CacheSetProduct(product)
	return product, err
}

// 단일 물품 정보 추출
func (o *DB) GetAucProductById(productId int64) (*context_auc.ProductInfo, error) {
	var err error
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_products WHERE product_id=%v", productId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	product := &context_auc.ProductInfo{}
	for rows.Next() {
		product, err = o.ScanProduct(rows)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	o.CacheSetProduct(product)
	return product, err
}

// 물품 리스트 수집 (페이징)
func (o *DB) GetAucProductList(pageInfo *context_auc.ProductList) ([]context_auc.ProductInfo, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_products ORDER BY product_id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	products := make([]context_auc.ProductInfo, 0)
	for rows.Next() {
		product, err := o.ScanProduct(rows)
		if err != nil {
			log.Error(err)
			continue
		}
		o.CacheSetProduct(product)
		products = append(products, *product)
	}

	totalCount, err := o.GetTotalAucProductSize()

	return products, totalCount, err
}

func (o *DB) ScanProduct(rows *sql.Rows) (*context_auc.ProductInfo, error) {
	var title, desc, prices, content, media sql.NullString
	var nftId sql.NullInt64
	var nftContract, nftCreateHash, nftUri sql.NullString

	product := &context_auc.ProductInfo{}
	if err := rows.Scan(&product.Id, &product.SNo, &title, &product.CreateTs, &desc,
		&product.OwnerNickName, &product.OwnerWalletAddr, &product.CreatorNickName, &product.CreatorWalletAddr,
		&nftContract, &nftId, &nftCreateHash, &nftUri, &product.NftState,
		&prices, &content,
		&product.CardInfo.BackgroundColor, &product.CardInfo.BorderColor, &product.CardInfo.CardGrade, &product.CardInfo.Tier,
		&product.Company.IpOwnerShip, &product.Company.IpOwnerShipLogoUrl, &product.Company.IpCategory,
		&media); err != nil {
		log.Error("ScanProduct error: ", err)
		return nil, err
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

	return product, nil
}

func (o *DB) GetTotalAucProductSize() (int64, error) {
	var dataCount int64
	if err := o.Mysql.QueryRow("SELECT COUNT(*) as count FROM auc_products", &dataCount); err != nil {
		log.Error("GetTotalAucProductSize : ", err)
		return 0, err
	}
	return dataCount, nil
}
