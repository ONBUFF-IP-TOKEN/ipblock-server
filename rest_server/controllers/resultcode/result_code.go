package resultcode

const (
	Result_Success                = 0
	Result_RequireWalletAddress   = 12000
	Result_RequireTitle           = 12001
	Result_RequireTokenType       = 12002
	Result_RequireThumbnailUrl    = 12003
	Result_RequireValidTokenPrice = 12004
	Result_RequireValidExpireDate = 12005
	Result_RequireCreator         = 12006
	Result_RequireValidItemId     = 12007
	Result_RequireValidPageOffset = 12008
	Result_RequireValidPageSize   = 12009
	Result_RequireDescription     = 12010
	Result_InvalidWalletAddress   = 12011
	Result_RequiredPurchaseTxHash = 12012
	Result_RequireEmailInfo       = 12013

	Result_Product_RequiredTitle        = 12500
	Result_Product_RequiredThumbnailUrl = 12501
	Result_Product_RequiredProductType  = 12502
	Result_Product_RequiredTokenType    = 12503
	Result_Product_RequiredCreator      = 12504
	Result_Product_RequiredDesc         = 12505
	Result_Product_RequireQuantityTotal = 12506
	Result_Product_RequireVaildId       = 12507
	Result_Product_RequireValidState    = 12508
	Result_Product_NotOnSale            = 12509
	Result_Product_LackOfQuantity       = 12510

	Result_DBError           = 13000
	Result_DBNotExistItem    = 13001
	Result_DBNotExistProduct = 13002
	Result_DBNotExistAuction = 13003

	Result_TokenError               = 14000
	Result_TokenERC721CreateError   = 14001
	Result_TokenERC721BurnError     = 14002
	Result_TokenERC721TransferError = 14003
	Result_Reused_Txhash            = 14004

	Result_Auc_Product_Requiredtitle        = 15000
	Result_Auc_Product_RequireDescription   = 15001
	Result_Auc_Product_RequireMediaOriginal = 15002
	Result_Auc_Product_RequireMediaThumnail = 15003
	Result_Auc_Product_RequireOwnerInfo     = 15004
	Result_Auc_Product_RequireCreatorInfo   = 15005
	Result_Auc_Product_RequirePriceInfo     = 15006
	Result_Auc_Product_RequireProductId     = 15007
	Result_Auc_Product_RequireIPOwnerShip   = 15008

	Result_Auc_Auction_RequireBidStartAmount = 15101
	Result_Auc_Auction_RequireBidUnit        = 15102
	Result_Auc_Auction_RequireStartTs        = 15103
	Result_Auc_Auction_RequireEndTs          = 15104
	Result_Auc_Auction_RequireRound          = 15105
	Result_Auc_Auction_RequireActiveState    = 15106
	Result_Auc_Auction_RequireProductId      = 15107
	Result_Auc_Auction_RequireAucId          = 15108
	Result_Auc_Auction_NotPeriod             = 15109 // 경매기간이 아니다.
	Result_Auc_Auction_NotOverYet            = 15110 // 아직 경매가 끝나지 않았다.

	Result_Auc_Bid_RequireAucId         = 15201
	Result_Auc_Bid_RequireWalletAddress = 15202
	Result_Auc_Bid_RequireAmount        = 15203
	Result_Auc_Bid_InvalidWalletAddress = 15204
	Result_Auc_Bid_RequireDepositTxHash = 15205
	Result_Auc_Bid_AlreadyBestAttendee  = 15206 // 이미 최고 입찰자일때 에러
	Result_Auc_Bid_NotBestBidAmount     = 15207 // 입찰 가격이 기존 가격보다 낮을때 에러
	Result_Auc_Bid_RequireDeposit       = 15208 // 입찰전에 보증금 지불을 하지 먼저 하지 않았을때 에러
	Result_Auc_Bid_AlreadyDeposit       = 15209 // 이미 입찰 보증금을 지불하였다.
	Result_Auc_Bid_NotWinner            = 15210 // 낙찰자가 아니다.
	Result_Auc_Bid_RequireDepoistAgree  = 15211 // 보증금 납부 약관 동의가 필요하다.
	Result_Auc_Bid_RequirePrivacyAgree  = 15212 // 개인정보 이용 동의 약관이 필요하다.

	Result_Auc_Notice_RequireNotice = 15301 // 공지사항이 한개라도 필요 하다.
	Result_Auc_Notice_RequireId     = 15302 // 공지사항의 id값이 필요하다.

	Result_Auth_RequireMessage    = 20000
	Result_Auth_RequireSign       = 20001
	Result_Auth_InvalidLoginInfo  = 20002
	Result_Auth_DontEncryptJwt    = 20003
	Result_Auth_InvalidJwt        = 20004
	Result_Auth_InvalidWalletType = 20005
)

