package gomythic

import (
	"log"
	"fmt"
	"encoding/json"
	"time"
	"context"
	"encoding/base64"
	"strings"
	"errors"
	"net/http"
	
	"github.com/hasura/go-graphql-client"

)

// Login sets up a new Mythic instance and tries to authenticate with the provided credentials or API token.
func Login(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout, loggingLevel int) (*Mythic, error) {
	mythic := &Mythic{}
	
	// SetMythicDetails initializes the server details for the new instance.
	mythic.SetMythicDetails(serverIP, serverPort, username, password, apiToken, ssl, timeout)

	// These lines configure the logger to print the filename and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// If an API token is not provided, the method tries to authenticate with the provided username and password.
	if apiToken == "" {
		// AuthenticateToMythic sends a login request to the server.
		if err := mythic.AuthenticateToMythic(); err != nil {
			return nil, err
		}
		// HandleAPITokens tries to fetch the API token from the server.
		if err := mythic.HandleAPITokens(); err != nil {
			return nil, err
		}
	} else {
		// If an API token is provided, it is directly used for authentication.
		mythic.APIToken = apiToken
	}
	return mythic, nil
}

// ExecuteCustomQuery sends a custom GraphQL query to the server and unmarshals the response into result.
func (m *Mythic) ExecuteCustomQuery(query string, variables map[string]interface{}, result interface{}) error {
	// The method first checks if the provided query is not empty.
	if strings.TrimSpace(query) == "" {
		return errors.New("query string is empty")
	}

	// GetHTTPTransport gets the http.Transport and server URL.
	transport, serverURL := m.GetHTTPTransport()

	// A new GraphQL client is initialized with the server URL and http.Transport.
	client := graphql.NewClient(serverURL, &http.Client{Transport: transport})

	// The context for the request is initialized.
	ctx := context.Background()

	// The query is executed with the provided variables and result.
	err := client.Exec(ctx, query, result, variables)
	if err != nil {
		log.Printf("Hit an exception within ExecuteCustomQuery: %v", err)
		return err
	}

	return nil
}


// TODO:
// func (m *Mythic) SubscribeCustomQuery(query string, variables map[string]interface{}, timeout int) error { }



// # ########### Callback Functions #############

// GetAllCallbacks sends a GraphQL query to fetch all callbacks from the server.
func (m *Mythic) GetAllCallbacks(customReturnAttributes []string) ([]map[string]interface{}, error) {
	var query CallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}

	// If customReturnAttributes is not nil, filter the response.
	if customReturnAttributes != nil {
		filteredCallbacks, err := FilterResponse(CallbacksToInterfaces(query.Callback), customReturnAttributes)		
		if err != nil {
			return nil, err
		}
		return filteredCallbacks, nil
	}
	
	// Else, convert the Callbacks to map[string]interface{} and return.
	callbacksMap := make([]map[string]interface{}, len(query.Callback))
	for i, callback := range query.Callback {
		jsonCb, _ := json.Marshal(callback)
		cbMap := make(map[string]interface{})
		json.Unmarshal(jsonCb, &cbMap)
		callbacksMap[i] = cbMap
	}
	
	return callbacksMap, nil
}


// GetAllActiveCallbacks sends a GraphQL query to fetch all active callbacks from the server.
func (m *Mythic) GetAllActiveCallbacks(customReturnAttributes []string) ([]map[string]interface{}, error) {
	var query ActiveCallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}
	
	// If customReturnAttributes is not nil, filter the response.
	if customReturnAttributes != nil {
	filteredCallbacks, err := FilterResponse(CallbacksToInterfaces(query.Callback), customReturnAttributes)
	if err != nil {
		return nil, err
	}
		return filteredCallbacks, nil
	}

	// Convert the Callbacks to map[string]interface{} and return.
	callbacksMap := make([]map[string]interface{}, len(query.Callback))
	for i, callback := range query.Callback {
		jsonCb, _ := json.Marshal(callback)
		cbMap := make(map[string]interface{})
		json.Unmarshal(jsonCb, &cbMap)
		callbacksMap[i] = cbMap
	}

	return callbacksMap, nil
}


