package fiGo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
)

const (
	baseURL     = "https://api.figo.me"
	authUserURL = "/auth/user"
)

var (
	// ErrUserAlreadyExists - code 30002
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrHTTPUnauthorized - code 90000
	ErrHTTPUnauthorized = errors.New("invalid authorization")
)

// Connection represent a connection to figo
type Connection struct {
	AuthString string
}

// NewFigoConnection creates a new connection.
// -> You need a clientID and a clientSecret. (You will get this from figo.io)
func NewFigoConnection(clientID string, clientSecret string) *Connection {
	authInfo := clientID + ":" + clientSecret
	authString := "Basic " + base64.URLEncoding.EncodeToString([]byte(authInfo))
	return &Connection{AuthString: authString}
}

// CreateUser creates a new user.
// Ask User for name, email and a password
func (connection *Connection) CreateUser(name string, email string, password string) ([]byte, error) {
	// build url
	url := baseURL + authUserURL

	// build jsonBody
	requestBody := map[string]string{
		"name":     name,
		"email":    email,
		"password": password}
	jsonBody, err := json.Marshal(requestBody)

	// build request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", connection.AuthString)

	// setup client
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// get response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check response for errors
	jsonParsed, err := gabs.ParseJSON(body)
	value, ok := jsonParsed.Path("error.code").Data().(float64)
	if ok {
		if value == 30002 {
			return body, ErrUserAlreadyExists
		} else if value == 90000 {
			return body, ErrHTTPUnauthorized
		}
	}

	return body, nil
}