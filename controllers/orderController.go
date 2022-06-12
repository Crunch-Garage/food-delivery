package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	_ = json.NewDecoder(r.Body).Decode(&order)

	if order.ProfileID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("UserID field is required")
		return
	}

	if order.Delivery_address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Delivery address field is required")
		return
	}

	if order.OrderItem == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order Item(s) field is required")
		return
	}

	// get sum of order items
	total_price_sum := 0.0
	for _, orderItem := range order.OrderItem {
		total_price_sum += orderItem.Unit_price * float64(orderItem.Quantity)
	}

	// get estimated delivery charge from calaculate delivery charge based on location ai
	delivery_charge_estimate := 102.0

	// get total amount of order items and deleivery charge
	total_amount_sum := total_price_sum + delivery_charge_estimate

	order_ := models.Order{
		ProfileID:        order.ProfileID,
		Delivery_address: order.Delivery_address,
		Order_status:     "PENDING",
		Order_Date:       time.Now().Format(time.RFC3339),
		Total_price:      total_price_sum,
		Delivery_charge:  delivery_charge_estimate,
		Total_amount:     total_amount_sum,
	}

	createdOrder := database.DB.Create(&order_)
	err := createdOrder.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	for _, food_item := range order.OrderItem {
		order_item_ := models.OrderItem{
			Quantity:     food_item.Quantity,
			Unit_price:   food_item.Unit_price,
			OrderID:      int(order_.ID),
			FoodID:       food_item.FoodID,
			RestaurantID: food_item.RestaurantID,
		}

		database.DB.Create(&order_item_)
	}

	var order_item []models.OrderItem

	database.DB.Model(&order_).Related(&order_item)
	order_.OrderItem = order_item

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order_)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var order models.Order
	var orderItem []models.OrderItem
	var invoice models.Invoice

	database.DB.First(&order, id)
	database.DB.Model(&order).Related(&orderItem)
	database.DB.Model(&order).Related(&invoice)

	if order.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	var orderItemHolder []map[string]interface{}

	for i, _ := range orderItem {

		var food models.Food

		database.DB.Model(&orderItem[i]).Related(&food)

		foodData := map[string]interface{}{
			"id":          food.ID,
			"name":        food.Name,
			"price":       food.Price,
			"image":       food.Food_image,
			"description": food.Description,
			"status":      food.Status,
		}

		orderItemData := map[string]interface{}{
			"id":          orderItem[i].ID,
			"quantity":    orderItem[i].Quantity,
			"unit_price":  orderItem[i].Unit_price,
			"food":        foodData,
			"total_price": orderItem[i].Unit_price * float64(orderItem[i].Quantity),
		}

		orderItemHolder = append(orderItemHolder, orderItemData)
	}

	/**invoice interface*/
	invoiceData := map[string]interface{}{
		"id":             invoice.ID,
		"payment_date":   invoice.Payment_date,
		"payment_status": invoice.Payment_status,
		"payment_method": invoice.Payment_method,
		"payment_amount": invoice.Amount,
	}

	orderData := map[string]interface{}{
		"id":               order.ID,
		"CreatedAt":        order.CreatedAt,
		"UpdatedAt":        order.UpdatedAt,
		"customer_id":      order.ProfileID,
		"order_items":      orderItemHolder,
		"delivery_address": order.Delivery_address,
		"order_status":     order.Order_status,
		"driver_id":        order.DriverID,
		"order_date":       order.Order_Date,
		"total_price":      order.Total_price,
		"delivery_charge":  order.Delivery_charge,
		"total_amount":     order.Total_amount,
		"payment_details":  invoiceData,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderData)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var order models.Order
	var dbOrder models.Order
	var orderItem []models.OrderItem

	database.DB.First(&dbOrder, id)
	database.DB.Model(&dbOrder).Related(&orderItem)

	if dbOrder.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&order)

	if order.ProfileID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("customer id field is required")
		return
	}

	if order.Delivery_address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Delivery address field is required")
		return
	}

	if order.OrderItem == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order Item(s) field is required")
		return
	}

	// get sum of order items
	total_price_sum := 0.0
	for _, orderItem := range order.OrderItem {
		total_price_sum += orderItem.Unit_price * float64(orderItem.Quantity)
	}

	// get estimated delivery charge from calaculate delivery charge based on location ai
	delivery_charge_estimate := 102.0

	// get total amount of order items and deleivery charge
	total_amount_sum := total_price_sum + delivery_charge_estimate

	dbOrder.ProfileID = order.ProfileID
	dbOrder.Delivery_address = order.Delivery_address
	dbOrder.Order_status = order.Order_status
	dbOrder.Order_Date = time.Now().Format(time.RFC3339)
	dbOrder.Total_price = total_price_sum
	dbOrder.Delivery_charge = delivery_charge_estimate
	dbOrder.Total_amount = total_amount_sum

	updated_order := database.DB.Save(&dbOrder)
	err := updated_order.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	for i, _ := range order.OrderItem {

		UpdateOrderItemFunc(order.OrderItem[i], id)

	}

	var order_item []models.OrderItem

	database.DB.Model(&dbOrder).Related(&order_item)
	dbOrder.OrderItem = order_item

	// var orderItemHolder []map[string]interface{}

	// for i, _ := range order.OrderItem {

	// 	var food models.Food

	// 	database.DB.Model(&order.OrderItem[i]).Related(&food)

	// 	foodData := map[string]interface{}{
	// 		"id":          food.ID,
	// 		"name":        food.Name,
	// 		"price":       food.Price,
	// 		"image":       food.Food_image,
	// 		"description": food.Description,
	// 		"status":      food.Status,
	// 	}

	// 	orderItemData := map[string]interface{}{
	// 		"id":          order.OrderItem[i].ID,
	// 		"quantity":    order.OrderItem[i].Quantity,
	// 		"unit_price":  order.OrderItem[i].Unit_price,
	// 		"food":        foodData,
	// 		"total_price": order.OrderItem[i].Unit_price * float64(order.OrderItem[i].Quantity),
	// 	}

	// 	orderItemHolder = append(orderItemHolder, orderItemData)
	// }

	// orderData := map[string]interface{}{
	// 	"id":               dbOrder.ID,
	// 	"customer_id":      dbOrder.ProfileID,
	// 	"order_items":      orderItemHolder,
	// 	"delivery_address": dbOrder.Delivery_address,
	// 	"order_status":     dbOrder.Order_status,
	// 	"driver_id":        dbOrder.DriverID,
	// 	"order_date":       dbOrder.Order_Date,
	// 	"total_price":      dbOrder.Total_price,
	// 	"delivery_charge":  dbOrder.Delivery_charge,
	// 	"total_amount":     dbOrder.Total_amount,
	// }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbOrder)
}
