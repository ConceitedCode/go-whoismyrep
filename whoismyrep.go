package whoismyrep

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_ENDPOINT = "whoismyrepresentative.com"
)

type WIMR struct {
	httpClient *http.Client
}

func Open() *WIMR {
	return &WIMR{httpClient: &http.Client{}}
}

// make an api request
func (wimr *WIMR) api(method string, path string, fields url.Values) (body []byte, err error) {
	var req *http.Request
	url := fmt.Sprintf("http://%s%s", API_ENDPOINT, path)

	if method == "POST" && fields != nil {
		req, err = http.NewRequest(method, url, strings.NewReader(fields.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		if fields != nil {
			url += "?" + fields.Encode()
		}
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}
	rsp, err := wimr.httpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		err = fmt.Errorf("whoismyrep error: %d %s", rsp.StatusCode, body)
	}
	return
}

type Rep struct {
	Name     string `json:"name"`
	Party    string `json:"party"`
	State    string `json:"state"`
	District string `json:"district"`
	Phone    string `json:"phone"`
	Office   string `json:"office"`
	Link     string `json:"link"`
}

func (wimr *WIMR) RepsByZip(zip string, zip4 string) (res []Rep, err error) {
	v := url.Values{}
	v.Set("zip", zip)
	v.Set("zip4", zip4)
	v.Set("output", "json")

	body, err := wimr.api("GET", "/getall_mems.php", v)
	if err != nil {
		return
	}

	var j struct {
		Items []Rep `json:"results"`
	}

	err = json.Unmarshal(body, &j)
	if err != nil {
		return
	}
	res = j.Items
	return

}
