package mythic_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Mythic struct {
	MythicClasses
	MythicUtilities *MythicUtilities
}

func NewMythic(username, password, serverIP string, serverPort int, apiToken string, ssl bool, timeout int) *Mythic {
	mythicClasses := MythicClasses{
		Username:           username,
		Password:           password,
		APIToken:           apiToken,
		AccessToken:        "",
		RefreshToken:       "",
		ServerIP:           serverIP,
		ServerPort:         serverPort,
		SSL:                ssl,
		HTTP:               "http://",
		WS:                 "ws://",
		GlobalTimeout:      timeout,
		ScriptingVersion:   "0.1.2",
		CurrentOperationID: 0,
		Schema:             "",
	}

	mythicUtilities := NewMythicUtilities(&mythicClasses)

	return &Mythic{
		MythicClasses:    mythicClasses,
		MythicUtilities: mythicUtilities,
	}
}

type GraphqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
	} `json:"errors"`
}

func (m *Mythic) GraphqlPost(query string) (GraphqlResponse, error) {
	url := fmt.Sprintf("%s%s:%d/graphql", m.HTTP, m.ServerIP, m.ServerPort)
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+m.AccessToken)
	headers.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(query))
	if err != nil {
		return GraphqlResponse{}, err
	}
	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return GraphqlResponse{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GraphqlResponse{}, err
	}
	var gqlResp GraphqlResponse
	err = json.Unmarshal(body, &gqlResp)
	if err != nil {
		return GraphqlResponse{}, err
	}
	return gqlResp, nil
}

func (m *Mythic) HttpPost(url string, data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dataResp map[string]interface{}
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, err
	}
	return dataResp, nil
}

func (m *Mythic) HttpPostForm(data url.Values, url string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = make(http.Header)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (m *Mythic) HttpGetDictionary(url string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (m *Mythic) HttpGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (m *Mythic) HttpGetChunked(url string, chunkSize int) (<-chan []byte, error) {
	client := &http.Client{Timeout: time.Duration(m.GlobalTimeout) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
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

func (m *Mythic) GraphQLPost(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	resp, err := m.MythicUtilities.GraphQLPost(query, variables)
	if err != nil {
		return nil, err
	}
	return resp, nil
}



func (m *Mythic) GraphQLSubscription(query string, variables map[string]interface{}, timeout int) (<-chan map[string]interface{}, error) {
	return m.MythicUtilities.GraphQLSubscription(query, variables, timeout)
}

func (m *Mythic) FetchGraphQLSchema() (string, error) {
	return m.MythicUtilities.FetchGraphQLSchema()
}

func (m *Mythic) LoadMythicSchema() bool {
	return m.MythicUtilities.LoadMythicSchema()
}

func Login(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout, loggingLevel int) (*Mythic, error) {
	mythic := NewMythic(username, password, serverIP, serverPort, apiToken, ssl, timeout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if apiToken == "" {
		url := fmt.Sprintf("%s%s:%d/auth", mythic.HTTP, mythic.ServerIP, mythic.ServerPort)
		data := map[string]interface{}{
			"username":          mythic.Username,
			"password":          mythic.Password,
			"scripting_version": mythic.ScriptingVersion,
		}
		log.Printf("[*] Logging into Mythic as scripting_version %s", mythic.ScriptingVersion)
		response, err := mythic.HttpPost(url, data)
		if err != nil {
			log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
			return nil, err
		}
		mythic.AccessToken = response["access_token"].(string)
		mythic.RefreshToken = response["refresh_token"].(string)
		user := response["user"].(map[string]interface{})
		mythic.CurrentOperationID = int(user["current_operation_id"].(float64))
		currentTokens, err := mythic.GraphqlPost(GetAPITokensQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to make GraphqlPost request: %v", err)
		}
		data, ok := currentTokens.Data.(map[string]interface{})
		if !ok {
			// Handle the error
			log.Fatal("Data is not a map[string]interface{}")
			return nil, fmt.Errorf("Data is not a map[string]interface{}")
		}

		apitokens, ok := data["apitokens"].([]interface{})
		if !ok {
			// Handle the error
			log.Fatal("apitokens is not a []interface{}")
			return nil, fmt.Errorf("apitokens is not a []interface{}")
		}
		if len(apitokens) > 0 {
			tokenMap, ok := apitokens[0].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("apitoken is not a map[string]interface{}")
			}
			tokenValue, ok := tokenMap["token_value"].(string)
			if !ok {
				return nil, fmt.Errorf("token_value is not a string")
			}
			mythic.APIToken = tokenValue
		} else {
			// If there are no current tokens, create a new one
			newToken, _ := mythic.GraphqlPost(CreateAPITokenMutation)
			data, ok := newToken.Data.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Data is not a map[string]interface{}")
			}
			if statusData, ok := data["createAPIToken"].(map[string]interface{}); ok {
				if statusData["status"].(string) == "success" {
					mythic.APIToken = statusData["token_value"].(string)
				} else {
					errMsg := statusData["error"].(string)
					err := fmt.Errorf("Failed to get or generate an API token to use from Mythic\n%s", errMsg)
					log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
					return nil, err
				}
			}
		}
	}
	return mythic, nil
}
