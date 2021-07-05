package token

var gToken *IToken
var gNullAddress = "0x0000000000000000000000000000000000000000"

const (
	token_state_pending  = "pending"
	token_state_mint     = "mint"
	token_state_transfer = "transfer from"
	token_state_burn     = "burn"
)
