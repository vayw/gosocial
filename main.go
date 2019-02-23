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
const ServerReqURL = APIServer + "/method/groups.getLongPollServer"
const Wait = "25"
const EventsReqURL = "%s?act=a_check&key=%s&ts=%s&wait=" + Wait

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
	TS       string        `json:"ts"`
	Updates  []string      `json:"updates"`
	Failed   int           `json:"failed"`
}

type VKPollClient struct {
	APIKey, GroupID  string
	SKey, Server, TS string
}

func (vkcli *VKPollClient) GetLongPollServer() {
	var answ APIResponse
	path := "?group_id=" + vkcli.GroupID + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	url := ServerReqURL + path
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	vkcli.SKey = answ.Response.Key
	vkcli.Server = answ.Response.Server
	vkcli.TS = answ.Response.TS
}

func (vkcli VKPollClient) GetUpdates() APIResponse {
	var updates APIResponse
	url := fmt.Sprintf(EventsReqURL, vkcli.Server, vkcli.SKey, vkcli.TS)
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &updates)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return (updates)
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

	vkcli := VKPollClient{APIKey: config.API_KEY, GroupID: config.GroupID}
	vkcli.GetLongPollServer()

	for {
		fmt.Println(vkcli.GetUpdates())
	}
}
