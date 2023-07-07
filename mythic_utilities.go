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


// GetHTTPTransport function sets up the HTTP transport using the defined server configuration in Mythic struct. 
// It also sets up custom headers for the requests using the GetHeaders function.
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

// GraphqlPost function performs a GraphQL query or mutation using the Hasura GraphQL client.
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


// roundTripperFunc is a type definition for a function that acts as a roundTripper in the http client transport.
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip function calls the roundTripper function.
func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// HttpPost function sends a POST request to the specified URL with the provided data as JSON. 
// The request is sent using the HTTP transport set up in the GetHTTPTransport function.
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

    var response map[string]interface{}
    err = json.Unmarshal(responseData, &response)
    if err != nil {
        return nil, err
    }

    return response, nil
}


// GetHeaders function sets up the custom headers for the request based on the provided API token or Access token in the Mythic struct.
func (m *Mythic) GetHeaders() http.Header {
	headers := http.Header{}
	if m.APIToken != "" {
		headers.Set("Apitoken", strings.TrimSpace(m.APIToken))
	} else if m.AccessToken != "" {
		headers.Set("Authorization", "Bearer "+strings.TrimSpace(m.AccessToken))
	}

	return headers
}

// HeaderToMap is assumed to convert http.Header to a map[string]interface{} type. 
// Replace this with your actual function if it is different.
func HeaderToMap(header http.Header) map[string]interface{} {
	// Implement this function based on your requirements.
	return make(map[string]interface{})
}


// UNTESTED
// HttpPostForm function sends a POST request to the specified URL with the provided data as form values.
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

// UNTESTED
// HttpGetDictionary function sends a GET request to the specified URL and returns the response as a dictionary.
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

// UNTESTED
// HttpGet function sends a GET request to the specified URL and returns the raw response.
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

// UNTESTED
// HttpGetChunked function sends a GET request to the specified URL and returns the response in chunks.
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

