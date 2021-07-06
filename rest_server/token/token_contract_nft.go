package token

import (
	"math/big"
)

func (o *Token) Nft_CreateERC721(wallertAddr, uri string) (string, error) {
	return o.Nft_CreateERC721Token(o.conf.ServerWalletAddr, wallertAddr, uri, o.conf.ServerPrivateKey)
}

func (o *Token) Nft_TransferERC721(fromAddr, toAddr string, tokenId int64) (string, error) {
	return o.Nft_TransferERC721Token(o.conf.ServerWalletAddr, fromAddr, toAddr, o.conf.ServerPrivateKey, tokenId)
}

func (o *Token) Nft_Burn(tokenId int64) (string, error) {
	return o.Nft_BurnToken(o.conf.ServerWalletAddr, o.conf.ServerPrivateKey, tokenId)
}

func (o *Token) Nft_LoadContract(tokenAddr string) error {
	err := o.eth.LoadContract(tokenAddr)
	return err
}

// 기본 정보 가져오기
func (o *Token) Nft_LoadContractInfo() error {
	var err error
	o.tokenName, err = o.eth.GetName()
	if err != nil {
		return err
	}
	o.tokenSymbol, err = o.eth.GetSymbol()
	if err != nil {
		return err
	}
	return err
}

// 선택한 지갑주소에 보유한 코인 개수
func (o *Token) Nft_GetBalanceOf(address string) (int64, error) {
	balance, err := o.eth.GetBalanceOf(address)
	return balance.Int64(), err
}

// 선택한 토큰 id가 존재하는지 체크
func (o *Token) Nft_IsExistToken(tokenId int64) (bool, error) {
	return o.eth.ExistOf(big.NewInt(tokenId))
}

// 선택한 토큰 id의 uri 정보 추출
func (o *Token) Nft_GetUriInfo(tokenId int64) (string, error) {
	return o.eth.GetTokenUri(big.NewInt(tokenId))
}

// 선택한 토큰 id의 owner 정보 추출
func (o *Token) Nft_GetOwnerOf(tokenId int64) (string, error) {
	return o.eth.OwnerOf(big.NewInt(tokenId))
}

// 토큰 생성
func (o *Token) Nft_CreateERC721Token(fromAddr, toAddr, uri, privateKey string) (string, error) {
	return o.eth.CreateERC721func(fromAddr, toAddr, uri, privateKey)
}

// 토큰 전송
func (o *Token) Nft_TransferERC721Token(adminAddr, fromAddr, toAddr, privateKey string, tokenId int64) (string, error) {
	return o.eth.Transfer(privateKey, adminAddr, fromAddr, toAddr, big.NewInt(tokenId))
}

// 토큰 삭제
func (o *Token) Nft_BurnToken(fromAddr, privateKey string, tokenId int64) (string, error) {
	return o.eth.Burn(fromAddr, privateKey, big.NewInt(tokenId))
}

// 토큰 승인
func (o *Token) Nft_Approve(fromAddr, privateKey, toAddr string, tokenId int64) (string, error) {
	return o.eth.Approve(fromAddr, privateKey, toAddr, big.NewInt(tokenId))
}
