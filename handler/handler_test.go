package handler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var h Handler = *NewHandler(2 * time.Second)

func init() {
	// store order one
	h.storeOrder("ABC123", []OrderDetails{
		{
			Latitude:  123.11,
			Longitude: 10.99,
		},
		{
			Latitude:  709.12,
			Longitude: 25.9,
		},
		{
			Latitude:  502.12,
			Longitude: 55.2,
		},
	})

	// store order two
	h.storeOrder("ABC124", []OrderDetails{
		{
			Latitude:  90.11,
			Longitude: 710.99,
		},
		{
			Latitude:  199.00,
			Longitude: 35.7,
		},
		{
			Latitude:  119.00,
			Longitude: 135.2,
		},
	})
}

type HandlerTestCase struct {
	Case             string
	req              interface{}
	ExpectedResponse interface{}
}

type request struct {
	orderNumber string
	max         int
}

func Test_GetOrder(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Case: "Test successful order and return history with 2 length",
			req: request{
				orderNumber: "ABC124",
				max:         2,
			},
			ExpectedResponse: OrderResponse{
				OrderID: "ABC124",
				History: []OrderDetails{
					{Latitude: 90.11, Longitude: 710.99},
					{Latitude: 199.00, Longitude: 35.7},
				},
			},
		},
		{
			Case: "Test invalid order number",
			req: request{
				orderNumber: "ABC125",
			},
			ExpectedResponse: "not found",
		},
	}

	for _, tc := range testCases {
		getOrder(t, &tc)
	}
}

func Test_Delete(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Case: "Test successful delete an order",
			req: request{
				orderNumber: "ABC123",
				max:         2,
			},
			ExpectedResponse: "not found",
		},
	}

	for _, tc := range testCases {
		deleteOrder(t, &tc)
	}
}

func getOrder(t *testing.T, tc *HandlerTestCase) {
	r := tc.req.(request)
	resp, err := h.GetOrder(r.orderNumber, r.max)
	if err != nil {
		assert.EqualError(t, err, tc.ExpectedResponse.(string))
	}

	if resp != nil {
		assert.Equal(t, tc.ExpectedResponse, *resp)
		assert.Equal(t, len(resp.History), r.max)
	}
}

func deleteOrder(t *testing.T, tc *HandlerTestCase) {
	r := tc.req.(request)
	resp, err := h.GetOrder(r.orderNumber, r.max)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(resp.History), r.max)

	h.DeleteOrder(r.orderNumber)

	getOrder(t, tc)

}

func Test_MakeOrder(t *testing.T) {
	orderNum := "test1"
	cases := []OrderDetails{
		{Latitude: 1, Longitude: 710.99},
		{Latitude: 2, Longitude: 35.7},
		{Latitude: 3, Longitude: 35.7},
		{Latitude: 4, Longitude: 35.7},
		{Latitude: 5, Longitude: 35.7},
		{Latitude: 6, Longitude: 35.7},
		{Latitude: 7, Longitude: 35.7},
		{Latitude: 8, Longitude: 35.7},
		{Latitude: 9, Longitude: 35.7},
	}

	wt := sync.WaitGroup{}

	for _, tc := range cases {
		wt.Add(1)
		go func() {
			h.MakeOrder(orderNum, tc.Latitude, tc.Longitude)
			wt.Done()
		}()
	}

	wt.Wait()

	history, ok := h.loadOrder(orderNum)
	assert.True(t, ok)

	assert.Equal(t, len(cases), len(history))
}
