package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	transfer "chuan.place/fp/transfer"
	transferMock "chuan.place/fp/transfer/mock"
)

func Test_transferHandle(t *testing.T) {

	// We only need to test the happy and sad path
	tests := []struct {
		name       string
		params     string
		mockErr    error
		mockArg    transfer.Payload
		mockReturn string
		expectRes  string
	}{
		{
			name:       "valid payload, successful transfer",
			params:     "amount=2&fr=a&to=b",
			mockArg:    transfer.Payload{Fr: "a", To: "b", Amount: 2},
			mockReturn: "success",
			expectRes:  "a -> b: 2 [success]\n",
		},
		{
			name:      "valid payload, failed transfer",
			params:    "amount=2&fr=a&to=b",
			mockArg:   transfer.Payload{Fr: "a", To: "b", Amount: 2},
			mockErr:   fmt.Errorf("error"),
			expectRes: "error\n",
		},
		{
			name:      "invalid payload",
			expectRes: "query `amount` should be an integer\n",
		},
	}
	for _, tt := range tests {
		transferMockObj := transferMock.TestObj{}
		// only mocking transfer and not prepare
		transferMockObj.On("Transfer", tt.mockArg).Return(tt.mockReturn, tt.mockErr).Once()
		handle := transferHandle(transferMockObj.Transfer, transfer.Prepare)

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/transfer?"+tt.params, nil)
			rec := httptest.NewRecorder()
			handle(rec, req)
			assert.Equal(t, tt.expectRes, rec.Body.String())
		})
	}

	// edge case
	t.Run("don't allow GET", func(t *testing.T) {
		transferMockObj := transferMock.TestObj{}
		handle := transferHandle(transferMockObj.Transfer, transfer.Prepare)
		req := httptest.NewRequest(http.MethodGet, "/transfer", nil)
		rec := httptest.NewRecorder()
		handle(rec, req)
		assert.Equal(t, "Method not allowed\n", rec.Body.String())
	})
}