var ResultCodeText = map[int]string{
	Result_Success:                "success",
	Result_RequireWalletAddress:   "Wallet address is required",
	Result_RequireTitle:           "Item title is required",
	Result_RequireTokenType:       "Token type is required",
	Result_RequireThumbnailUrl:    "Thumbnail url is required",
	Result_RequireValidTokenPrice: "Valid token price is required",
	Result_RequireValidExpireDate: "Valid expire date is required",
	Result_RequireCreator:         "Creator is required",
	Result_RequireValidItemId:     "Valid item id is required",
	Result_RequireValidPageOffset: "Valid page offset is required",
	Result_RequireValidPageSize:   "Valid page size is required",
	Result_RequireDescription:     "Description is required",
	Result_InvalidWalletAddress:   "Invalid Wallet Address",
	Result_RequiredPurchaseTxHash: "Require purchase tx hash info",
	Result_RequireEmailInfo:       "Require email info",

	Result_Product_RequiredTitle:        "Require product title",
	Result_Product_RequiredThumbnailUrl: "Require product thumbnail url",
	Result_Product_RequiredProductType:  "Require product type info",
	Result_Product_RequiredTokenType:    "Require token type",
	Result_Product_RequiredCreator:      "Require creator info",
	Result_Product_RequiredDesc:         "Require description info",
	Result_Product_RequireQuantityTotal: "Require total quantity",
	Result_Product_RequireVaildId:       "Require valid product id",
	Result_Product_RequireValidState:    "Require valid product state",
	Result_Product_NotOnSale:            "Not on sale",
	Result_Product_LackOfQuantity:       "Lack of Quantity",

	Result_DBError:           "Internal DB error",
	Result_DBNotExistItem:    "Not exist item",
	Result_DBNotExistProduct: "Not exist Product",
	Result_DBNotExistAuction: "Not exist Auction",

	Result_TokenError:               "Internal Token error",
	Result_TokenERC721CreateError:   "ERC721 create error",
	Result_TokenERC721BurnError:     "ERC721 burn error",
	Result_TokenERC721TransferError: "ERC721 transfer error",
	Result_Reused_Txhash:            "Reused hash",

	Result_Auc_Product_Requiredtitle:        "Require product title",
	Result_Auc_Product_RequireDescription:   "Require description",
	Result_Auc_Product_RequireMediaOriginal: "Require Media original",
	Result_Auc_Product_RequireMediaThumnail: "Require Media thumnail",
	Result_Auc_Product_RequireOwnerInfo:     "Require owner info",
	Result_Auc_Product_RequireCreatorInfo:   "Require creator info",
	Result_Auc_Product_RequireProductId:     "Require product id",
	Result_Auc_Product_RequireIPOwnerShip:   "Require ip ownership",

	Result_Auc_Auction_RequireBidStartAmount: "Require bid start amount",
	Result_Auc_Auction_RequireBidUnit:        "Require bid unit",
	Result_Auc_Auction_RequireStartTs:        "Require bid start time",
	Result_Auc_Auction_RequireEndTs:          "Require bid end time",
	Result_Auc_Auction_RequireRound:          "Require bid round",
	Result_Auc_Auction_RequireActiveState:    "Require bid active state",
	Result_Auc_Auction_RequireProductId:      "Require bid product id",
	Result_Auc_Auction_RequireAucId:          "Require auc id",
	Result_Auc_Auction_NotPeriod:             "Not auction period",
	Result_Auc_Auction_NotOverYet:            "Auction is not over yet",

	Result_Auc_Bid_RequireAucId:         "Require bid auc id",
	Result_Auc_Bid_RequireWalletAddress: "Require bid wallet address",
	Result_Auc_Bid_RequireAmount:        "Require bid amount",
	Result_Auc_Bid_InvalidWalletAddress: "Invalid wallet address",
	Result_Auc_Bid_RequireDepositTxHash: "Require deposit tx hash",
	Result_Auc_Bid_AlreadyBestAttendee:  "Alreay best bid attendee",
	Result_Auc_Bid_NotBestBidAmount:     "You not best bid amount",
	Result_Auc_Bid_RequireDeposit:       "Require deposit send",
	Result_Auc_Bid_AlreadyDeposit:       "Alreay submit deposit",
	Result_Auc_Bid_NotWinner:            "Not winner",
	Result_Auc_Bid_RequireDepoistAgree:  "Consent for deposit payment is required",
	Result_Auc_Bid_RequirePrivacyAgree:  "Consent to the privacy policy is required",

	Result_Auc_Notice_RequireNotice: "At least one notice is required.",
	Result_Auc_Notice_RequireId:     "I need the id of the notice",

	Result_Auth_RequireMessage:    "Message is required",
	Result_Auth_RequireSign:       "Sign info is required",
	Result_Auth_InvalidLoginInfo:  "Invalid login info",
	Result_Auth_DontEncryptJwt:    "Auth token create fail",
	Result_Auth_InvalidJwt:        "Invalid jwt token",
	Result_Auth_InvalidWalletType: "Invalid wallet type",
}