// newWebsocketConn function sets up the WebSocket connection for GraphQL subscription.
func (m *Mythic) newWebsocketConn(sc *graphql.SubscriptionClient) (graphql.WebsocketConn, error) {
    var endpoint = "/graphql/"
	
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


// ReadJSON function reads JSON data from the WebSocket connection.
func (h *MythicWebSocketHandler) ReadJSON(v interface{}) error {
    return h.Conn.ReadJSON(v)
}

// WriteJSON function writes JSON data to the WebSocket connection.
func (h *MythicWebSocketHandler) WriteJSON(v interface{}) error {
    return h.Conn.WriteJSON(v)
}

// Close function closes the WebSocket connection.
func (h *MythicWebSocketHandler) Close() error {
    return h.Conn.Close()
}

// SetReadLimit function sets the read limit for the WebSocket connection.
func (h *MythicWebSocketHandler) SetReadLimit(limit int64) {
    h.Conn.SetReadLimit(limit)
}

// GetCloseStatus function gets the close status of the WebSocket connection.
func (h *MythicWebSocketHandler) GetCloseStatus(err error) int32 {
    // You can modify this to return the actual close status if possible
    return 1000  // Normal closure status
}


// GraphQLSubscription function sets up the GraphQL subscription using the provided parameters.
func (m *Mythic) GraphQLSubscription(ctx context.Context, subscription interface{}, variables map[string]interface{}, timeout int) (<-chan interface{}, error) {
    var endpoint = "/graphql/"

    // Convert headers to map[string]interface{}
    headersMap := make(map[string]interface{})
    for key, values := range m.GetHeaders() {
        if len(values) > 0 {
            headersMap[key] = values[0]
        }
    }

    // Prepare the client
    client := graphql.NewSubscriptionClient(endpoint).
        WithConnectionParams(headersMap).
        WithTimeout(time.Duration(timeout) * time.Second).
        WithWebSocket(m.newWebsocketConn)
		

    // Create a channel to receive responses
    events := make(chan interface{})
	
    ctx, cancel := context.WithCancel(context.Background())
	
    handleResult := func(data []byte, err error) error {
        if err != nil {
            log.Println("Error in GraphQL subscription:", err)
            return err
        }


        switch subscription.(type) {
        case *TaskWaitForStatusSubscription:
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
                if task.Status == "completed" {
                    close(events)
                    cancel()
                    return nil
                }
            }
        case *TaskWaitForOutputSubscription:
            var event TaskWaitForOutputSubscription
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

// UNTESTED
// FetchGraphQLSchema function sends an HTTP GET request to the Mythic server to retrieve the GraphQL schema and returns it as a string.
func (m *Mythic) FetchGraphQLSchema() (string, error) {
	response, err := m.HttpGet(m.HTTP + m.ServerIP + ":" + strconv.Itoa(m.ServerPort) + "/graphql/schema.json")
	if err != nil {
		return "", err
	}
	return string(response), nil
}

//UNTESTED
// LoadMythicSchema function fetches the GraphQL schema from the Mythic server and stores it in the Mythic struct.
func (m *Mythic) LoadMythicSchema() bool {
	schema, err := m.FetchGraphQLSchema()
	if err != nil {
		log.Println("Failed to fetch Mythic schema:", err)
		return false
	}

	m.Schema = schema
	return true
}

// SetMythicDetails function sets up the Mythic struct with the provided server details and returns a new GraphQL client.
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

// AuthenticateToMythic function sends a POST request to authenticate to the Mythic server with the provided username and password.
func (mythic *Mythic) AuthenticateToMythic() error {
    url := fmt.Sprintf("%s://%s:%d/auth", mythic.HTTP, mythic.ServerIP, mythic.ServerPort)  
	data := map[string]interface{}{
		"username":          mythic.Username,
		"password":          mythic.Password,
		"scripting_version": mythic.ScriptingVersion,
	}
	response, err := mythic.HttpPost(url, data)
	if err != nil {
		log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
		return err
	}

	mythic.AccessToken = response["access_token"].(string)
	mythic.RefreshToken = response["refresh_token"].(string)
	user := response["user"].(map[string]interface{})
	mythic.CurrentOperationID = int(user["current_operation_id"].(float64))

	return nil
}

// HandleAPITokens function fetches the API tokens associated with the user from the Mythic server using a GraphQL query.
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


// CreateNewAPIToken function sends a GraphQL mutation to the Mythic server to create a new API token.
func (m *Mythic) CreateNewAPIToken() error {

	variables := CreateAPITokenVariables{
		TokenType: "User",
	}
	
	variableMap := map[string]interface{}{
		"token_type": variables.TokenType,
	}
	
	var response CreateAPITokenMutation
	err := m.GraphqlPost(&response, variableMap, "mutation")
	if err != nil {
		log.Printf("[-] Failed to execute mutation: \n%s", err)
		return err
	}

	
	if response.CreateAPIToken.Status == "success" {
		m.APIToken = response.CreateAPIToken.TokenValue
	} else {
		errMsg := response.CreateAPIToken.Error
		err := fmt.Errorf("Failed to get or generate an API token to use from Mythic\n%s", errMsg)
		log.Printf("[-] Failed to authenticate to Mythic: \n%s", err)
		return err
	}

	return nil
}

// getWebSocketTransport function sets up a WebSocket connection to the Mythic server with the specified path.
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

// structToMap function converts a struct to a map.
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


func FilterResponse(data []interface{}, fields []string) ([]map[string]interface{}, error) {
    filtered := []map[string]interface{}{}

    for _, item := range data {
        entry := make(map[string]interface{})
        itm := structToMap(item) // Convert struct to map
        for _, field := range fields {
            if val, ok := itm[field]; ok {
                entry[field] = val
            }
        }
        filtered = append(filtered, entry)
    }

    return filtered, nil
}

// Convert []Callback to []interface{}
func CallbacksToInterfaces(callbacks []Callback) []interface{} {
    res := make([]interface{}, len(callbacks))
    for i, v := range callbacks {
        res[i] = v
    }
    return res
}

// Convert []TaskFragment to []interface{}
func TasksToInterfaces(tasks []TaskFragment) []interface{} {
    res := make([]interface{}, len(tasks))
    for i, v := range tasks {
        res[i] = v
    }
    return res
}

func TaskFragmentsToInterfaces(taskFragments []TaskFragment) []interface{} {
    result := make([]interface{}, len(taskFragments))
    for i, v := range taskFragments {
        result[i] = v
    }
    return result
}





