package function

import (
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"net/http"
	"io/ioutil"
	"log"
	"crypto/tls"
	"io"
	"time"
	"encoding/json"
	"github.com/joeshaw/iso8601"
)

func getEmail(aContext *alexa.Context) string{
	respNew := doRequestOauth(http.MethodGet,aContext.System.APIEndpoint + "/v2/accounts/~current/settings/Profile.email",aContext.System.APIAccessToken,nil)
	b,_ := ioutil.ReadAll(respNew.Body)
	log.Println(string(b))
	return string(b)
}

func doRequestOauth(method, apiURL, tokenBearer string, body io.Reader ) *http.Response{
	reqNew, _ := http.NewRequest(method, apiURL, body)
	reqNew.Header.Set("Content-Type", "application/json")
	reqNew.Header.Set("Authorization", "Bearer "+tokenBearer)
	log.Println("Doing The request: ")
	log.Println(reqNew)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	respNew, _ := client.Do(reqNew)
	log.Println("Getting the next response: ")
	log.Println(respNew)
	return respNew
}

func getTimeNow() string {
	t := time.Now()

	// ISO8601 JSON format
	// "2014-03-25T16:15:25"
	data, _ := json.Marshal(iso8601.Time(t))
	return string(data)
}