package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

func main() {
	config, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	t := T{}
	err = yaml.Unmarshal([]byte(config), &t)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Welcome to ")
	fmt.Print("Input Access_Token: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n') // convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	fmt.Println(text)
}

func check_token(token string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("ctfd_92233962e4dc146288edf8602fa5d42790b1f172f275ed63fb253289f44fd636").
		Get("http://ctf.actvn.edu.vn/api/v1/users/me")
	if err != nil {
		fmt.Println("Err:", err)
		return "", err
	}
	isSuccess := gjson.Get(resp.String(), "success")
	fmt.Println(isSuccess.String())
	return "", nil
}
