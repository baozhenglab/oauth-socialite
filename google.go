package oauthsocialite

import (
	"errors"
	"flag"
	"fmt"
)

type googleSocialite struct {
	clientID     string
	clientSecret string
	name         string
}

const (
	graphGGUrl = "https://www.googleapis.com/oauth2"
	versionGG  = "v2"
)

//UserFromToken return user User,status int,err
func (gg *googleSocialite) UserFromToken(token string) (user User, status int, panErr error) {
	meURL := graphGGUrl + "/" + versionGG + "/userinfo"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
	var errRes error
	var response map[string]interface{}
	response, errRes = RequestGet(meURL, headers)
	if errRes != nil {
		return User{}, 500, errRes
	}
	return converResponse(response, token)
}
func mapUserGoogleFromObject(user map[string]interface{}) User {
	return User{
		ID:     fmt.Sprintf("%s", user["id"]),
		Name:   fmt.Sprintf("%s", user["name"]),
		Email:  fmt.Sprintf("%s", user["email"]),
		Token:  fmt.Sprintf("%s", user["token"]),
		Avatar: fmt.Sprintf("%s", user["picture"]),
	}
}

func converResponse(response map[string]interface{}, token string) (user User, status int, panErr error) {
	if _, ok := response["email"]; !ok {
		return User{}, 400, errors.New("Token required email scope")
	}
	if response["status"] != 200 {
		rd := response["error"].(map[string]interface{})
		return User{}, response["status"].(int), errors.New(fmt.Sprintf("%s", rd["message"]))
	}
	response["token"] = token

	return mapUserGoogleFromObject(response), 200, nil
}

func (gg *googleSocialite) Get() interface{} {
	return gg
}

func (gg *googleSocialite) Name() string {
	return gg.name
}
func (gg *googleSocialite) GetPrefix() string {
	return gg.name
}

func (gg *googleSocialite) InitFlags() {
	prefix := fmt.Sprintf("%s-", gg.Name())
	flag.StringVar(&gg.clientID, prefix+"client-id", "", "oauth client id google")
	flag.StringVar(&gg.clientSecret, prefix+"client-secrect", "", "oauth client secrect google")
}
