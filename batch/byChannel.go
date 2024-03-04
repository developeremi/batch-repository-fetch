package batch

import (
	"context"
	"go-find-by-id-conncurrency/model"
	"gorm.io/gorm"
	"time"
)

type ByChannel struct {
	db     *gorm.DB
	buffer chan queue
}

func NewByChannel(db *gorm.DB) *ByChannel {
	// run worker
	b := &ByChannel{
		db:     db,
		buffer: make(chan queue, 10000),
	}
	go b.worker()
	return b
}

func (b *ByChannel) ReadByID(ctx context.Context, id uint) model.Campaign {

	// make private channel for this request
	promise := make(chan queueResponse)
	// send id and promise to buffer to fetch and return promise callback
	iqs := queue{
		ID:      id,
		Promise: promise,
	}

	select {
	case b.buffer <- iqs:
	default:
		return model.Campaign{}
	}

	select {
	case <-ctx.Done():

		return model.Campaign{}
	case result := <-promise:

		return result.campaign
	}
}

func (b *ByChannel) worker() {

	for {
		// execute every 50 millisecond
		<-time.After(time.Millisecond * 20)

		// if send parameter into buffer
		l := len(b.buffer)
		if l > 0 {

			// ids for fetch from database
			buffer := make([]queue, l)
			campaigns := make(map[uint]queueResponse)
			for i := 0; i < l; i++ {

				item := <-b.buffer // we have to send each buffer into array

				buffer[i] = item // assign random buffer item into array

				campaigns[item.ID] = queueResponse{} // prepare map model

			}

			var ids []uint
			for id, _ := range campaigns {
				ids = append(ids, id)
			}

			var results []model.Campaign
			b.db.Where("id IN ?", ids).Find(&results)

			for _, result := range results {
				response := campaigns[result.ID]
				response.campaign = result
				campaigns[result.ID] = response
			}

			for _, final := range buffer {
				response := campaigns[final.ID]
				final.Promise <- response
			}
		}

	}

}

// queue definitions
type queueResponse struct {
	campaign model.Campaign
}

type queue struct {
	ID      uint
	Promise chan queueResponse
}
