package dto

import "product-mall/internal/model"

type Order struct {
	ID         uint        `json:"id"`
	OrderNo    string      `json:"order_no"`
	UserID     uint        `json:"user_id"`
	AddressID  uint        `json:"address_id"`
	TotalPrice float64     `json:"total_price"`
	Status     int         `json:"status"`
	StatusText string      `json:"status_text"`
	Remark     string      `json:"remark,omitempty"`
	Items      []OrderItem `json:"items,omitempty"`
	CreatedAt  int64       `json:"created_at"`
}

type OrderItem struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Num         uint    `json:"num"`
	Price       float64 `json:"price"`
}

var statusTextMap = map[int]string{
	0: "待支付",
	1: "已支付",
	2: "已发货",
	3: "已完成",
	4: "已取消",
}

func BuildOrder(order model.Order, items []model.OrderItem) Order {
	orderDTO := Order{
		ID:         order.ID,
		OrderNo:    order.OrderNo,
		UserID:     order.UserID,
		AddressID:  order.AddressID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		StatusText: statusTextMap[order.Status],
		Remark:     order.Remark,
		CreatedAt:  order.CreatedAt.Unix(),
	}

	for _, item := range items {
		orderDTO.Items = append(orderDTO.Items, OrderItem{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Num:         item.Num,
			Price:       item.Price,
		})
	}
	if orderDTO.Items == nil {
		orderDTO.Items = make([]OrderItem, 0)
	}
	return orderDTO
}

func BuildOrders(orders []model.Order) []Order {
	result := make([]Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, BuildOrder(o, nil))
	}
	return result
}
