package queue

import (
	"github.com/kavehjamshidi/arvan-challenge/domain"
	"github.com/kavehjamshidi/arvan-challenge/queue/contract"
)

type noOpQueue struct{}

func NewNoOpQueue() contract.FileUploadQueue {
	return noOpQueue{}
}

func (n noOpQueue) Store(file domain.File) error {
	return nil
}
