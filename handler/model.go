package handler

type OrderDetails struct {
	Latitude  string `json:"lat" binding:"required"`
	Longitude string `json:"lng" binding:"required"`
}

type OrderResponse struct {
	OrderID string         `json:"order_id"`
	History []OrderDetails `json:"history"`
}

type Order struct {
	OrderDetails       []OrderDetails
	ExpiredAtTimeStamp int64
}