//TODO:
// func (m *Mythic) SubscribeNewCallbacks(batchSize int, timeout int, customReturnAttributes []string)) ([]Callback, error){}

// func (m *Mythic) SubscribeAllActiveCallbacks(timeout int, customReturnAttributes []string)) ([]Callback, error){}

/* func (m *Mythic) UpdateCallback(callbackDisplayID *int, active bool, sleepInfo string, locaked bool, description string,
	ips []string, user string, host string, os string, architecture string, extraInfo string, pid int, processName string,
	integrityLevel int, domain string) ([]Callback, error){}
*/


// ############ Task Functions #################

// GetAllTasks sends a GraphQL query to fetch all tasks associated with a callback from the server.
func (m *Mythic) GetAllTasks(callbackDisplayID *int, customReturnAttributes []string) ([]map[string]interface{}, error) {
    // Depending on whether a callback display ID is provided, a different GraphQL query is sent.
    if callbackDisplayID != nil {
        var query TaskQueryWithCallback
        err := m.GraphqlPost(&query, map[string]interface{}{
            "callbackDisplayID": *callbackDisplayID,
        }, "query")
        if err != nil {
            return nil, err
        }

        // If customReturnAttributes is not nil, filter the response.
        if customReturnAttributes != nil {
            filteredTasks, err := FilterResponse(TaskFragmentsToInterfaces(query.Task), customReturnAttributes)
            if err != nil {
                return nil, err
            }
            return filteredTasks, nil
        }

        // Else, convert the Tasks to map[string]interface{} and return.
        tasksMap := make([]map[string]interface{}, len(query.Task))
        for i, task := range query.Task {
            jsonTask, _ := json.Marshal(task)
            taskMap := make(map[string]interface{})
            json.Unmarshal(jsonTask, &taskMap)
            tasksMap[i] = taskMap
        }
        
        return tasksMap, nil
    } else {
        var query TaskQuery
        err := m.GraphqlPost(&query, nil, "query")
        if err != nil {
            return nil, err
        }

        // If customReturnAttributes is not nil, filter the response.
        if customReturnAttributes != nil {
            filteredTasks, err := FilterResponse(TaskFragmentsToInterfaces(query.Task), customReturnAttributes)
            if err != nil {
                return nil, err
            }
            return filteredTasks, nil
        }

        // Else, convert the Tasks to map[string]interface{} and return.
        tasksMap := make([]map[string]interface{}, len(query.Task))
        for i, task := range query.Task {
            jsonTask, _ := json.Marshal(task)
            taskMap := make(map[string]interface{})
            json.Unmarshal(jsonTask, &taskMap)
            tasksMap[i] = taskMap
        }
        
        return tasksMap, nil
    }
}


// IssueTask sends a task to a callback and optionally waits for it to complete.
func (m *Mythic) IssueTask(commandName string, parameters interface{}, callbackDisplayID int, tokenID *int, originalParams interface{}, parameterGroupName interface{}, waitForComplete bool, timeout *int) (*CreateTaskMutation, error) {
	// The parameters are first converted to a JSON string.
	var parameterString string
	switch parameters := parameters.(type) {
	case string:
		parameterString = parameters
	case map[string]interface{}:
		parametersBytes, err := json.Marshal(parameters)
		if err != nil {
			return nil, err
		}
		parameterString = string(parametersBytes)
	default:
		return nil, fmt.Errorf("parameters must be a string or map[string]interface{}")
	}
	
	// Depending on the type of parameters, the tasking location is set.
	taskingLocation := "command_line"
	if _, ok := parameters.(map[string]interface{}); ok {
		taskingLocation = "scripting"
	}
	
	// A new task mutation is created and sent to the server.
	var query CreateTaskMutation
	
	variables := map[string]interface{}{
		"callback_id":             callbackDisplayID,
		"command":                 commandName,
		"params":                  parameterString,
		"tasking_location":        taskingLocation,
	}
	
	if tokenID != nil {
    variables["token_id"] = *tokenID
	} else {
		variables["token_id"] = 0
	}

	if originalParams != "" {
		variables["original_params"] = originalParams
	} else {
		variables["original_params"] = ""
	}

	if parameterGroupName != "" {
		variables["parameter_group_name"] = parameterGroupName
	} else {
		variables["parameter_group_name"] = ""
	}

	err := m.GraphqlPost(&query, variables, "mutation")
	if err != nil {
		return nil, err
	}

	// If the task is created successfully and the caller requested to wait for it to complete, the method waits for the task to complete.
	if query.CreateTask.Status.Equals("success") {
		if waitForComplete {
			taskDisplayID := query.CreateTask.DisplayID
			if err != nil {
				return nil, fmt.Errorf("failed to convert display_id to integer: %v", err)
			}

			taskResult, err := m.WaitForTaskComplete(taskDisplayID, nil, timeout) // Assuming you don't need customReturnAttributes
			if err != nil {
				return nil, fmt.Errorf("failed to wait for task complete: %v", err)
			}
			return taskResult, nil
		}
		return &query, nil
	}

	return nil, fmt.Errorf("failed to create task: %s", query.CreateTask.Error)
}

