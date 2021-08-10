package token

import "strconv"

type NftUriShose struct {
	Idx           int64  `json:"idx" validate:"required"`
	Name          string `json:"name" validate:"required"`
	SerialNo      string `json:"serial_no" validate:"required"`
	Info          string `json:"info"`
	Certification string `json:"certification" validate:"required"`
}

type NftIpblock struct {
	ProductId     int64  `json:"product_id"`
	QuantityIndex int64  `json:"quantity_index"`
	Name          string `json:"name"`
	Certification string `json:"certification" validate:"required"`
}

type NftUriInfo struct {
	UriType  string      `json:"type" validate:"required"`
	CreateTs int64       `json:"create_ts" validate:"required"`
	Data     interface{} `json:"data" validate:"required"`
}

func GetNftUri(domain string, productId, quantityIndex int64) string {
	return domain + strconv.FormatInt(productId, 10) + "_" + strconv.FormatInt(quantityIndex, 10) + ".json"
}
