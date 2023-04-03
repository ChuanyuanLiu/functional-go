package transfer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Payload struct {
	Fr     string
	To     string
	Amount int
}

type TransferFx func(Payload) (string, error)

func Transfer(p Payload) (string, error) {
	// assuming we have a partner
	return "success", nil
}

type PrepareFx func(url.Values) (Payload, error)

func Prepare(params url.Values) (Payload, error) {
	amount, err := strconv.Atoi(params.Get("amount"))
	if err != nil {
		return Payload{}, fmt.Errorf("query `amount` should be an integer")
	}
	payload := Payload{
		Fr:     params.Get("fr"),
		To:     params.Get("to"),
		Amount: amount,
	}
	msgs := []string{}
	if payload.Amount == 0 {
		msgs = append(msgs, "query `amount` has default value")
	}
	if payload.Fr == "" {
		msgs = append(msgs, "query `fr` has default value")
	}
	if payload.To == "" {
		msgs = append(msgs, "query `to` has default value")
	}
	if len(msgs) != 0 {
		return Payload{}, fmt.Errorf(strings.Join(msgs, "\n"))
	}
	return payload, nil
}
