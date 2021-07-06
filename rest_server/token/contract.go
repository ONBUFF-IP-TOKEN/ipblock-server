package token

// func (o *IToken) LoadContract(tokenAddr string) error {
// 	if err := o.ethClient.LoadContract(tokenAddr); err != nil {
// 		log.Error("LoadContract error :", err)
// 		return err
// 	}
// 	o.LoadContractInfo()
// 	return nil
// }

// // 기본 정보 가져오기
// func (o *IToken) LoadContractInfo() error {
// 	var err error
// 	o.tokenName, err = o.ethClient.GetName()
// 	if err != nil {
// 		return err
// 	}
// 	o.tokenSymbol, err = o.ethClient.GetSymbol()
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

// // 선택한 지갑주소에 보유한 코인 개수
// func (o *IToken) GetBalanceOf(address string) (int64, error) {
// 	balance, err := o.ethClient.GetBalanceOf(address)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return balance.Int64(), nil
// }

// // 선택한 토큰 id가 존재하는지 체크
// func (o *IToken) IsExistToken(tokenId int64) (bool, error) {
// 	return o.ethClient.ExistOf(big.NewInt(tokenId))
// }

// // 선택한 토큰 id의 uri 정보 추출
// func (o *IToken) GetUriInfo(tokenId int64) (string, error) {
// 	return o.ethClient.GetTokenUri(big.NewInt(tokenId))
// }

// // 선택한 토큰 id의 owner 정보 추출
// func (o *IToken) GetOwnerOf(tokenId int64) (string, error) {
// 	return o.ethClient.OwnerOf(big.NewInt(tokenId))
// }

// // 토큰 생성
// func (o *IToken) CreateERC721Token(fromAddr, toAddr, uri, privateKey string) (string, error) {
// 	return o.ethClient.CreateERC721func(fromAddr, toAddr, uri, privateKey)
// }

// // 토큰 전송
// func (o *IToken) TransferERC721Token(adminAddr, fromAddr, toAddr, privateKey string, tokenId int64) (string, error) {
// 	return o.ethClient.Transfer(privateKey, adminAddr, fromAddr, toAddr, big.NewInt(tokenId))
// }

// // 토큰 삭제
// func (o *IToken) BurnToken(fromAddr, privateKey string, tokenId int64) (string, error) {
// 	return o.ethClient.Burn(fromAddr, privateKey, big.NewInt(tokenId))
// }

// // 토큰 승인
// func (o *IToken) Approve(fromAddr, privateKey, toAddr string, tokenId int64) (string, error) {
// 	return o.ethClient.Approve(fromAddr, privateKey, toAddr, big.NewInt(tokenId))
// }
