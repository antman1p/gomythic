package Mythic_Go_Scripting

import (
	"bytes"
	"encoding/json"
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
	"reflect"
	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
	"github.com/gorilla/websocket"

)

func (m *Mythic) GetHTTPTransport() (http.RoundTripper, string) {
	url := fmt.Sprintf("%s://%s:%d/graphql/", m.HTTP, m.ServerIP, m.ServerPort)

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
	}), url
}


func (m *Mythic) GraphqlPost(operation interface{}, variables map[string]interface{}, operationType string) error {
	transport, serverURL := m.GetHTTPTransport()

	client := graphql.NewClient(serverURL, &http.Client{Transport: transport})

	// Check operation type and execute accordingly
	if operationType == "query" {
		err := client.Query(context.Background(), operation, variables)
		if err != nil {
			return err
		}
	} else if operationType == "mutation" {
		err := client.Mutate(context.Background(), operation, variables)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid operation type: %s", operationType)
	}

	return nil
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

    // Use underscore to ignore the URL return value
    transport, _ := m.GetHTTPTransport()

    client := &http.Client{
        Transport: transport,
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

// HeaderToMap is assumed to convert http.Header to a map[string]interface{} type. 
// Replace this with your actual function if it is different.
func HeaderToMap(header http.Header) map[string]interface{} {
	// Implement this function based on your requirements.
	return make(map[string]interface{})
}




func (m *Mythic) HttpPostForm(data url.Values, url string) (map[string]interface{}, error) {
	// Ignore the returned URL using underscore
	transport, _ := m.GetHTTPTransport()

	client := &http.Client{
		Transport: transport,
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
	// Ignore the returned URL using underscore
	transport, _ := m.GetHTTPTransport()
	
	client := &http.Client{
		Transport: transport,
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
	// Ignore the returned URL using underscore
	transport, _ := m.GetHTTPTransport()
	
	client := &http.Client{
		Transport: transport,
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
	// Ignore the returned URL using underscore
	transport, _ := m.GetHTTPTransport()
	
	client := &http.Client{
		Transport: transport,
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


func (m *Mythic) newWebsocketConn(sc *graphql.SubscriptionClient) (graphql.WebsocketConn, error) {
    var endpoint = "/graphql/"
    log.Println("newWebsocketConn endpoint:", endpoint)

    headers := m.GetHeaders()
    log.Println("newWebsocketConn headers:", headers)

    // Prepare the client
    client, err := m.getWebSocketTransport(endpoint)
    if err != nil {
        return nil, err
    }

    return &MythicWebSocketHandler{
        Conn:    client,
        timeout: sc.GetTimeout(),
    }, nil
}

type MythicWebSocketHandler struct {
    Conn    *websocket.Conn
    timeout time.Duration
}

func (h *MythicWebSocketHandler) ReadJSON(v interface{}) error {
    return h.Conn.ReadJSON(v)
}

func (h *MythicWebSocketHandler) WriteJSON(v interface{}) error {
    return h.Conn.WriteJSON(v)
}

func (h *MythicWebSocketHandler) Close() error {
    return h.Conn.Close()
}

func (h *MythicWebSocketHandler) SetReadLimit(limit int64) {
    h.Conn.SetReadLimit(limit)
}

func (h *MythicWebSocketHandler) GetCloseStatus(err error) int32 {
    // You can modify this to return the actual close status if possible
    return 1000  // Normal closure status
}



func (m *Mythic) GraphQLSubscription(ctx context.Context, subscription interface{}, variables map[string]interface{}, timeout int) (<-chan interface{}, error) {
    var endpoint = "/graphql/"
    log.Println("GraphQLSubscription endpoint:", endpoint)

    // Convert headers to map[string]interface{}
    headersMap := make(map[string]interface{})
    for key, values := range m.GetHeaders() {
        if len(values) > 0 {
            headersMap[key] = values[0]
        }
    }

    log.Println("GraphQLSubscription headers:", headersMap)
	
    // Marshal the subscription into JSON
    jsonQuery, err := json.Marshal(subscription)
    if err != nil {
        log.Printf("Error marshaling subscription: %v", err)
        return nil, err
    }

    // Marshal the variables into JSON
    jsonVariables, err := json.Marshal(variables)
    if err != nil {
        log.Printf("Error marshaling variables: %v", err)
        return nil, err
    }

    // Log the subscription and variables
    log.Printf("GraphQL subscription: %s", string(jsonQuery))
    log.Printf("GraphQL variables: %s", string(jsonVariables))

    // Prepare the client
    client := graphql.NewSubscriptionClient(endpoint).
        WithConnectionParams(headersMap).
        WithTimeout(time.Duration(timeout) * time.Second).
        WithWebSocket(m.newWebsocketConn)
		
    log.Printf("DEBUG: Starting subscription with the following parameters:\nSubscription: %+v\nVariables: %+v\nTimeout: %v", subscription, variables, timeout)

    // Create a channel to receive responses
    events := make(chan interface{})
	
    ctx, cancel := context.WithCancel(context.Background())
	
    handleResult := func(data []byte, err error) error {
        if err != nil {
            log.Println("Error in GraphQL subscription:", err)
            return err
        }

        log.Println("DEBUG: Received JSON data: ", string(data)) // add this line

        switch v := subscription.(type) {
        case *TaskWaitForStatusSubscription:
			log.Printf("Subscription is of type: %T\n", v)
            var event TaskWaitForStatusSubscription
            err = jsonutil.UnmarshalGraphQL(data, &event)
            if err != nil {
                log.Println("Error parsing GraphQL subscription event:", err)
                close(events) // close the events channel
                return err
            }
            events <- &event

            // Close the events channel if the task is completed
            for _, task := range event.TaskStream {
                log.Printf("Received event: %+v", event)
                if task.Status == "completed" {
                    close(events)
                    cancel()
                    return nil
                }
            }
        case *TaskWaitForOutputSubscription:
			log.Printf("Subscription is of type: %T\n", v)
            var event TaskWaitForOutputSubscription
			log.Printf("Received raw data: %s", string(data))
			log.Printf("Parsed event: %+v", event)
            err = jsonutil.UnmarshalGraphQL(data, &event)
            if err != nil {
                log.Println("Error parsing GraphQL subscription event:", err)
                close(events) // close the events channel
                return err
            }
            events <- &event

            // you might need to handle task completion logic for this event type here
        // case more subscription types as needed
        default:
            return fmt.Errorf("unsupported subscription type")
        }

        return nil
    }

    // Subscribe with the prepared request
    subscriptionId, err := client.Subscribe(subscription, variables, handleResult)

	if err != nil {
		log.Println("Error in GraphQL subscription:", err)
		return nil, err
	}


    // Run the client in a separate goroutine
	go func() {
		defer close(events)

		running := true
		for running {
			select {
			case <-ctx.Done():
				client.Unsubscribe(subscriptionId) // unsubscribe when context is done
				running = false
			default:
				client.Run()
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

func (mythic *Mythic) SetMythicDetails(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout int) *graphql.Client {
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
	
	url := fmt.Sprintf("%s://%s:%d/graphql/", mythic.HTTP, mythic.ServerIP, mythic.ServerPort)

	var transport http.RoundTripper
	if ssl {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		transport = http.DefaultTransport
	}

	client := graphql.NewClient(url, &http.Client{Transport: transport})

	return client
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
	log.Printf("Sending GraphqlPost request...\n") 

	var query GetAPITokensQuery

	err := mythic.GraphqlPost(&query, map[string]interface{}{}, "query")
	if err != nil {
		//DEBUG
		log.Printf("GraphqlPost ERROR: %s", err)
		return err
	}

	if len(query.APITokens) > 0 {
		//DEBUG 
		log.Printf("query.APITokens > 0: %v", query.APITokens)
		tokenValue := query.APITokens[0].TokenValue
		mythic.APIToken = tokenValue
	} else {
		err := mythic.CreateNewAPIToken()
		if err != nil {
			log.Fatal("Failed to create a new API token: ", err)
			return err
		}
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
	transport, serverURL := mythic.GetHTTPTransport()
	client := graphql.NewClient(serverURL, &http.Client{Transport: transport})

	var response CreateAPITokenMutation
	err := client.Query(context.Background(), &response, map[string]interface{}{})
	if err != nil {
		return err
	}
	

	if response.CreateAPIToken.Status.Equals("completed") {
		mythic.APIToken = response.CreateAPIToken.TokenValue
	} else {
		errMsg := response.CreateAPIToken.Error
		err := fmt.Errorf("Failed to get or generate an API token to use from Mythic\n%s", errMsg)
		log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
		return err
	}

	return nil
}





func (m *Mythic) getWebSocketTransport(path string) (*websocket.Conn, error) {
    u := url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", m.ServerIP, m.ServerPort), Path: path}

    dialer := websocket.Dialer{
        HandshakeTimeout:  time.Minute,
        Subprotocols:      []string{"graphql-ws"},
        TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
        EnableCompression: true,
    }

    headers := http.Header{}
    for key, value := range m.GetHeaders() {
        headers.Add(key, value[0])
    }

    c, _, err := dialer.Dial(u.String(), headers)
    if err != nil {
        return nil, err
    }

    return c, nil
}







func structToMap(obj interface{}) map[string]interface{} {
    out := make(map[string]interface{})
    v := reflect.ValueOf(obj)

    // If pointer get the underlying element≤
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    for i := 0; i < v.NumField(); i++ {
        field := v.Type().Field(i)
        value := v.Field(i).Interface()
        out[field.Name] = value
    }

    return out
}



