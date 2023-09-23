package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

type InstancesTable struct {
	gorm.Model
	Id            int    `gorm:"primaryKey;autoIncrement"`
	UserId        int    `json:"user_id"`
	ChallengeName string `json:"challenge_name"`
	Port          int    `json:"port"`
	//Timestamp     time.Time `json:"timestamp"`
}

var (
	t = T{}
)

func checkOpenChallenge(db *gorm.DB, challenge_name string, user_id int) (*InstancesTable, bool) {
	var instanceOpen InstancesTable
	result := db.Where("challenge_name=?", challenge_name).Where("user_id=?", user_id).First(&instanceOpen)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		} else {
			log.Fatalln(result.Error)
		}
	}
	return &instanceOpen, true
}

func checkToken(token string) (string, error) {
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

func getIp() (string, error) {
	client := resty.New()
	resp, err := client.R().
		Get("http://ip-api.com/json/")
	if err != nil {
		fmt.Println("Err:", err)
		return "", err
	}
	return gjson.Get(resp.String(), "query").String(), nil
}

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
	resp, err := checkToken(inputToken)
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
	listChallenge := []string{}
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
	db, err := gorm.Open(sqlite.Open("sql.db"), &gorm.Config{
		QueryFields: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&InstancesTable{})
	challengeName := listChallenge[inputChallenge-1]
	instanceOpen, status := checkOpenChallenge(db, challengeName, user.Id)
	ip, _ := getIp()
	if !status {
		///
		prefixPort := 20 + len(challengeName)
		port, _ := strconv.Atoi(fmt.Sprintf("%d%03d", prefixPort, user.Id))
		result := db.Create(&InstancesTable{
			ChallengeName: challengeName,
			UserId:        user.Id,
			Port:          port,
			//Timestamp:     time.Now(),
		})

		if result.RowsAffected > 0 {
			timeout := t.Challenges[challengeName].Timeout
			cmd := exec.Command("./instance.sh",
				t.Challenges[challengeName].Build,
				strconv.Itoa(user.Id),
				strconv.Itoa(port),
				strconv.Itoa(timeout),
				challengeName)
			cmd.Start()
			time.Sleep(5 * time.Second)
			fmt.Printf("Challenge: http://%s:%d\n", ip, port)
			fmt.Println("Time to close:", timeout, "minutes")
			time.Sleep(time.Duration(timeout) * time.Minute)
			fmt.Println("Ended")
			os.Exit(0)
		}
	}
	if time.Now().Unix()-instanceOpen.CreatedAt.Unix() < int64(t.Challenges[challengeName].Timeout*60) {
		fmt.Printf("Challenge: http://%s:%d\n", ip, instanceOpen.Port)
		timeout := float64(int64(t.Challenges[challengeName].Timeout*60)-time.Now().Unix()+instanceOpen.CreatedAt.Unix()) / 60
		fmt.Println("Time to close:", fmt.Sprintf("%.2f", timeout), "minutes")
		time.Sleep(time.Duration(timeout) * time.Minute)
		fmt.Println("Ended")
		os.Exit(0)
	} else {
		result := db.Where("id=?", instanceOpen.Id).Delete(&InstancesTable{})
		if result.RowsAffected > 0 {
			prefixPort := 20 + len(challengeName)
			port, _ := strconv.Atoi(fmt.Sprintf("%d%03d", prefixPort, user.Id))
			resultCreate := db.Create(&InstancesTable{
				ChallengeName: challengeName,
				UserId:        user.Id,
				Port:          port,
				//Timestamp:     time.Now(),
			})
			if resultCreate.RowsAffected > 0 {
				timeout := t.Challenges[challengeName].Timeout
				cmd := exec.Command("./instance.sh",
					t.Challenges[challengeName].Build,
					strconv.Itoa(user.Id),
					strconv.Itoa(port),
					strconv.Itoa(timeout),
					challengeName)
				cmd.Start()
				time.Sleep(5 * time.Second)
				fmt.Printf("Challenge: http://%s:%d\n", ip, port)
				fmt.Println("Time to close:", timeout, "minutes")
				time.Sleep(time.Duration(timeout) * time.Minute)
				fmt.Println("Ended")
				os.Exit(0)
			}
		}
	}
}
