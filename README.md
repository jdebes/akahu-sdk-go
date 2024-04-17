# Akahu API Client Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/jdebes/akahu-sdk-go/akahu.svg)](https://pkg.go.dev/github.com/jdebes/akahu-sdk-go/akahu)

This is an unofficial and incomplete Go SDK for the [Akahu API](https://developers.akahu.nz/). The SDK provides a client library for interacting with the Akahu API as well as some utilities for the authentication flow.

## Getting Started

The easiest way to get access to the Akahu API is to create a personal app. More details can be found [here](https://developers.akahu.nz/docs/personal-apps).

If you're using a personal app, you can find your **App Token** and **User Token** [here](https://my.akahu.nz/developers).

## Installation

```
go get github.com/jdebes/akahu-sdk-go/akahu
```

## Usage

```
import "github.com/jdebes/akahu-sdk-go/akahu"
```

Create a client with your app token, app secret, and redirect URI:

```go
appToken := os.Getenv("YOUR_AKAHU_APP_TOKEN")
userToken := os.Getenv("YOUR_AKAHU_USER_TOKEN")
appSecret := os.Getenv("YOUR_AKAHU_APP_SECRET")

client := akahu.NewClient(nil, appToken, appSecret, "https://example.com/auth/akahu")
```

Then go ahead and query the API for a connected user:

```go
accounts, resp, err := client.Accounts.List(context.TODO(), "USER_ACCESS_TOKEN")
```

### More Examples

Take a look [here](https://github.com/jdebes/akahu-sdk-go/tree/main/example) for more examples, that demonstrate how to use the SDK.

The tests also provide a good overview.

## Status

The SDK has everything you need to complete the OAuth2 authorization flow and interact with the Akahu API.

However not all endpoints are implemented yet. Here is a rough list of supported endpoints:

- Accounts (complete)
- Auth (complete)
- Connections (complete)
- Webhooks (complete)
- Me (complete)
- Transactions (incomplete)
    - [x] Get transactions
    - [x] Get pending transactions
    - [x] Get a transaction by ID
    - [x] Get transactions by IDs
    - [ ] Get transactions by account 
    - [ ] Get pending transactions by account

See Akahu's full API reference [here](https://developers.akahu.nz/docs).