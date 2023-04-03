package main

import (
	"ecommerce/domain/order"
	orderImpl "ecommerce/domain/order/impl"
	"ecommerce/domain/user"
	userImpl "ecommerce/domain/user/impl"
	"ecommerce/pkg/database"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	database.InitDatabase()

	userRepo := userImpl.NewRepo(database.DB)
	orderRepo := orderImpl.NewRepo(database.DB)

	generateCSV(userRepo, orderRepo)
}

func generateCSV(userRepo user.Repository, orderRepo order.Repository) {
	allOrders, err := orderRepo.GetAllOrders()
	if err != nil {
		fmt.Println("error get orders by ids", err)
		return
	}
	userIds := make([]int, 0)
	orderByUserIds := make(map[int][]*order.Order)
	for _, po := range allOrders {
		userIds = append(userIds, po.CustomerId)

		if _, ok := orderByUserIds[po.CustomerId]; !ok {
			orderByUserIds[po.CustomerId] = make([]*order.Order, 0)
		}
		orderByUserIds[po.CustomerId] = append(orderByUserIds[po.CustomerId], &po)
	}

	users, err := userRepo.GetByIds(userIds)
	if err != nil {
		fmt.Println("error get users by ids", err)
		return
	}
	userIdsMap := make(map[int]*user.User)
	for _, u := range users {
		userIdsMap[u.ID] = &u
	}

	csvData := [][]string{
		{"order_id", "customer_name", "order_date", "total_price", "status"},
	}
	for _, o := range allOrders {
		u := userIdsMap[o.CustomerId]

		rowData := []string{
			strconv.Itoa(o.ID), u.Name, o.CreatedAt.Format(time.RFC3339), strconv.Itoa(o.TotalAmount), string(o.Status),
		}
		csvData = append(csvData, rowData)
	}

	now := time.Now()

	fileName := fmt.Sprintf("order-data-%s.csv", now.Format(time.RFC3339))
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("failed creating file: ", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	for _, row := range csvData {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvFile.Close()
}
