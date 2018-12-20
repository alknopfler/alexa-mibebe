package function

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
)

func getUserId(aContext *alexa.Context) string{
	return aContext.System.User.UserID
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
	return time.Now().Format("2006-01-02")
}

func formatNewTime(d time.Time) string{
	return fmt.Sprintf(d.Format("2006-01-02"))

}

func getTimestamp() string {
	return time.Now().Format("20060102150405")
}

func splitFloat(f string) (string,string){
	r := strings.Split(f,".")
	return r[0],r[1]
}
