package main

import (
	"fmt"
	"log"
	"net/http"

	"chuan.place/fp/transfer"
)

type Resp struct {
	Data string
}

// Handler's responsiblity is to 1. call the respective helpers, 2. read and write
func transferHandle(transfer transfer.TransferFx, prepare transfer.PrepareFx) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow post
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		payload, err := prepare(r.URL.Query())
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		state, err := transfer(payload)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintf(w, "%s -> %s: %d [%s]\n", payload.Fr, payload.To, payload.Amount, state)
	}
}

// e2e test for this
func main() {
	http.HandleFunc("/transfer", transferHandle(transfer.Transfer, transfer.Prepare))
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
