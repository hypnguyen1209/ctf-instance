package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type Challenge struct {
	Build   string `yaml:"build"`
	Timeout int    `yaml:"timeout"`
}

type T struct {
	Title      string               `yaml:"title"`
	Site_Token string               `yaml:"site_token"`
	Challenges map[string]Challenge `yaml:"challenges"`
}

type User struct {
	Email string `json:"email"`
	Score int32  `json:"score"`
	Id    int    `json:"id"`
	Name  string `json:"name"`
}

var (
	t = T{}
)

func main() {
	config, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal([]byte(config), &t)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Welcome to", t.Title)
	fmt.Println("==============================================")
	fmt.Print("Input Access_Token: ")
	var inputToken string
	fmt.Scanf("%s", &inputToken)
	resp, err := check_token(inputToken)
	if err != nil {
		log.Fatalln(err)
	}
	user := User{}
	err = json.Unmarshal([]byte(resp), &user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Hello,", user.Name, fmt.Sprintf("(score: %d)", user.Score))
	fmt.Println("List of challenges:")
	index := 1
	for name := range t.Challenges {
		fmt.Printf("	%d. %s\n", index, name)
		index++
	}
	fmt.Printf("Press the target challenge [1,..%d]: ", len(t.Challenges))
	var inputChallenge int
	fmt.Scanf("%d", &inputChallenge)
	if inputChallenge > len(t.Challenges) {
		log.Fatalln("Target invalid")
	}
	fmt.Println("Waiting for start....")
	
}

func check_token(token string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		Get(t.Site_Token + "/api/v1/users/me")
	if err != nil {
		fmt.Println("Err:", err)
		return "", err
	}
	isSuccess := gjson.Get(resp.String(), "success")
	if isSuccess.Bool() {
		dataResp := gjson.Get(resp.String(), "data")
		return dataResp.String(), nil
	}
	return "", errors.New(gjson.Get(resp.String(), "message").String())
}
