package gosocial

import (
	"encoding/json"
	"errors"
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
}

func (vkcli *VKClient) GetUserData(uids string) ([]User, error) {
	var answ struct {
		Response []User `json:"response"`
	}
	var params string
	method := "/method/users.get?"
	params = "user_ids=" + uids + "&fields=photo_100,sex"
	url := APIServer + method + params + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		logger.Print("[ERR} ", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		logger.Print("[ERR]", jsonErr)
		return []User{}, errors.New("Failed to get userdata")
	}
	return answ.Response, nil
}
