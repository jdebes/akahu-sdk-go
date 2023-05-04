package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jdebes/akahu-sdk-go/akahu"
)

func main() {
	appToken := os.Getenv("AKAHU_APP_TOKEN")
	appSecret := os.Getenv("AKAHU_APP_SECRET")
	userToken := os.Getenv("AKAHU_USER_TOKEN")

	client := akahu.NewClient(nil, appToken, appSecret, "https://example.com/auth/akahu")

	startTime, _ := time.Parse(time.DateOnly, "2022-01-01")
	endTime := time.Now()

	transactions, resp, _ := client.Transactions.List(context.TODO(), userToken, startTime, endTime)

	fmt.Printf("Response status: %d\n", resp.StatusCode)

	for i, transaction := range transactions {
		fmt.Printf("%v: %v\n", i+1, transaction.Id)
	}
}
