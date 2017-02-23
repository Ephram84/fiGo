package fiGo

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
)

var (
	// ErrCantFindUser means "Cant find user in fake-figo"
	ErrCantFindUser = errors.New("cant find user")
)

// FakeConnection is a fakes the figoAPI
type FakeConnection struct {
	Users []map[string]interface{}
}

// NewFakeConnection creates a new "fak-Connection" only in memory
func NewFakeConnection() *FakeConnection {
	return &FakeConnection{Users: make([]map[string]interface{}, 0)}
}

// CreateUser "store" a user in this fake-Connection
func (fakeConnection *FakeConnection) CreateUser(name string, email string, password string) ([]byte, error) {
	recoveryPassword := randStringRunes(10)
	// "store" in fakeConnection
	fakeConnection.Users = append(fakeConnection.Users, map[string]interface{}{
		"name":              name,
		"email":             email,
		"password":          password,
		"recovery_password": recoveryPassword,
		"access_token":      randStringRunes(10),
		"expires_in":        3600,
		"refresh_token":     randStringRunes(10),
		"scope":             "accounts=ro balance=ro transactions=ro offline",
		"token_type":        "Bearer"})

	responseMap := map[string]string{"recovery_password": recoveryPassword}
	return json.Marshal(responseMap)
}

// CredentialLogin returns a token
// -> Notice: first add a user via CreateUser!
func (fakeConnection *FakeConnection) CredentialLogin(username string, password string) ([]byte, error) {
	// search user
	for _, user := range fakeConnection.Users {
		log.Println(user)
		if user["email"] == username && user["password"] == password {
			response := map[string]interface{}{
				"access_token":  user["access_token"],
				"expires_in":    user["expires_in"],
				"refresh_token": user["refresh_token"],
				"scope":         user["scope"],
				"token_type":    user["token_type"],
			}
			return json.Marshal(response)
		}
	}
	return nil, ErrCantFindUser
}

// SetupNewBankAccount sets up a new bank account for an existing account.
// -> Notice: you need an accessToken from CredentialLogin
func (fakeConnection *FakeConnection) SetupNewBankAccount(accessToken string, bankCode string, country string, credentials []string) ([]byte, error) {
	for _, user := range fakeConnection.Users {
		if user["access_token"] == accessToken {
			user["banks"] = map[string]interface{}{
				"bankCode":    bankCode,
				"country":     country,
				"credentials": credentials,
			}
			response := map[string]string{"task_token": randStringRunes(10)}
			return json.Marshal(response)
		}
	}
	return nil, ErrCantFindUser
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}