package vault

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ask for database credential
// curl -X GET http://172.17.0.2:8200/v1/miles-mongo/creds/<role> -H "X-Vault-Token: ..."

// step to get credential
// - login to vault using user and password -> get access token
// - use recieved token to access vault secrets -> get user and password of database
// DONE

type DbCredsData struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type DbCreds struct	{
	Data DbCredsData `json:"data"`
}

func GetDatabaseCred(token string) {
	VAULT_ADDR	:= os.Getenv("VAULT_ADDR")
	VAULT_PORT	:= os.Getenv("VAULT_PORT")

	login_url := fmt.Sprintf("%s:%s/v1/miles-mongo/creds/%s", VAULT_ADDR, VAULT_PORT, "example")

	req,err := http.NewRequest( http.MethodGet, login_url, nil)

	if err!=nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", token)

	res,err := http.DefaultClient.Do(req)

	if err!=nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		return 
	}

	var dbcreds DbCreds
	json.Unmarshal(resBody,&dbcreds)

	os.Setenv("DB_USER", dbcreds.Data.Username)
	os.Setenv("DB_PASSWD", dbcreds.Data.Password)
}