package services

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
)

func PushNotification(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Alarm Test")

	// 계정 확인

	// 알람 푸쉬

}
