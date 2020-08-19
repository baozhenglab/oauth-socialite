package oauthsocialite

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type facebookSocialite struct {
	clientID     string
	clientSecret string
	name         string
}

const (
	graphUrl = "https://graph.facebook.com"
	version  = "v4.0"
)

//UserFromToken return user User,status int,err
func (fb *facebookSocialite) UserFromToken(token string) (user User, status int, error error) {
	var fields = []string{"name", "email", "gender", "verified", "link"}
	meURL := graphUrl + "/" + version + "/me?access_token=" + token + "&fields=" + strings.Join(fields, ",")
	if fb.clientSecret != "" {
		appSecretProof := GenerateHashMac(fb.clientSecret, token)
		meURL += "&appsecret_proof=" + appSecretProof
	}
	response, err := RequestGet(meURL, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return User{}, 500, err
	}
	if response["status"] != 200 {
		rd := response["error"].(map[string]interface{})
		return User{}, response["status"].(int), errors.New(fmt.Sprintf("%s", rd["message"]))
	}
	if _, ok := response["email"]; !ok {
		return User{}, 400, errors.New("Token required email scope")
	}
	response["token"] = token

	return mapUserFromObject(response), 200, nil
}
func mapUserFromObject(user map[string]interface{}) User {
	avatarURL := graphUrl + "/" + version + "/" + fmt.Sprintf("%s", user["id"]) + "/picture"
	return User{
		ID:     fmt.Sprintf("%s", user["id"]),
		Name:   fmt.Sprintf("%s", user["name"]),
		Email:  fmt.Sprintf("%s", user["email"]),
		Token:  fmt.Sprintf("%s", user["token"]),
		Avatar: avatarURL + "?type=normal",
	}
}

func (fb *facebookSocialite) Name() string {
	return fb.name
}

func (fb *facebookSocialite) GetPrefix() string {
	return fb.name
}

func (fb *facebookSocialite) InitFlags() {
	prefix := fmt.Sprintf("%s-", fb.Name())
	flag.StringVar(&fb.clientID, prefix+"client-id", "", "oauth client id facebook")
	flag.StringVar(&fb.clientSecret, prefix+"client-secrect", "", "oauth client secrect facebook")
}

func (fb *facebookSocialite) Get() interface{} {
	return fb
}
