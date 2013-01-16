package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Credentials struct {
	// Used to produce the request body
	Email     string `json:"email"`
	Password  string `json:"password"`
	TokenName string `json:"token_name"`
}

type SSOData struct {
	// Used to catch the values in the response body
	TokenSecret    string `json:"token_secret"`
	TokenKey       string `json:token_key`
	TokenName      string `json:"token_name"`
	ConsumerSecret string `json:"consumer_secret"`
	ConsumerKey    string `json:"consumer_key"`
}

func GetToken(email, password string) (*SSOData, error) {
	// Get a valid access token provided email and password
	msg := Credentials{email, password, "juju"}
	json_msg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	response, err := http.Post(
		"https://login.staging.ubuntu.com/api/v2/tokens",
		"application/json",
		strings.NewReader(string(json_msg)))
	if err != nil {
		return nil, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	response_body := SSOData{}
	err = json.Unmarshal(body, &response_body)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return &response_body, nil
}

type OAuth struct {
	BaseURL        string
	TokenKey       string
	TokenSecret    string
	ConsumerKey    string
	ConsumerSecret string
	TokenName      string
}

func (oauth *OAuth) Sign(req *http.Request) error {
	// Sign the provided request.
	auth := `OAuth realm="https://api.launchpad.net/", ` +
		`oauth_consumer_key="` + url.QueryEscape(oauth.ConsumerKey) + `", ` +
		`oauth_token="` + url.QueryEscape(oauth.TokenKey) + `", ` +
		`oauth_signature_method="PLAINTEXT", ` +
		`oauth_signature="` + url.QueryEscape(`&`+oauth.TokenSecret) + `", ` +
		`oauth_timestamp="` + strconv.FormatInt(time.Now().Unix(), 10) + `", ` +
		`oauth_nonce="` + strconv.Itoa(int(rand.Int31())) + `", ` +
		`oauth_version="1.0"`
	req.Header.Add("Authorization", auth)
	return nil
}

func do_request(ssodata *SSOData) {
	// FIXME remove it
	oauth := OAuth{
		"https://login.ubuntu.com/api/v2/accounts/" + ssodata.ConsumerKey,
		ssodata.TokenKey,
		ssodata.TokenSecret,
		ssodata.ConsumerKey,
		ssodata.ConsumerSecret,
		ssodata.TokenName,
	}

	request, err := http.NewRequest(
		"POST",
		"https://login.staging.ubuntu.com/api/v2/accounts/"+ssodata.TokenKey,
		nil)

	err = oauth.Sign(request)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	// run the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("response: %+v\n", response)
	fmt.Printf("Status: %s", response.Status)
	// resp_body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("response body: %s\n", resp_body)
}

func main() {
	// FIXME Just to test remove it
	ssodata, err := GetToken(
		"", "")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("Tokens: %+v\n", ssodata)
	do_request(ssodata)

}
