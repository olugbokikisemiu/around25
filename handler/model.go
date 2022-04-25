package handler

type OrderDetails struct {
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
}

type OrderResponse struct {
	OrderID string         `json:"order_id"`
	History []OrderDetails `json:"history"`
}

type Order struct {
	OrderDetails       []OrderDetails
	ExpiredAtTimeStamp int64
}
