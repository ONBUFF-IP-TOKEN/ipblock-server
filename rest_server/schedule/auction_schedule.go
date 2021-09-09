package schedule

import (
	"sync"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

var gAuctionScheduler *AuctionScheduler
var once sync.Once

type AuctionScheduler struct {
	auctions []*context_auc.AucAuction
	mutex    *sync.Mutex
}

func GetScheduler() *AuctionScheduler {
	once.Do(func() {
		gAuctionScheduler = new(AuctionScheduler)
		gAuctionScheduler.mutex = new(sync.Mutex)
		gAuctionScheduler.auctions = make([]*context_auc.AucAuction, 0)
		gAuctionScheduler.Run()
	})

	return gAuctionScheduler
}

func (o *AuctionScheduler) ResetAuctionScheduler() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	var err error
	o.auctions, err = model.GetDB().GetAucAuctionListForSchedule()
	if err != nil {
		log.Error("ResetAuctionScheduler DB error ", err)
	}
}

func (o *AuctionScheduler) Run() error {
	go func() {
		ticker := time.NewTicker(time.Duration(5) * time.Second)

		for {
			// 경매 시간에 맞게 경매 auc_state 변경처리
			o.CheckAuctionState()
			<-ticker.C
		}
	}()

	return nil
}

func (o *AuctionScheduler) CheckAuctionState() {
	if len(o.auctions) == 0 {
		return
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	curT := datetime.GetTS2MilliSec()
	bChange := false
	for _, auction := range o.auctions {
		if auction.AucStartTs < curT && auction.AucEndTs > curT {
			// 경매 중으로 변경
			if auction.AucState != context_auc.Auction_auc_state_start {
				log.Info("change schedule : ", auction.Id, " auc_state:", context_auc.Auction_auc_state_start)
				bChange = true
				auction.AucState = context_auc.Auction_auc_state_start
				model.GetDB().UpdateAucAuctionAucState(auction.Id, context_auc.Auction_auc_state_start, false)
				model.GetDB().DeleteAuctionCache(auction.Id)
			}
		} else if auction.AucEndTs < curT {
			// 경매 종료
			if auction.AucState != context_auc.Auction_auc_state_finish {
				log.Info("change schedule : ", auction.Id, " auc_state:", context_auc.Auction_auc_state_finish)
				bChange = true
				auction.AucState = context_auc.Auction_auc_state_finish
				model.GetDB().UpdateAucAuctionAucState(auction.Id, context_auc.Auction_auc_state_finish, false)
				model.GetDB().DeleteAuctionCache(auction.Id)
			}
		}
	}

	if bChange {
		model.GetDB().DeleteAuctionList()
	}
}
