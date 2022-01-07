package utils

import (
	"encoding/base64"
	"testing"
)

func TestKeyUtil(t *testing.T) {
	appKey := "NUhONBRTxPxFtFkH78P9AJ2EDUJ1EeoaFzGVoJUz5BcYFtqiag0baRw61y1ycoZaYpkxp9BC08K2F8h2II4tyQ=="
	appKeyBytes, _ := base64.StdEncoding.DecodeString(appKey)
	apikey := GetMsgEncAesIV(appKeyBytes)
	println(base64.StdEncoding.EncodeToString(apikey))

}
