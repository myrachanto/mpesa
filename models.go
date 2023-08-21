package mpesa

import "net/http"

// Mpesa is an application that will be making a transaction
type Mpesa struct {
	consumerKey    string
	consumerSecret string
	baseURL        string
	client         *http.Client
}

// MpesaOpts stores all the configuration keys we need to set up a Mpesa app,
type MpesaOpts struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string
}

// MpesaAccessTokenResponse is the response sent back by Safaricom when we make a request to generate a token
type MpesaAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RequestID    string `json:"requestId"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// STKPushRequestBody is the body with the parameters to be used to initiate an STK push request
type STKPushRequestBody struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// STKPushRequestResponse is the response sent back after initiating an STK push request.
type STKPushRequestResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
	RequestID           string `json:"requestId"`
	ErrorCode           string `json:"errorCode"`
	ErrorMessage        string `json:"errorMessage"`
}
