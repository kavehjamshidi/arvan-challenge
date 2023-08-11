package contract

import (
	"github.com/kavehjamshidi/arvan-challenge/domain"
)

type FileUploadQueue interface {
	Store(file domain.File) (err error)
}
