package transfer_mock

import (
	"chuan.place/fp/transfer"
	mock "github.com/stretchr/testify/mock"
)

type TestObj struct {
	mock.Mock
}

func (o TestObj) Transfer(p transfer.Payload) (string, error) {
	args := o.Called(p)
	return args.Get(0).(string), args.Error(1)
}
