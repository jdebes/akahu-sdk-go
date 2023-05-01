package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jdebes/akahu-sdk-go/akahu"
)

func main() {
	appToken := os.Getenv("AKAHU_APP_TOKEN")
	userToken := os.Getenv("AKAHU_USER_TOKEN")
	appSecret := os.Getenv("AKAHU_APP_SECRET")

	client := akahu.NewClient(nil, appToken, appSecret, "https://example.com/auth/akahu")
	accounts, resp, err := client.Accounts.List(context.TODO(), userToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)

	for i, account := range accounts {
		fmt.Printf("%v: %v\n", i+1, account.Name)
	}
}
