// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "drdgvhbh/discordbot/internal/entity"
import mock "github.com/stretchr/testify/mock"

// StockRepository is an autogenerated mock type for the StockRepository type
type StockRepository struct {
	mock.Mock
}

// UpsertStock provides a mock function with given fields: stock
func (_m *StockRepository) UpsertStock(stock *entity.Stock) (*entity.Stock, error) {
	ret := _m.Called(stock)

	var r0 *entity.Stock
	if rf, ok := ret.Get(0).(func(*entity.Stock) *entity.Stock); ok {
		r0 = rf(stock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Stock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*entity.Stock) error); ok {
		r1 = rf(stock)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}