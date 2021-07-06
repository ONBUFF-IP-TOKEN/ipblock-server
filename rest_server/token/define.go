package token

var gToken *IToken
var gNullAddress = "0x0000000000000000000000000000000000000000"

const (
	Token_nft  = 0
	Token_onit = 1

	token_state_pending  = "pending"
	token_state_mint     = "mint"
	token_state_transfer = "transfer from"
	token_state_burn     = "burn"
)

var tokenTypes = map[int]string{
	Token_nft:  "NFT",
	Token_onit: "ONIT",
}
