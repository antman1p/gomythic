package Mythic_Go_Scripting

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"crypto/tls"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type MythicUtilities struct {
	Mythic *Mythic
}

func NewMythicUtilities(mythic *Mythic) *MythicUtilities {
	return &MythicUtilities{
		Mythic: mythic,
	}
}

func (u *MythicUtilities) GetHeaders() http.Header {
	headers := make(http.Header)
	if u.Mythic.APIToken != "" {
		headers.Set("apitoken", u.Mythic.APIToken)
	} else if u.Mythic.AccessToken != "" {
		headers.Set("Authorization", "Bearer "+u.Mythic.AccessToken)
	}
	return headers
}

func (u *MythicUtilities) GetHTTPTransport() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (u *MythicUtilities) GetWSTransport() *websocket.Dialer {
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return dialer
}

// HTTPPost performs a POST request to the specified URL and returns the response
func (u *MythicUtilities) HTTPPost(data map[string]interface{}, url string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: u.GetHTTPTransport(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header = u.GetHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}


func (u *MythicUtilities) HTTPPostForm(data url.Values, url string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: u.GetHTTPTransport(),
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = u.GetHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *MythicUtilities) HTTPGetDictionary(url string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: u.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = u.GetHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *MythicUtilities) HTTPGet(url string) ([]byte, error) {
	client := &http.Client{
		Transport: u.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = u.GetHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func (u *MythicUtilities) HTTPGetChunked(url string, chunkSize int) (<-chan []byte, error) {
	client := &http.Client{
		Transport: u.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = u.GetHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)

	go func() {
		defer close(ch)
		buf := make([]byte, chunkSize)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				ch <- buf[:n]
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading HTTP response:", err)
				break
			}
		}
		resp.Body.Close()
	}()

	return ch, nil
}

func (u *MythicUtilities) GraphQLPost(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	response, err := u.HTTPPost(data, u.Mythic.HTTP+u.Mythic.ServerIP+":"+strconv.Itoa(u.Mythic.ServerPort)+"/graphql")
	if err != nil {
		return nil, err
	}

	if responseErrors, ok := response["errors"]; ok {
		errMsg := ""
		for _, err := range responseErrors.([]interface{}) {
			errMsg += fmt.Sprintf("%s\n", err.(map[string]interface{})["message"])
		}
		return nil, errors.New(errMsg)
	}

	return response, nil
}



func (u *MythicUtilities) GraphQLSubscription(query string, variables map[string]interface{}, timeout int) (<-chan map[string]interface{}, error) {
	data := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	response, err := u.HTTPPost(data, u.Mythic.HTTP+u.Mythic.ServerIP+":"+strconv.Itoa(u.Mythic.ServerPort)+"/subscriptions")
	if err != nil {
		return nil, err
	}

	if responseErrors, ok := response["errors"]; ok {
		errMsg := ""
		for _, err := range responseErrors.([]interface{}) {
			errMsg += fmt.Sprintf("%s\n", err.(map[string]interface{})["message"])
		}
		return nil, errors.New(errMsg)
	}

	events := make(chan map[string]interface{})
	go func() {
		defer close(events)
		for {
			select {
			case <-time.After(time.Duration(timeout) * time.Second):
				return
			default:
				response, err := u.HTTPGet(u.Mythic.HTTP + u.Mythic.ServerIP + ":" + strconv.Itoa(u.Mythic.ServerPort) + "/graphql/events")
				if err != nil {
					log.Println("Error receiving GraphQL subscription event:", err)
					return
				}
				var event map[string]interface{}
				err = json.Unmarshal(response, &event)
				if err != nil {
					log.Println("Error parsing GraphQL subscription event:", err)
					return
				}
				events <- event
			}
		}
	}()

	return events, nil
}

func (u *MythicUtilities) FetchGraphQLSchema() (string, error) {
	response, err := u.HTTPGet(u.Mythic.HTTP + u.Mythic.ServerIP + ":" + strconv.Itoa(u.Mythic.ServerPort) + "/graphql/schema.json")
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (u *MythicUtilities) LoadMythicSchema() bool {
	schema, err := u.FetchGraphQLSchema()
	if err != nil {
		log.Println("Failed to fetch Mythic schema:", err)
		return false
	}

	u.Mythic.Schema = schema
	return true
}
