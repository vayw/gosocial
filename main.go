package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const APIv = "5.92"
const APIServer = "https://api.vk.com"
const BotAPIUrl = APIServer + "/method/groups.getLongPollServer"

type Configuration struct {
	API_KEY, GroupID string
}

type LongPollParam struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     string `json:"ts"`
}

type APIResponse struct {
	Response LongPollParam `json:"response"`
}

func GetLongPollServer(api_key string, group_id string) (LongPollParam, error) {
	var params APIResponse
	path := "?group_id=" + group_id + "&access_token=" + api_key + "&v=" + APIv
	url := BotAPIUrl + path
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &params)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return params.Response, err
}

func test(resp string) {
	w := APIResponse{}
	e := []byte(resp)
	jsonErr := json.Unmarshal(e, &w)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Println(w)
}

func main() {
	var (
		config Configuration
	)
	file, _ := os.Open("conf.json")
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(byteValue, &config)

	if err != nil {
		fmt.Println("error loading configuration:", err)
	}

	params, _ := GetLongPollServer(config.API_KEY, config.GroupID)
	fmt.Println(params.Key)
}
