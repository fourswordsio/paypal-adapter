package main

import (
	"os"

	"github.com/logpacker/PayPal-Go-SDK"
)

func main() {

	c := getPayPalClient()

}

func getPayPalClient() PayPal-Go-SDK {
	c, err := paypal.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
	c.SetLog(os.Stdout)

	accessToken, err := c.GetAccessToken()
}

type BridgeRequest struct {
	JobRunID    string `json:"id"`
	data        int
	ResponseURL string `json:"responseURL"`
}

type BridgeResponse struct {
	JobRunID string `json:"jobRunID"`
	data     int
	Status   string `json:"status"`
	Error    string `json:"error"`
	Pending  bool   `json:"pending"`
}
