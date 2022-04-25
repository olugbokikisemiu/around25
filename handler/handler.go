package handler

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Handler struct {
	lock     sync.RWMutex
	data     map[string]Order
	cacheTTL time.Duration
}

func NewHandler(cacheTTL time.Duration) *Handler {
	return &Handler{data: make(map[string]Order), cacheTTL: cacheTTL}
}

func (h *Handler) MakeOrder(orderNumber string, lat float64, lng float64) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	order := OrderDetails{
		Latitude:  lat,
		Longitude: lng,
	}

	t := time.Now().Add(h.cacheTTL)

	orderData := Order{
		ExpiredAtTimeStamp: t.Unix(),
	}

	res, ok := h.data[orderNumber]
	if !ok {
		orderDetails := []OrderDetails{}

		orderDetails = append(orderDetails, order)
		orderData.OrderDetails = orderDetails
	} else {
		res.OrderDetails = append(res.OrderDetails, order)
		orderData.OrderDetails = res.OrderDetails
	}

	h.data[orderNumber] = orderData

}

func (h *Handler) GetOrder(orderNumber string, max int) (resp *OrderResponse, err error) {
	res, ok := h.loadOrder(orderNumber)
	if !ok {
		err = fmt.Errorf("not found")
		return
	}

	resp = &OrderResponse{
		OrderID: orderNumber,
	}

	if max == 0 || len(res) <= max {
		resp.History = res
		return

	}

	counter := 0

	details := []OrderDetails{}

	for counter < max {
		order := OrderDetails{
			Latitude:  res[counter].Latitude,
			Longitude: res[counter].Longitude,
		}
		details = append(details, order)
		counter++
	}

	resp.History = details

	return
}

func (h *Handler) DeleteOrder(orderNumber string) (err error) {
	_, ok := h.loadOrder(orderNumber)
	if !ok {
		err = fmt.Errorf("not found")
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.data, orderNumber)

	return
}

func (h *Handler) loadOrder(orderNumber string) (resp []OrderDetails, ok bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	res, ok := h.data[orderNumber]
	if !ok {
		ok = false
		return
	}

	ok = true
	resp = res.OrderDetails
	return
}

func (h *Handler) storeOrder(orderNumber string, data []OrderDetails) {
	h.lock.Lock()
	defer h.lock.Unlock()

	t := time.Now().Add(h.cacheTTL)

	orderData := Order{
		OrderDetails:       data,
		ExpiredAtTimeStamp: t.Unix(),
	}

	h.data[orderNumber] = orderData
}

func (h *Handler) RemoveExpiredOrder() {
	log.Println("runner for removeExpiredOrder!!")

	for k := range h.data {
		if h.data[k].ExpiredAtTimeStamp < time.Now().Unix() {
			delete(h.data, k)
		}
	}
}
