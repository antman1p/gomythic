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
	"context"
	
	"github.com/machinebox/graphql"
)

func NewMythic(username, password, serverIP string, serverPort int, apiToken string, ssl bool, timeout int) *Mythic {
	protocol := "http"
	if ssl {
		protocol = "https"
	}

	return &Mythic{
		Username:         username,
		Password:         password,
		APIToken:         apiToken,
		ServerIP:         serverIP,
		ServerPort:       serverPort,
		SSL:              ssl,
		HTTP:             protocol,
		WS:               "ws",
		GlobalTimeout:    timeout,
		ScriptingVersion: "0.1.4",
	}
}

func (m *Mythic) GetHTTPTransport() http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		// Set custom headers here
		req.Header = m.GetHeaders()

		// Use different transport depending on SSL setting
		var transport http.RoundTripper
		if m.SSL {
			transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		} else {
			transport = http.DefaultTransport
		}

		return transport.RoundTrip(req)
	})
}

func (m *Mythic) GraphqlPost(query string, variables map[string]interface{}) (interface{}, error) {
	// Set up GraphQL client
	url := fmt.Sprintf("%s://%s:%d/graphql/", m.HTTP, m.ServerIP, m.ServerPort)

	client := graphql.NewClient(url, graphql.WithHTTPClient(&http.Client{
		Transport: m.GetHTTPTransport(),
	}))

	// Prepare the request
	req := graphql.NewRequest(query)
	
	// Set the headers
	req.Header = m.GetHeaders()

	// Set variables if any
	if variables != nil {
		for key, value := range variables {
			req.Var(key, value)
		}
	}

	// Prepare a context with timeout
	ctx := context.Background()
	if m.GlobalTimeout >= 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(m.GlobalTimeout)*time.Second)
		defer cancel()
	}

	// Execute the request
	var res map[string]interface{}
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	
	//DEBUG
	log.Printf("Response: %v\n", res)


	// Check for errors in the response
	if responseErrors, ok := res["errors"]; ok {
		errMsg := ""
		for _, err := range responseErrors.([]interface{}) {
			errMsg += fmt.Sprintf("%s\n", err.(map[string]interface{})["message"])
		}
		return nil, fmt.Errorf("%s", errMsg)
	}

	// If there are no errors, return the response
	return res, nil
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func (m *Mythic) HttpPost(url string, data map[string]interface{}) (map[string]interface{}, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    client := &http.Client{
        Transport: m.GetHTTPTransport(),
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header = m.GetHeaders()
    req.Header.Set("Content-Type", "application/json")
	


    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
	
	log.Printf("Response body: %s\n", responseData)  // log the response body DEBUG


    var response map[string]interface{}
    err = json.Unmarshal(responseData, &response)
    if err != nil {
        return nil, err
    }
	

    return response, nil
}


func (m *Mythic) GetHeaders() http.Header {
	headers := http.Header{}
	if m.APIToken != "" {
		headers.Set("apitoken", m.APIToken)
	} else if m.AccessToken != ""{
		headers.Set("Authorization", "Bearer "+m.AccessToken)
	}

	return headers
}



func (m *Mythic) HttpPostForm(data url.Values, url string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: m.GetHTTPTransport(),
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = m.GetHeaders()
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

func (m *Mythic) HttpGetDictionary(url string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: m.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = m.GetHeaders()

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

func (m *Mythic) HttpGet(url string) ([]byte, error) {
	client := &http.Client{
		Transport: m.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = m.GetHeaders()

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

func (m *Mythic) HttpGetChunked(url string, chunkSize int) (<-chan []byte, error) {
	client := &http.Client{
		Transport: m.GetHTTPTransport(),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = m.GetHeaders()

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


func (m *Mythic) GraphQLSubscription(query string, variables map[string]interface{}, timeout int) (<-chan map[string]interface{}, error) {
	data := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	response, err := m.HttpPost(m.HTTP+m.ServerIP+":"+strconv.Itoa(m.ServerPort)+"/graphql", data)
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
				response, err := m.HttpGet(m.HTTP + m.ServerIP + ":" + strconv.Itoa(m.ServerPort) + "/graphql/events")
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

func (m *Mythic) FetchGraphQLSchema() (string, error) {
	response, err := m.HttpGet(m.HTTP + m.ServerIP + ":" + strconv.Itoa(m.ServerPort) + "/graphql/schema.json")
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (m *Mythic) LoadMythicSchema() bool {
	schema, err := m.FetchGraphQLSchema()
	if err != nil {
		log.Println("Failed to fetch Mythic schema:", err)
		return false
	}

	m.Schema = schema
	return true
}

func (mythic *Mythic) SetMythicDetails(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout int) {
	mythic.Username = username
	mythic.Password = password
	mythic.ServerIP = serverIP
	mythic.ServerPort = serverPort
	mythic.APIToken = apiToken
	mythic.SSL = ssl
	mythic.GlobalTimeout = timeout
	
	if ssl {
		mythic.HTTP = "https"
	} else {
		mythic.HTTP = "http"
	}
	
	mythic.Schema = "https"
	
	// Set the scripting version here
	mythic.ScriptingVersion = "0.1.4"
	
}


func (mythic *Mythic) AuthenticateToMythic() error {
    url := fmt.Sprintf("%s://%s:%d/auth", mythic.HTTP, mythic.ServerIP, mythic.ServerPort)  
	log.Printf("[*] URL: %s\n", url) // Add this line
	data := map[string]interface{}{
		"username":          mythic.Username,
		"password":          mythic.Password,
		"scripting_version": mythic.ScriptingVersion,
	}
	log.Printf("[*] Logging into Mythic as scripting_version: %s", mythic.ScriptingVersion)
	response, err := mythic.HttpPost(url, data)
	if err != nil {
		log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
		// DEBUG
		responseBody, _ := json.Marshal(response)
		log.Printf("HTTP Response from server: %s\n", responseBody)
		return err
	}

	mythic.AccessToken = response["access_token"].(string)
	mythic.RefreshToken = response["refresh_token"].(string)
	user := response["user"].(map[string]interface{})
	mythic.CurrentOperationID = int(user["current_operation_id"].(float64))

	return nil
}

func (mythic *Mythic) HandleAPITokens() error {
	
	log.Printf("Sending GraphqlPost request...\n") //DEBUG
	
	// Handle data as a generic interface{} first, then check and convert to map or array
	currentTokens, err := mythic.GraphqlPost(GetAPITokensQuery, map[string]interface{}{})

	// Check if error is nil
	if err != nil {
		// Handle error
		log.Printf("GraphqlPost Error: %v", err) // DEBUG
		return fmt.Errorf("failed to make GraphqlPost request: %v", err)
	} else if currentTokens == nil {
		// Handle nil response
		log.Printf("GraphqlPost returned nil response")
		return fmt.Errorf("GraphqlPost returned nil response")
	}
	
	// Try to convert response to a map
	responseMap, ok := currentTokens.(map[string]interface{})
	if !ok {
		log.Fatal("response is not a map[string]interface{}")
		return fmt.Errorf("response is not a map[string]interface{}")
	}

	// Extract 'apitokens' from response
	apitokens, _ := responseMap["apitokens"]

	// Handle apitokens
	switch apitokens := apitokens.(type) {
	case []interface{}:
		if len(apitokens) > 0 {
			// Try to convert the first item to a map
			firstToken, ok := apitokens[0].(map[string]interface{})
			if !ok {
				log.Fatal("first token is not a map[string]interface{}")
				return fmt.Errorf("first token is not a map[string]interface{}")
			}

			// Try to convert the 'token_value' field to a string
			tokenValue, ok := firstToken["token_value"].(string)
			if !ok {
				log.Fatal("token_value is not a string")
				return fmt.Errorf("token_value is not a string")
			}

			// Store the token value in the Mythic struct
			mythic.APIToken = tokenValue
		} else {
			// If there are no current tokens, you could create a new one here
			// Note that you'll need to handle the error from this function
			err := mythic.CreateNewAPIToken()
			if err != nil {
				log.Fatal("Failed to create new API token: ", err)
				return err
			}
		}
	default:
		log.Printf("Unexpected data type: %T\n", apitokens) // Log the actual type of data
		log.Fatal("Data is neither a map nor an array")
		return fmt.Errorf("Data is neither a map nor an array")
	}

	return nil
}

func (mythic *Mythic) HandleAPITokenMap(data map[string]interface{}) error {
	apitokens, ok := data["apitokens"].([]interface{})
	if !ok {
		// Handle the error
		log.Fatal("apitokens is not a []interface{}")
		return fmt.Errorf("apitokens is not a []interface{}")
	}
	if len(apitokens) > 0 {
		return mythic.HandleExistingAPIToken(apitokens)
	} else {
		// If there are no current tokens, create a new one
		return mythic.CreateNewAPIToken()
	}
}

func (mythic *Mythic) HandleExistingAPIToken(apitokens []interface{}) error {
	tokenMap, ok := apitokens[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("apitoken is not a map[string]interface{}")
	}
	tokenValue, ok := tokenMap["token_value"].(string)
	if !ok {
		return fmt.Errorf("token_value is not a string")
	}
	mythic.APIToken = tokenValue

	return nil
}

func (mythic *Mythic) CreateNewAPIToken() error {
    response, _ := mythic.GraphqlPost(CreateAPITokenMutation, map[string]interface{}{})
    
    // Add a type assertion to convert the interface{} to map[string]interface{}
    newToken, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response is not a map[string]interface{}")
	}

	createAPIToken, ok := newToken["createAPIToken"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("createAPIToken is not a map[string]interface{}")
	}

	if statusData := createAPIToken; statusData["status"].(string) == "success" {
		mythic.APIToken = statusData["token_value"].(string)
	} else {
		errMsg := statusData["error"].(string)
		err := fmt.Errorf("Failed to get or generate an API token to use from Mythic\n%s", errMsg)
		log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
		return err
	}

    return nil
}