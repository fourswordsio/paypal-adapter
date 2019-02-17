package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/logpacker/PayPal-Go-SDK"
)

var ppc = getPayPalClient()

func main() {
	r := gin.Default()
	r.POST("/", Payouts)

	log.Fatal(r.Run())
	// d, _ := json.MarshalIndent(resp, "", "  ")
	// log.Println(string(d))

	// resp, _ = c.GetPayout(resp.BatchHeader.PayoutBatchID)
	// d, _ = json.MarshalIndent(resp, "", "  ")
	// log.Println(string(d))
}

func Payouts(c *gin.Context) {
	br := BridgeRequest{}
	if err := c.BindJSON(&br); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	_, err := createSinglePayout(br.Data.Email, br.Data.Amount)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(200, BridgeResponse{
			JobRunID: br.JobRunID,
			Status:   "pending_bridge",
			Pending:  true,
		})
	}
}

func createSinglePayout(email, amount string) (*paypalsdk.PayoutResponse, error) {
	return ppc.CreateSinglePayout(paypalsdk.Payout{
		SenderBatchHeader: &paypalsdk.SenderBatchHeader{
			EmailSubject: "Subject",
		},
		Items: []paypalsdk.PayoutItem{
			paypalsdk.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      email,
				Amount: &paypalsdk.AmountPayout{
					Value:    amount,
					Currency: "USD",
				},
				Note:         "thanks for the secure data shaggy",
				SenderItemID: "1337",
			},
		},
	})

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
	Data        Payout `json:"data"`
	ResponseURL string `json:"responseURL"`
}

type Payout struct {
	Email  string `json:"email"`
	Amount string `json:"amount"`
}

type BridgeResponse struct {
	JobRunID string `json:"jobRunID"`
	Data     gin.H  `json:"data"`
	Status   string `json:"status"`
	Error    string `json:"error"`
	Pending  bool   `json:"pending"`
}
