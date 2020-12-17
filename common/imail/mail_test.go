package imail

import (
	"github.com/jordan-wright/email"
	"log"
	"testing"
	"time"
)

func Test_Send(t *testing.T) {
	Init(&IOption{
		Service:  "smtp.126.com",
		AuthCode: "KZMAPPOVFWIJJSXC",
		Sender:   "p01platform@126.com",
	})

	err := Ins.Send(&email.Email{
		From:        "p01platform@126.com",
		To:          []string{"1296985852@qq.com"},
		Subject:     "this test",
		Text:        []byte("hello,hai gg"),
		HTML:        nil,
		Sender:      "",
		Headers:     nil,
		Attachments: nil,
		ReadReceipt: nil,
	}, time.Second*3)
	if err != nil {
		log.Fatal(err)
	}
}
