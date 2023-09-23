package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type Instances struct {
	gorm.Model
	Id            int       `gorm:"primaryKey"`
	UserId        int       `json:"user_id"`
	ChallengeName string    `json:"challenge_name"`
	Port          int       `json:"port"`
	Timestamp     time.Time `json:"timestamp"`
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
	listChallenge := make([]string, len(t.Challenges))
	index := 1
	for name := range t.Challenges {
		listChallenge = append(listChallenge, name)
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
	db, err := gorm.Open(sqlite.Open("sql.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Instances{})
	check_open_challenge(db, listChallenge[inputChallenge-1], user.Id)

}

func check_open_challenge(db *gorm.DB, challenge_name string, user_id int) {

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