// WaitForTaskComplete subscribes to task updates and waits for the task to complete or fail.
func (m *Mythic) WaitForTaskComplete(taskDisplayID int, customReturnAttributes *string, timeout *int) (*CreateTaskMutation, error) {
    var subscription TaskWaitForStatusSubscription

    variables := TaskWaitForStatusSubscriptionVariables{
        DisplayID: taskDisplayID,
    }

    variableMap := structToMap(variables)

	// The method subscribes to task updates with a given timeout.
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    results, err := m.GraphQLSubscription(ctx, &subscription, variableMap, *timeout)
    if err != nil {
        log.Printf("Error subscribing to task updates: %v", err)
        return nil, err
    }
	
	// The method waits for the task to complete or fail and returns the task result.
    for result := range results {
        // Here we are asserting that the result is of type *TaskWaitForStatusSubscription
        resultAssertion, ok := result.(*TaskWaitForStatusSubscription)
        if !ok {
            return nil, fmt.Errorf("unexpected type from results channel")
        }

        for _, taskFragment := range resultAssertion.TaskStream {

            if taskFragment.Status.Equals("error") || taskFragment.Status.Equals("completed") {
                taskResult := &CreateTaskMutation{
                    CreateTask: struct {
                        Status    MythicStatus `graphql:"status"`
                        ID        int `graphql:"id"`
                        DisplayID int `graphql:"display_id"`
                        Error     string `graphql:"error"` // assign the Error value here
                    } {
                        Status:    taskFragment.Status,
                        ID:        taskFragment.ID,
                        DisplayID: taskFragment.DisplayID,
                        Error:     "", // this is a placeholder. Please replace this with the actual Error field from your taskFragment if it exists
                    },
                }
                // Logging the completion status of the task
                log.Printf("Task completed with status: %s", taskResult.CreateTask.Status)
                return taskResult, nil
            }
        }
    }
    log.Printf("Task did not complete within the given timeout")
    return nil, fmt.Errorf("task not completed")
}


// IssueTaskAndWaitForOutput sends a task to a callback, waits for it to complete, and then retrieves its output.
func (m *Mythic) IssueTaskAndWaitForOutput(commandName string, parameters interface{}, callbackDisplayId int, tokenId int, originalParams interface{}, parameterGroupName interface{}, waitForComplete bool, timeout int) ([]byte, error) {
	tokenIdPtr := &tokenId
	task, err := m.IssueTask(commandName, parameters, callbackDisplayId, tokenIdPtr, originalParams, parameterGroupName, true, &timeout)
	if err != nil {
		return nil, err
	}
	
	taskDisplayId := task.CreateTask.DisplayID
	if taskDisplayId == 0 {
		return nil, fmt.Errorf("invalid task display id")
	}

	return m.WaitForTaskOutput(taskDisplayId, &timeout)
}


// TODO:

// func (m *Mythic) IssueTaskAllActiveCallbacks(commandName string, parameters interface{}) ([]byte, error) {}

// func (m *Mythic) SubscribeNewTasks(batchSize int, timeout int, callbackDisplayId int, customReturnAttributes []string)) ([]byte, error) {}

