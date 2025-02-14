package utils

import (
	"github.com/aidarkhanov/nanoid"
)

//go:generate mockery --name=NanoIDGenerator --output=./mocks
type NanoIDGenerator interface {
	Generate() (string, error)
}

type nanoidGeneratorImpl struct {
	size int
}

func NewNanoIDGenerator(size int) NanoIDGenerator {
	return &nanoidGeneratorImpl{size: size}
}

func (i *nanoidGeneratorImpl) Generate() (string, error) {
	return nanoid.Generate(nanoid.DefaultAlphabet, i.size)
}
