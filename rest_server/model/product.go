package model

import (
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertProduct(product *context.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO product(product_title, product_thumbnail_url, product_price, product_type, token_type," +
		"create_ts, creator, description, content, quantity_total, quantity_remaining, state) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.Title, product.Thumbnail, product.Price, product.ProductType, product.TokenType,
		product.CreateTs, product.Creator, product.Desc, product.Content, product.QuantityTotal, product.QuantityRemaining, product.State)
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

func (o *DB) DeleteProduct(product *context.UnregisterProduct) (bool, error) {
	sqlQuery := "DELETE FROM product WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.ProductId)
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

func (o *DB) UpdateProduct(product *context.ProductInfo) (int64, error) {
	sqlQuery := "UPDATE product set product_title=?, product_thumbnail_url=?, product_price=?, product_type=?, " +
		"token_type=?, creator=?, description=?, content=?, state=? WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.Title, product.Thumbnail, product.Price, product.ProductType,
		product.TokenType, product.Creator, product.Desc, product.Content, product.State, product.Id)
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

func (o *DB) UpdateProductState(product *context.ProductUpdateState) (int64, error) {
	sqlQuery := "UPDATE product set state=? WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.State, product.ProductId)
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

func (o *DB) GetProductList(pageInfo *context.ProductList) ([]context.ProductInfo, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM product ORDER BY product_id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var creator, thumbnail, content sql.NullString
	products := make([]context.ProductInfo, 0)
	for rows.Next() {
		product := context.ProductInfo{}
		if err := rows.Scan(&product.Id, &product.Title, &thumbnail, &product.Price, &product.ProductType, &product.TokenType,
			&product.CreateTs, &creator, &product.Desc, &content, &product.QuantityTotal, &product.QuantityRemaining, &product.State); err != nil {
			log.Error(err)
		}
		product.Thumbnail = thumbnail.String
		product.Creator = creator.String
		product.Content = content.String
		products = append(products, product)
	}

	totalCount, err := o.GetTotalProductSize()

	return products, totalCount, err
}

func (o *DB) UpdateProductRemain(increase bool, productId int64) (int64, error) {
	var sqlQuery string
	if increase {
		sqlQuery = "UPDATE product set quantity_remaining=quantity_remaining+1 WHERE product_id=?"
	} else {
		sqlQuery = "UPDATE product set quantity_remaining=quantity_remaining-1 WHERE product_id=?"
	}

	result, err := o.Mysql.PrepareAndExec(sqlQuery, productId)
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

func (o *DB) GetTotalProductSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM product")
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

func (o *DB) GetProductInfo(productId int64) (*context.ProductInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM product WHERE product_id=%v", productId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var creator, thumbnail, content sql.NullString
	product := &context.ProductInfo{}
	for rows.Next() {
		if err := rows.Scan(&product.Id, &product.Title, &thumbnail, &product.Price, &product.ProductType, &product.TokenType,
			&product.CreateTs, &creator, &product.Desc, &content, &product.QuantityTotal, &product.QuantityRemaining, &product.State); err != nil {
			log.Error(err)
		}
		product.Thumbnail = thumbnail.String
		product.Creator = creator.String
		product.Content = content.String
	}

	return product, nil
}
