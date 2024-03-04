package batch

import (
	"context"
	"go-find-by-id-conncurrency/model"
)

type Batch struct {
	BatchInterface BatchInterface
}

func NewBatch(batchInterface BatchInterface) *Batch {
	return &Batch{
		BatchInterface: batchInterface,
	}
}

func (b *Batch) ReadByID(ctx context.Context, id uint) model.Campaign {
	return b.BatchInterface.ReadByID(ctx, id)
}
