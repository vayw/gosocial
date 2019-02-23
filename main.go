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

var logger = log.New(os.Stdout, "gosocial ", log.Ltime)

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
	Updates  []UpdateEvent `json:"updates"`
	Failed   int           `json:"failed"`
}

type VKPollClient struct {
	APIKey, GroupID  string
	SKey, Server, TS string
}

type UpdateEvent struct {
	Type     string     `json:"type"`
	GroupID  string     `json:"group"`
	EventObj GroupEvent `json:"object"`
}

type GroupEvent struct {
	UID       int    `json:"user_id"`
	JoinType  string `json:"join_type,omitempty"`
	LeaveType int    `json:"self,omitempty"`
}

func (vkcli *VKPollClient) GetLongPollServer() {
	var answ APIResponse
	path := "?group_id=" + vkcli.GroupID + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	url := ServerReqURL + path
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		logger.Print("[ERR} ", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		logger.Print("[ERR]", jsonErr)
	}
	vkcli.SKey = answ.Response.Key
	vkcli.Server = answ.Response.Server
	vkcli.TS = answ.Response.TS
}

func (vkcli *VKPollClient) GetUpdates() ([]UpdateEvent, int) {
	var api_resp APIResponse
	url := fmt.Sprintf(EventsReqURL, vkcli.Server, vkcli.SKey, vkcli.TS)
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		logger.Print("[ERR]", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &api_resp)
	if jsonErr != nil {
		logger.Print("[ERR]", jsonErr)
	}

	switch api_resp.Failed {
	case 1:
		// history is outdated or partly lost, try again with TS
		// from current answer
		vkcli.TS = api_resp.TS
		return nil, 1
	case 2:
		// session key is expired
		vkcli.GetLongPollServer()
		return nil, 2
	case 3:
		// history is lost, request new key and ts
		vkcli.GetLongPollServer()
		return nil, 3
	}
	vkcli.TS = api_resp.TS
	return api_resp.Updates, 0
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
