package main

import (
	_ "github.com/logpacker/PayPal-Go-SDK"
)

func main() {

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
