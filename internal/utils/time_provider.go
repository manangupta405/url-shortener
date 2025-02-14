package utils

import "time"

//go:generate mockery --name=TimeProvider --output=./mocks
type TimeProvider interface {
	Now() time.Time
}

type timeProviderImpl struct{}

func NewTimeProvider() TimeProvider {
	return &timeProviderImpl{}
}

func (t *timeProviderImpl) Now() time.Time {
	return time.Now()
}
