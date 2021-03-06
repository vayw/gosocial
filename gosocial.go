package gosocial

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const APIv = "5.95"
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

type APIResponseF struct {
	TS     int64 `json:"ts"`
	Failed int   `json:"failed"`
}

type VKClient struct {
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

func (vkcli *VKClient) GetLongPollServer() error {
	var answ APIResponse
	path := "?group_id=" + vkcli.GroupID + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	url := ServerReqURL + path
	logger.Print("::GetLongPollServer:: Query: ", url)
	res, err := http.Get(url)
	if err != nil {
		logger.Print("::GetLongPollServer:: [ERR} ", err)
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		logger.Print("::GetLongPollServer:: [ERR]", jsonErr)
		return errors.New("Failed to parse answer")
	}
	logger.Print("::GetLongPollServer:: Response: ", answ.Response)
	vkcli.SKey = answ.Response.Key
	vkcli.Server = answ.Response.Server
	vkcli.TS = answ.Response.TS
	return nil
}

func (vkcli *VKClient) GetUpdates() ([]UpdateEvent, error) {
	var api_resp APIResponse
	url := fmt.Sprintf(EventsReqURL, vkcli.Server, vkcli.SKey, vkcli.TS)
	logger.Print("::GetUpdates:: query= ", url)
	res, err := http.Get(url)
	if err != nil {
		logger.Print("[ERR] ::GetUpdates:: ", err)
		return nil, fmt.Errorf("VKClient GetUpdates: %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	logger.Print("::GetUpdates:: response string: ", string(body))
	jsonErr := json.Unmarshal([]byte(body), &api_resp)
	if jsonErr != nil {
		logger.Print("[INFO] ::GetUpdates:: ", jsonErr)
		api_resp.Failed = 1
	}
	if api_resp.Failed != 0 {
		var api_resp_f APIResponseF
		jsonErr := json.Unmarshal([]byte(body), &api_resp_f)
		if jsonErr != nil {
			logger.Print("[ERR] ::GetUpdates:: ", jsonErr)
			return nil, fmt.Errorf("VKClient GetUpdates: %v", jsonErr)
		}
		switch api_resp_f.Failed {
		case 1:
			// history is outdated or partly lost, try again with TS
			// from current answer
			vkcli.TS = strconv.FormatInt(api_resp_f.TS, 10)
			return nil, errors.New("VKClient GetUpdates: history is outdated or partly lost")
		case 2:
			// session key is expired
			return nil, errors.New("VKClient GetUpdates: session key has expired")
		case 3:
			// history is lost, request new key and ts
			return nil, errors.New("VKClient GetUpdates: history is lost")
		}
	}

	logger.Print("::GetUpdates:: updates=", api_resp)
	vkcli.TS = api_resp.TS
	return api_resp.Updates, nil
}
