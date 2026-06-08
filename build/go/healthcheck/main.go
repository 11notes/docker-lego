package main

import (
	"encoding/json"
  "github.com/11notes/go-eleven"
)

type Account struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	KeyType      string `json:"keyType"`
	Server       string `json:"server"`
	Origin       string `json:"origin"`
	Registration struct {
		Status     string `json:"status"`
		AccountURL string `json:"accountURL"`
	} `json:"registration"`
}

func main(){
	file, err := eleven.Util.ReadFile("/lego/var/accounts/acme-v02.api.letsencrypt.org/default/account.json")
	if err != nil {
		eleven.LogFatal(err.Error())
	}
	var account Account
	err = json.Unmarshal([]byte(file), &account)
	if err != nil {
		eleven.LogFatal("not valid ACME account json!")
	}
	if account.Registration.Status != "valid" {
		eleven.LogFatal("account status is %s", account.Registration.Status)
	}
}