package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// using userpass
// login
// curl -X POST -d @pwd.json http://172.17.0.2:8200/v1/auth/mongouser/login/mongo
// res token is in res.auth.client_token

// ask for database credential
// curl -X GET http://172.17.0.2:8200/v1/miles-mongo/creds/<role> -H "X-Vault-Token: ..."

// step to get credential
// - login to vault using user and password -> get access token
// - use recieved token to access vault secrets -> get user and password of database
// DONE

type userpassPayload struct {
	Password string `json:"password"`
}

type userpassAuth struct {
	Token string `json:"client_token"`
}

type Userpass struct {
	Auth userpassAuth `json:"auth"`
}

func Login() {
	VAULT_ADDR	:= os.Getenv("VAULT_ADDR")
	VAULT_PORT	:= os.Getenv("VAULT_PORT")
	VAULT_USER	:= os.Getenv("VAULT_USER") 
	payload	:= userpassPayload{ Password: os.Getenv("VAULT_PWD") }

	login_url := fmt.Sprintf("%s:%s/v1/auth/mongouser/login/%s", VAULT_ADDR, VAULT_PORT, VAULT_USER)

	body, err := json.Marshal(payload)
	if err!=nil {
		return
	}

	req,err := http.NewRequest( http.MethodPost, login_url, bytes.NewReader(body) )

	if err!=nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res,err := http.DefaultClient.Do(req)

	if err!=nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		return 
	}

	var userpass Userpass
	json.Unmarshal(resBody,&userpass)

	// this is tokena
	os.Setenv("VAULT_TOKEN", userpass.Auth.Token)
}