// func (m *Mythic) SubscribeNewTasksAndUpdates(batchSize int, timeout int, callbackDisplayId int, customReturnAttributes []string)) ([]byte, error) {}

// func (m *Mythic) SubscribeAllTasks(timeout int, callbackDisplayId int, customReturnAttributes []string)) ([]byte, error) {}

// func (m *Mythic) SubscribeAllTasksAndUpdates(timeout int, callbackDisplayId int, customReturnAttributes []string)) ([]byte, error) {}

// func (m *Mythic) AddMitreAttackToTask(timeout int, taskDisplayID int, MitreAttackNumbers []string ) ([]byte, error) {}





// # ######### File Browser Functions ###########

//func (m *Mythic) GetAllFileBrowser(host string, batchSize int, customReturnAttributes []string)) ([]byte, error) {}

//func (m *Mythic) SubscribeNewFileBrowser(host string, batchSize int, timeout int, customReturnAttributes []string)) ([]byte, error) {}

//func (m *Mythic) SubscribeAllFileBrowser(host string, timeout int, batchSize int, customReturnAttributes []string)) ([]byte, error) {}




// # ######### Command Functions ##############

// func (m *Mythic) GetAllCommandsForPayloadType(payloadTypeName string, customReturnAttributes []string)) ([]byte, error) {}



// # ######### Payload Functions ##############

/* func (m *Mythic) CreatePayload(payloadTypeName string, filename string, operatingSystem string, c2Profiles interface{}, commands []string, 
	buildParameters interface{}, description string, returnOnComplete bool, timeout int, customReturnAttributes []string), includeAllCommands bool) ([]byte, error) {}
*/

/*
func (m *Mythic) CreateWrapperPayload(payloadTypeName string, filename string, operatingSystem string, wrappedPayloadUUID string,
	buildParameters interface{}, description string, returnOnComplete bool, timeout int, customReturnAttributes []string) ([]byte, error) {}
*/

// func (m *Mythic) WaitForPayloadComplete(payloadUUID string, timeout int, customReturnAttributes []string) ([]byte, error) {}

// func (m *Mythic) GetAllPayloads(customReturnAttributes []string) ([]byte, error) {}

// func (m *Mythic) GetPayloadByUUID(payloadUUID string, customReturnAttributes []string) ([]byte, error) {}

// func (m *Mythic) DownloadPayload(payloadUUID string) ([]byte, error) {}

// func (m *Mythic) PayloadCheckConfig(payloadUUID string) ([]byte, error) {}

// func (m *Mythic) PayloadRedirectRules(payloadUUID string) ([]byte, error) {}


// # ######### Task Output Functions ###########


/* This function is responsible for waiting for the output of a specific task
   identified by its display ID.
   The output will be collected and returned as a byte array.
   A timeout can be provided to stop waiting after a specific duration.
*/ 
func (m *Mythic) WaitForTaskOutput(taskDisplayID int, timeout *int) ([]byte, error) {
	// A subscription object for getting task output
    var subscription TaskWaitForOutputSubscription

	// Variables required for the subscription
    variables := TaskWaitForOutputSubscriptionVariables{
        DisplayID: taskDisplayID,
    }
	// Convert the variables to a map
    variableMap := structToMap(variables)
	
	// Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10) // Set timeout to 10 seconds
    defer cancel() // Make sure to cancel the context when we're done to free resources
	
	// Subscribe to the task's updates
    events, err := m.GraphQLSubscription(ctx, &subscription, variableMap, *timeout)
    if err != nil {
        log.Printf("Error subscribing to task updates: %v", err)
        return nil, err
    }
	
	// This will hold the final output
    finalOutput := make([]byte, 0)

    for {
        select {
		// When a new event is received
        case event := <-events: 
            v, ok := event.(*TaskWaitForOutputSubscription)
            if !ok {
                return nil, fmt.Errorf("unexpected type: %T", event)
            }
			// Loop over all task streams in the event
            for _, taskStream := range v.TaskStream { 
				// Loop over all responses in the task stream
                for _, response := range taskStream.Responses { 
					// Decode the response text from base64
                    outputBytes, err := base64.StdEncoding.DecodeString(response.ResponseText) 
                    if err != nil {
                        return nil, fmt.Errorf("failed to decode base64 response text: %v", err)
                    }
					// Append the output to the final output
                    finalOutput = append(finalOutput, outputBytes...) 
                }
            }
		// When the context is done (timeout)
        case <-ctx.Done():
            // Now retrieve all subtask IDs
            subtaskIds, err := m.GetAllSubtaskIDs(taskDisplayID, true)
            if err != nil {
                return nil, fmt.Errorf("failed to get subtask IDs: %v", err)
            }
            // Get and aggregate the output of all subtasks
            for _, subtaskId := range subtaskIds {
                subtaskOutput, err := m.GetAllTaskOutputByID(subtaskId)
                if err != nil {
                    return nil, fmt.Errorf("failed to get subtask output: %v", err)
                }
                for _, r := range subtaskOutput {
                    outputBytes, err := base64.StdEncoding.DecodeString(r.ResponseText)
                    if err != nil {
                        return nil, fmt.Errorf("failed to decode base64 response text: %v", err)
                    }
                    finalOutput = append(finalOutput, outputBytes...)
                }
            }
			// Return the final output
            return finalOutput, nil
        }
    }
}


