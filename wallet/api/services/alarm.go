package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/databases"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/wallet/api/utils"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	resty "gopkg.in/resty.v1"
)

/*

JSON structure that is required when sending push notification

{
  "notifications": [
    {
      "tokens": ["dRtiaZY5JzI:APA91bEEVWV7AbszQutnuAlFZfn9aXucZUCo_sbTltmKB7F1_l3n2TtlR31HmPx04xSw6kl0V0Fafjn4koAqydPMR8heKv8n_9Zr0bLIjsjVzfOFXC2jjWbfTxhURSDTnW0_Zvh1s6J5"],
      "platform": 2,
      "message": "you received atom with txid 06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
      "title": "Received 11.434532Atom",
      "data": {"notifyto" : "cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7","txid" : "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827"}
    }
  ]
}
*/

// PushNotification receives an address from push-tx-parser and
// sends push notification to its respective device
func PushNotification(db *pg.DB, cf *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	// lower case
	address = strings.ToLower(address)

	// query account information
	var account models.Account
	account, _ = databases.QueryAccount(w, db, address)

	// return when data is empty
	if account.AlarmToken == "" {
		return
	}

	// check user's alarm status
	if !account.AlarmStatus {
		u.Result(w, false, "user's alarm status is false")
		return
	}

	// send push notification
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{	
  "notifications": [	
    {	
      "tokens": [
      	"dRtiaZY5JzI:APA91bEEVWV7AbszQutnuAlFZfn9aXucZUCo_sbTltmKB7F1_l3n2TtlR31HmPx04xSw6kl0V0Fafjn4koAqydPMR8heKv8n_9Zr0bLIjsjVzfOFXC2jjWbfTxhURSDTnW0_Zvh1s6J5"
      ],	
      "platform": 2,	
      "title": "Received 11.434532Atom",	
      "message": "you received atom with txid 06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",	
      "data": {
      	"notifyto" : "cosmos188z9th39f54sqmexvgj5wvjg4sus9qmm9w665e",
      	"txid" : "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827"
      }	
    }	
  ]	
}`).
		Post(cf.Web.PushServerURL)

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println(resp)

	u.Result(w, true, "successfully sent push notification")
	return
}

func PushTest(db *pg.DB, cf *config.Config, w http.ResponseWriter, r *http.Request) {

	payload := models.Payload{
		NotifyTo: "cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7",
		Txid:     "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
	}

	var notification []models.Notification
	tempNotification := models.Notification{
		Tokens:   []string{"APA91bEEVWV7AbszQutnuAlFZfn9aXucZUCo_sbTltmKB7F1_l3n2TtlR31HmPx04xSw6kl0V0Fafjn4koAqydPMR8heKv8n_9Zr0bLIjsjVzfOFXC2jjWbfTxhURSDTnW0_Zvh1s6J5"},
		Platform: 2,
		Message:  "Hello World",
		Title:    "Title",
		Data:     payload,
	}
	notification = append(notification, tempNotification)

	pushNotification := models.PushNotification{
		Notifications: notification,
	}

	fmt.Println("notification", notification)
	fmt.Println("payload", payload)
	fmt.Println("pushNotification", pushNotification)

	u.Respond(w, pushNotification)
	return
}
