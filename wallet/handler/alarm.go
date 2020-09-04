package handler

import "net/http"

/*

JSON structure that is required when sending push notification

{
  "notifications": [
    {
      "tokens": ["dRtiaZY5JzI:APA91bEEVWV7AbszQutnuAlFZfn9aXucZUCo_sbTltmKB7F1_l3n2TtlR31HmPx04xSw6kl0V0Fafjn4koAqydPMR8heKv8n_9Zr0bLIjsjVzfOFXC2jjWbfTxhURSDTnW0_Zvh1s6J5"],
      "platform": 2,
      "message": "you received atom with txid 06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
      "title": "Received 11.434532Atom",
      "data": {
		  "notifyto" : "cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7",
		  "txid" : "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
		  "type" : "send or receive"
		}
    }
  ]
}
*/

// PushNotification receives an address from push-tx-parser and
// sends push notification to its respective device
func PushNotification(w http.ResponseWriter, r *http.Request) {
	// var nrp models.NotificationReceivePayload

	// // get post data from request
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&nrp)
	// if err != nil {
	// 	errors.ErrBadRequest(w, http.StatusBadRequest)
	// 	return
	// }

	// // lower case
	// nrp.From = strings.ToLower(nrp.From)
	// nrp.To = strings.ToLower(nrp.To)
	// nrp.Txid = strings.ToLower(nrp.Txid)

	// // query account information
	// var account models.Account
	// account, _ = databases.QueryAccount(w, db, nrp.To)

	// // return when data is empty
	// if account.AlarmToken == "" {
	// 	return
	// }

	// // check user's alarm status
	// if !account.AlarmStatus {
	// 	u.Result(w, false, "user's alarm status is false")
	// 	return
	// }

	// // push notification payload
	// var notification []models.Notifications
	// tempNotification := models.Notifications{
	// 	Tokens:   []string{account.AlarmToken},
	// 	Platform: 2,
	// 	Title:    models.AlarmTitle + nrp.Amount,
	// 	Message:  models.AlarmMessage + nrp.Amount,
	// 	Data: models.Data{
	// 		NotifyTo: nrp.To,
	// 		Txid:     nrp.Txid,
	// 	},
	// }
	// notification = append(notification, tempNotification)

	// notificationPayload := models.NotificationPayload{
	// 	Notifications: notification,
	// }

	// // send push notification
	// _, err = resty.R().
	// 	SetHeader("Content-Type", "application/json").
	// 	SetBody(notificationPayload).
	// 	Post(cf.Alarm.PushServerURL)
	// if err != nil {
	// 	errors.ErrInternalServer(w, http.StatusInternalServerError)
	// 	return
	// }

	// u.Result(w, true, "successfully sent push notification")
	return
}