// This function gets all the subtask IDs for a particular task, identified by its DisplayID.
// If fetchDisplayIDInstead is true, it fetches DisplayIDs instead of normal IDs.
func (m *Mythic) GetAllSubtaskIDs(taskDisplayID int, fetchDisplayIDInstead bool) ([]int, error) {
	
	type TaskIdFromDisplayID struct {
		Task []struct {
			ID int `graphql:"id"`
		} `graphql:"task(where: {parent_task_id: {_eq: $task_id}})"`
	}
	
	type SubtaskList struct {
		Task []struct {
			ID         int `graphql:"id"`
			DisplayID int `graphql:"display_id"`
		} `graphql:"task(where: {parent_task_id: {_eq: $task_id}})"`
	}
	
	// Setting the variables for the GraphQL query
	variables := map[string]interface{}{
		"task_id": taskDisplayID,
	}

	// Initialize the initial struct to get the parent task ID
	var initial TaskIdFromDisplayID
	if err := m.GraphqlPost(&initial, variables, "query"); err != nil {
		return nil, err
	}

	subtaskIds := []int{}
	taskIdsToCheck := []int{}
	if len(initial.Task) > 0 {
		// Add the parent task ID to the list of task IDs to check
		taskIdsToCheck = append(taskIdsToCheck, initial.Task[0].ID)
	}
	
	// Loop through all task IDs to check
	for len(taskIdsToCheck) > 0 {
		currentTaskId := taskIdsToCheck[len(taskIdsToCheck)-1]
		taskIdsToCheck = taskIdsToCheck[:len(taskIdsToCheck)-1]
		
		// Update the 'task_id' in variables map
		variables["task_id"] = currentTaskId  // update the 'task_id' in variables map
		
		// Initialize the subtasks struct
		var subtasks SubtaskList
		if err := m.GraphqlPost(&subtasks, variables, "query"); err != nil {
			return nil, err
		}
		
		// Add all the subtask IDs or DisplayIDs to the list
		for _, t := range subtasks.Task {
			taskIdsToCheck = append(taskIdsToCheck, t.ID)
			if fetchDisplayIDInstead {
				subtaskIds = append(subtaskIds, t.DisplayID)
			} else {
				subtaskIds = append(subtaskIds, t.ID)
			}
		}
	}
	// Return the list of subtask IDs
	return subtaskIds, nil
}

// This function retrieves all the output for a specific task, identified by its DisplayID.
func (m *Mythic) GetAllTaskOutputByID(taskDisplayID int) ([]TaskOutputFragment, error) {
	var taskOutput TaskOutput
	
	// Setting the variables for the GraphQL query
	variables := map[string]interface{}{
		"task_display_id": taskDisplayID,
	}
	
	// Execute the GraphQL query and store the result in taskOutput
	err := m.GraphqlPost(&taskOutput, variables, "query")
	if err != nil {
		return nil, err
	}
	// Return the task output
	return taskOutput.Response, nil
}


// TODO:

