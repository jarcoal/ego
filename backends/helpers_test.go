// random test data from http://www.databasetestdata.com/
// hopefully it doesn't collide with any actual addresses :)
package backends

import (
	"github.com/jarcoal/ego"
	"net/mail"
)

func testAddresses() []*mail.Address {
	return []*mail.Address{
		{Name: "Nyasia Block", Address: "jade@austen.name"},
		{Name: "Erin Dare", Address: "otilia@hermina.io"},
		{Name: "Sandy Schulist", Address: "zane@anastacio.co.uk"},
		{Name: "Rocio Christiansen", Address: "retta.ankunding@fletcher.biz"},
		{Name: "Abigale Gleason", Address: "freida@orpha.info"},
		{Name: "Garland Spencer", Address: "corrine@remington.io"},
		{Name: "Tad Will", Address: "ludwig@paula.co.uk"},
		{Name: "Chad Ritchie", Address: "kathryne_ankunding@uriel.biz"},
		{Name: "Junius Boehm", Address: "baylee.fadel@ellis.info"},
		{Name: "Edison Kris", Address: "alyce_rutherford@gennaro.biz"},
	}
}

func testEmail() *ego.Email {
	addresses := testAddresses()

	return &ego.Email{
		To:       addresses[0:7],
		From:     addresses[8],
		ReplyTo:  addresses[9],
		Subject:  "Test Subject",
		HtmlBody: "<h1>Test Body</h1>",
		TextBody: "Test Body",
		Tags:     []string{"really", "important", "message"},
	}
}
