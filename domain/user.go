package domain

type User struct {
	ID                  string
	Usage               uint64
	QuotaStartTimestamp uint64
	QuotaEndTimestamp   uint64
}
