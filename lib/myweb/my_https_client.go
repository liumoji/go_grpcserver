package myweb

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"package": "myweb-my_https_client"})

func testStart() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8081")

	if err != nil {
		log.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func NewHttpsClient(host string, tlsFlag bool, caCertPath string) *http.Client {
	var tr *http.Transport
	if tlsFlag {
		pool := x509.NewCertPool()
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Println("NewHttpsClient ReadFile err:", err)
			return nil
		}
		pool.AppendCertsFromPEM(caCrt)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	return client
}

func DecodeRequest(resp *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}
	return dat, nil
}

func postRequest(client *http.Client) {
	var r http.Request
	r.ParseForm()
	r.Form.Add("uuid", "qwetyu")
	bodystr := strings.TrimSpace(r.Form.Encode())
	req, err := http.NewRequest("POST", "https://10.21.140.160:10000/v1/gcp", strings.NewReader(bodystr))
	if err != nil {
		fmt.Println("NewRequest err", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "Keep-Alive")
	resp, err := client.Do(req)
	//resp, err := client.("https://10.21.140.160:10000/v1/gcp")

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
}