// func (m *Mythic) GetAllTaskOutput(customReturnAttributes []string, batchSize int) ([]byte, error) {}

// func (m *Mythic) GetAllTaskandSubTaskOutputByID(taskDisplayID int, customReturnAttributes []string) ([]byte, error) {}

// func (m *Mythic) SubacribeNewTaskOutput(timeout int, customReturnAttributes []string, batchSize int) ([]byte, error) {}

// func (m *Mythic) SubacribeAllTaskOutput(timeout int, customReturnAttributes []string, batchSize int) ([]byte, error) {}



// TODO:

// # ########## Operator Functions ##############

// func (m *Mythic) CreateOperator(username string, password string) ([]byte, error) {}


// CreateAPIToken function sends a GraphQL mutation to the Mythic server to create a new API token.
func (m *Mythic) CreateAPIToken() error {

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


// TODO:

// func (m *Mythic) SetAdminStatus() ([]byte, error) {}

// func (m *Mythic) SetActiveStatus(username string, active bool) ([]byte, error) {}

// func (m *Mythic) SetPassword(username string, oldPassword string, newPassword string) ([]byte, error) {}

// func (m *Mythic) GetOperator(username string, customReturnAttributes []string) ([]byte, error) {}

// func (m *Mythic) getMe() ([]byte, error) {}



// # ########## File Functions ##############

// func (m *Mythic) RegisterFile() ([]byte, error) {}
// func (m *Mythic) DownloadFile() ([]byte, error) {}
// func (m *Mythic) DownloadFileChunked() ([]byte, error) {}
// func (m *Mythic) GetAllDownloadedFiles() ([]byte, error) {}
// func (m *Mythic) SubscribeNewDownloadedFiles() ([]byte, error) {}
// func (m *Mythic) SubscribeAllDownloadedFiles() ([]byte, error) {}
// func (m *Mythic) GetAllScreenshots() ([]byte, error) {}
// func (m *Mythic) GetAllUploadedFiles() ([]byte, error) {}
// func (m *Mythic) GetLatestUploadedFileByName() ([]byte, error) {}


// # ########## Operations Functions #############

// func (m *Mythic) GetOperations() ([]byte, error) {}
// func (m *Mythic) CreateOperation() ([]byte, error) {}
// func (m *Mythic) AddOperatorToOperation() ([]byte, error) {}
// func (m *Mythic) RemoveOperatorFromOperation() ([]byte, error) {}
// func (m *Mythic) UpdateOperatorInOperation() ([]byte, error) {}
// func (m *Mythic) UpdateOperation() ([]byte, error) {}
// func (m *Mythic) UpdateCurrentOperationForUser() ([]byte, error) {}


// # ############ Process Functions ##############

// func (m *Mythic) SubscribeNewProcess() ([]byte, error) {}
// func (m *Mythic) GetAllProcesses() ([]byte, error) {}
// func (m *Mythic) SubscribeAllProcesses() ([]byte, error) {}


// # ####### Credential Functions ############

// func (m *Mythic) CreateCredential() ([]byte, error) {}


// # ####### Analytic-based Functions ############

// func (m *Mythic) GetUniqueCompromisedHosts() ([]byte, error) {}
// func (m *Mythic) GetUniqueCompromisedAccounts() ([]byte, error) {}
// func (m *Mythic) GetUniqueCompromisedIPs() ([]byte, error) {}


// # ####### Event Feed functions ############

// func (m *Mythic) SendEventLogMessage() ([]byte, error) {}

// # ####### webhook ############

// func (m *Mythic) SendCustomWebHookMessage() ([]byte, error) {}

// # ####### C2 Functions #############

// func (m *Mythic) StartStopC2Profiles() ([]byte, error) {}
// func (m *Mythic) CreateSavedC2Instance() ([]byte, error) {}

// # ####### Tag Functions ############

// func (m *Mythic) CreateTagTypes() ([]byte, error) {}
// func (m *Mythic) UpdateTagTypes() ([]byte, error) {}
// func (m *Mythic) DeleteTagTypes() ([]byte, error) {}
// func (m *Mythic) GetTagTypes() ([]byte, error) {}
// func (m *Mythic) GetAllTagTypes() ([]byte, error) {}






