package gosocial

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const QueryCount = 1000

type Members struct {
	Count int   `json:"count"`
	Items []int `json:"items"`
}

func (vkcli *VKClient) GetMembers() (Members, error) {
	var sum Members

	count, err := vkcli.MembersCount()
	if err != nil {
		return Members{}, err
	}
	sum.Count = count

	for len(sum.Items) < sum.Count {
		resp, err := vkcli.QueryMembers(QueryCount, len(sum.Items))
		if err != nil {
			logger.Print("[ERR]", err)
			return Members{}, errors.New("Failed to get group members")
		}
		sum.Items = append(sum.Items, resp.Items...)
	}

	return sum, nil
}

func (vkcli *VKClient) MembersCount() (int, error) {
	resp, err := vkcli.QueryMembers(0, 0)
	if err != nil {
		logger.Print("::MemebersCount::", err)
		return -1, err
	}
	return resp.Count, nil
}

func (vkcli *VKClient) QueryMembers(count int, offset int) (Members, error) {
	var answ struct {
		Response Members `json:"response"`
	}
	method := "/method/groups.getMembers?"
	params := `group_id=%s&count=%d&offset=%d`
	url := APIServer + method + fmt.Sprintf(params, vkcli.GroupID, 300, 0) + "&access_token=" + vkcli.APIKey + "&v=" + APIv
	res, err := http.Get(url)
	if err != nil {
		logger.Print("[ERR} ", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	jsonErr := json.Unmarshal([]byte(body), &answ)
	if jsonErr != nil {
		logger.Print("[ERR]", jsonErr)
		return Members{}, errors.New("Failed to get group members")
	}
	return answ.Response, nil
}
