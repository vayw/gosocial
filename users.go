package gosocial

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	UID         int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Deactivated string `json:"deactivated"`
	Closed      bool   `json:"is_closed"`
	CanAccess   bool   `json:"can_access_closed"`
	Photo100    string `json:"photo_100"`
	Sex         int    `json:"sex"`
	About       string `json:"about"`
	Books       string `json:"books"`
	HomeTown    string `json:"home_town"`
	Interests   string `json:"interests"`
}

func (vkcli *VKClient) GetUserData(uids string, fields string) ([]User, error) {
	var answ struct {
		Response []User `json:"response"`
	}
	var params = "user_ids=%s&fields=%s"
	method := "/method/users.get?"
	if fields == "" {
		fields = "photo_100"
	}
	params = fmt.Sprintf(params, uids, fields)
	url := APIServer + method + params + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	res, err := http.Get(url)
	if err != nil {
		logger.Print("[ERR} ", err)
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		logger.Print("[ERR]", jsonErr)
		return []User{}, errors.New("Failed to get userdata")
	}
	return answ.Response, nil
}
