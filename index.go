package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// NewMpesa sets up and returns an instance of Mpesa
func newMpesa(m *MpesaOpts) *Mpesa {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Mpesa{
		consumerKey:    m.ConsumerKey,
		consumerSecret: m.ConsumerSecret,
		baseURL:        m.BaseURL,
		client:         client,
	}
}

// takes an amount float32 and phoneNumber 254700000000 and returns *STKPushRequestResponse, error
// env variables include appKey, appSecret,baseUrl,shortcode,passkey,partyA,callback,AccountReference,TransactionDesc
// partyA is your phone number 254700000000
// baseUrl ttps://sandbox.safaricom.co.ke
func Process(amount float32, PhoneNumber string) (*STKPushRequestResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("to load env files")
	}
	appKey := os.Getenv("appKey")
	appSecret := os.Getenv("appSecret")
	baseUrl := os.Getenv("baseUrl")
	shortcode := os.Getenv("shortcode")
	passkey := os.Getenv("passkey")
	partyA := os.Getenv("partyA")
	callback := os.Getenv("callback")
	AccountReference := os.Getenv("AccountReference")
	TransactionDesc := os.Getenv("TransactionDesc")
	mpesa := newMpesa(&MpesaOpts{
		ConsumerKey:    appKey,
		ConsumerSecret: appSecret,
		BaseURL:        baseUrl,
	})
	timestamp := time.Now().Format("20060102150405")
	passwordToEncode := fmt.Sprintf("%s%s%s", shortcode, passkey, timestamp)
	password := base64.StdEncoding.EncodeToString([]byte(passwordToEncode))

	accessTokenResponse, err := mpesa.generateAccessToken()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", accessTokenResponse)
	response, err := mpesa.initiateSTKPushRequest(&STKPushRequestBody{
		BusinessShortCode: shortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            fmt.Sprintf("%.2f", amount), // Amount to be charged when checking out
		PartyA:            partyA,                      // 2547XXXXXXXX
		PartyB:            shortcode,
		PhoneNumber:       PhoneNumber, // 2547XXXXXXXX
		CallBackURL:       callback,    // https://
		AccountReference:  AccountReference,
		TransactionDesc:   TransactionDesc,
	})

	if err != nil {
		return nil, err
	}
	return response, nil

}
func (m *Mpesa) generateAccessToken() (*MpesaAccessTokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", m.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(m.consumerKey, m.consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.makeRequest(req)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := &MpesaAccessTokenResponse{}

	if err := json.Unmarshal(resp, &accessTokenResponse); err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}

func (m *Mpesa) initiateSTKPushRequest(body *STKPushRequestBody) (*STKPushRequestResponse, error) {
	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", m.baseURL)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	accessTokenResponse, err := m.generateAccessToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessTokenResponse.AccessToken))

	resp, err := m.makeRequest(req)
	if err != nil {
		return nil, err
	}

	stkPushResponse := new(STKPushRequestResponse)
	if err := json.Unmarshal(resp, &stkPushResponse); err != nil {
		return nil, err
	}

	return stkPushResponse, nil
}

// makeRequest performs all the http requests for the specific app
func (m *Mpesa) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
