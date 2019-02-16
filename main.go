package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/logpacker/PayPal-Go-SDK"
)

func main() {
	r := gin.Default()
	r.GET("/", Index)
	r.POST("/payouts", Payouts)

	log.Fatal(r.Run())

	c := getPayPalClient()

	resp, err := c.CreateSinglePayout(paypalsdk.Payout{
		SenderBatchHeader: &paypalsdk.SenderBatchHeader{
			EmailSubject: "Subject",
		},
		Items: []paypalsdk.PayoutItem{
			paypalsdk.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      "test@chainlink.com",
				Amount: &paypalsdk.AmountPayout{
					Value:    "1.11",
					Currency: "USD",
				},
				Note:         "thanks for the secure data shaggy",
				SenderItemID: "1337",
			},
			paypalsdk.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      "tabba.ahmad-buyer@gmail.com",
				Amount: &paypalsdk.AmountPayout{
					Value:    ".50",
					Currency: "USD",
				},
				Note:         "Note",
				SenderItemID: "SenderItemID",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	d, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(d))

	resp, _ = c.GetPayout(resp.BatchHeader.PayoutBatchID)
	d, _ = json.MarshalIndent(resp, "", "  ")
	log.Println(string(d))
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{"message": "PayPal external adapter"})
}

func Payouts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Sent Payment to test@test.com"})
}

func getPayPalClient() *paypalsdk.Client {
	c, err := paypalsdk.NewClient(
		os.Getenv("PAYPAL_CLIENT_ID"),
		os.Getenv("PAYPAL_SECRET"),
		paypalsdk.APIBaseSandBox,
	)
	if err != nil {
		panic(err)
	}
	log.Println("Created paypal client with base url", c.APIBase)
	c.SetLog(os.Stdout)
	setToken(c)
	return c
}

func setToken(c *paypalsdk.Client) {
	if os.Getenv("PAYPAL_TOKEN") == "" {
		t, _ := c.GetAccessToken()
		log.Println("PAYPAL_TOKEN env variable not set.")
		log.Printf("Token %s\nExpires in %d", t.Token, t.ExpiresIn)
	} else {
		c.SetAccessToken(os.Getenv("PAYPAL_TOKEN"))
	}
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
