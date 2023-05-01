package main

import (
	"fmt"
	"os"

	"github.com/jdebes/akahu-sdk-go/akahu"
)

func main() {
	appToken := os.Getenv("AKAHU_APP_TOKEN")
	appSecret := os.Getenv("AKAHU_APP_SECRET")

	client := akahu.NewClient(nil, appToken, appSecret, "https://example.com/auth/akahu")

	options := akahu.AuthorizationURLOptions{}
	authUrl := client.Auth.BuildAuthorizationURL(options)

	fmt.Println(authUrl)
}
