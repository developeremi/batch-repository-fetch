package batch

import (
	"context"
	"go-find-by-id-conncurrency/model"
)

type BatchInterface interface {
	ReadByID(ctx context.Context, id uint) model.Campaign
}
