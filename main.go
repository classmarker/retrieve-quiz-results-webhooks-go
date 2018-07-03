package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
)

var ClassmarkerWebhookSecret = "YOUR_CLASSMARKER_WEBHOOK_SECRET_PHRASE"

func WebHook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("received webhook request")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonData := string(body)

	//ClassMarker sent signature to check against
	headerHmacSignature := r.Header.Get("X-Classmarker-Hmac-Sha256")
	verified := VerifyClassmarkerWebhook(jsonData, headerHmacSignature)

	if verified {
		// Save results in your database.
		// Important: Do not use a script that will take a long time to respond.

		//Notify ClassMarker you have received the Webhook.
		fmt.Println("200")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("400")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func VerifyClassmarkerWebhook(jsonData string, headerHmacSignature string) bool {
	calculatedSignature := ComputeHmac256(jsonData, ClassmarkerWebhookSecret)
	return headerHmacSignature == calculatedSignature
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	fmt.Println("Listening on port 8080")
    http.HandleFunc("/webhook", WebHook)
	http.ListenAndServe(":8080", nil)
}
