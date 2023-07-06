package Mythic_Go_Scripting

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

func Login(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout, loggingLevel int) (*Mythic, error) {
	mythic := &Mythic{}
	mythic.SetMythicDetails(serverIP, serverPort, username, password, apiToken, ssl, timeout)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if apiToken == "" {
		if err := mythic.AuthenticateToMythic(); err != nil {
			return nil, err
		}

		if err := mythic.HandleAPITokens(); err != nil {
			return nil, err
		}
	} else {
		mythic.APIToken = apiToken
	}
	return mythic, nil
}


func (m *Mythic) ExecuteCustomQuery(query string, variables map[string]interface{}, result interface{}) error {
	// Check if the query string is empty
	if strings.TrimSpace(query) == "" {
		return errors.New("query string is empty")
	}

	// Get the endpoint and http.Client
	transport, serverURL := m.GetHTTPTransport()

	// Create a new client
	client := graphql.NewClient(serverURL, &http.Client{Transport: transport})

	ctx := context.Background()

	// Execute the query
	err := client.Exec(ctx, query, result, variables)
	if err != nil {
		log.Printf("Hit an exception within ExecuteCustomQuery: %v", err)
		return err
	}

	return nil
}


// func (m *Mythic) SubscribeCustomQuery(query string, variables map[string]interface{}, timeout int) error { }



// # ########### Callback Functions #############

// CallbackAttributes represents the returned data structure of a callback
type CallbackAttributes map[string]interface{}

func (m *Mythic) GetAllCallbacks() ([]Callback, error) {
	var query CallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}

	return query.Callback, nil
}




func (m *Mythic) GetAllActiveCallbacks() ([]Callback, error) {
	var query ActiveCallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}

	return query.Callback, nil
}

//TODO:
// func (m *Mythic) SubscribeNewCallbacks(batchSize int, timeout int) ([]Callback, error){}

// func (m *Mythic) SubscribeAllActiveCallbacks(timeout int) ([]Callback, error){}

/* func (m *Mythic) UpdateCallback(callbackDisplayID *int, active bool, sleepInfo string, locaked bool, description string,
	ips []string, user string, host string, os string, architecture string, extraInfo string, pid int, processName string,
	integrityLevel int, domain string) ([]Callback, error){}
*/

// ############ Task Functions #################


func (m *Mythic) GetAllTasks(callbackDisplayID *int) ([]TaskFragment, error) {
	if callbackDisplayID != nil {
		var query TaskQueryWithCallback
		err := m.GraphqlPost(&query, map[string]interface{}{
			"callbackDisplayID": *callbackDisplayID,
		}, "query")
		if err != nil {
			return nil, err
		}
		return query.Task, nil
	} else {
		var query TaskQuery
		err := m.GraphqlPost(&query, nil, "query")
		if err != nil {
			return nil, err
		}
		return query.Task, nil
	}
}


func (m *Mythic) IssueTask(commandName string, parameters interface{}, callbackDisplayID int, tokenID *int, originalParams interface{}, parameterGroupName interface{}, waitForComplete bool, timeout *int) (*CreateTaskMutation, error) {
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

	taskingLocation := "command_line"
	if _, ok := parameters.(map[string]interface{}); ok {
		taskingLocation = "scripting"
	}

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


func (m *Mythic) WaitForTaskComplete(taskDisplayID int, customReturnAttributes *string, timeout *int) (*CreateTaskMutation, error) {
    var subscription TaskWaitForStatusSubscription

    variables := TaskWaitForStatusSubscriptionVariables{
        DisplayID: taskDisplayID,
    }

    variableMap := structToMap(variables)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    results, err := m.GraphQLSubscription(ctx, &subscription, variableMap, *timeout)
    if err != nil {
        log.Printf("Error subscribing to task updates: %v", err)
        return nil, err
    }
	

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
                
                log.Printf("Task completed with status: %s", taskResult.CreateTask.Status)
                return taskResult, nil
            }
        }
    }
    log.Printf("Task did not complete within the given timeout")
    return nil, fmt.Errorf("task not completed")
}



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

// func (m *Mythic) SubscribeNewTasks(batchSize int, timeout int, callbackDisplayId int) ([]byte, error) {}

// func (m *Mythic) SubscribeNewTasksAndUpdates(batchSize int, timeout int, callbackDisplayId int) ([]byte, error) {}

// func (m *Mythic) SubscribeAllTasks(timeout int, callbackDisplayId int) ([]byte, error) {}

// func (m *Mythic) SubscribeAllTasksAndUpdates(timeout int, callbackDisplayId int) ([]byte, error) {}

// func (m *Mythic) AddMitreAttackToTask(timeout int, taskDisplayID int, MitreAttackNumbers []string ) ([]byte, error) {}





// # ######### File Browser Functions ###########

//func (m *Mythic) GetAllFileBrowser(host string, batchSize int) ([]byte, error) {}

//func (m *Mythic) SubscribeNewFileBrowser(host string, batchSize int, timeout int) ([]byte, error) {}

//func (m *Mythic) SubscribeAllFileBrowser(host string, timeout int, batchSize int) ([]byte, error) {}




// # ######### Command Functions ##############

// func (m *Mythic) GetAllCommandsForPayloadType(payloadTypeName string) ([]byte, error) {}



// # ######### Payload Functions ##############

/* func (m *Mythic) CreatePayload(payloadTypeName string, filename string, operatingSystem string, c2Profiles interface{}, commands []string, 
	buildParameters interface{}, description string, returnOnComplete bool, timeout int, includeAllCommands bool) ([]byte, error) {}

*/


// # ######### Task Output Functions ###########


func (m *Mythic) WaitForTaskOutput(taskDisplayID int, timeout *int) ([]byte, error) {
    var subscription TaskWaitForOutputSubscription

    variables := TaskWaitForOutputSubscriptionVariables{
        DisplayID: taskDisplayID,
    }
    variableMap := structToMap(variables)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10) // Set timeout to 10 seconds
    defer cancel() // Make sure to cancel the context when we're done to free resources

    events, err := m.GraphQLSubscription(ctx, &subscription, variableMap, *timeout)
    if err != nil {
        log.Printf("Error subscribing to task updates: %v", err)
        return nil, err
    }

    finalOutput := make([]byte, 0)

    for {
        select {
        case event := <-events:
            v, ok := event.(*TaskWaitForOutputSubscription)
            if !ok {
                return nil, fmt.Errorf("unexpected type: %T", event)
            }
            for _, taskStream := range v.TaskStream {
                for _, response := range taskStream.Responses {
                    outputBytes, err := base64.StdEncoding.DecodeString(response.ResponseText)
                    if err != nil {
                        return nil, fmt.Errorf("failed to decode base64 response text: %v", err)
                    }
                    finalOutput = append(finalOutput, outputBytes...)
                }
            }
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
            return finalOutput, nil
        }
    }
}


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
	
	variables := map[string]interface{}{
		"task_id": taskDisplayID,
	}

	var initial TaskIdFromDisplayID
	if err := m.GraphqlPost(&initial, variables, "query"); err != nil {
		return nil, err
	}

	subtaskIds := []int{}
	taskIdsToCheck := []int{}
	if len(initial.Task) > 0 {
		taskIdsToCheck = append(taskIdsToCheck, initial.Task[0].ID)
	}

	for len(taskIdsToCheck) > 0 {
		currentTaskId := taskIdsToCheck[len(taskIdsToCheck)-1]
		taskIdsToCheck = taskIdsToCheck[:len(taskIdsToCheck)-1]

		variables["task_id"] = currentTaskId  // update the 'task_id' in variables map
		
		var subtasks SubtaskList
		if err := m.GraphqlPost(&subtasks, variables, "query"); err != nil {
			return nil, err
		}

		for _, t := range subtasks.Task {
			taskIdsToCheck = append(taskIdsToCheck, t.ID)
			if fetchDisplayIDInstead {
				subtaskIds = append(subtaskIds, t.DisplayID)
			} else {
				subtaskIds = append(subtaskIds, t.ID)
			}
		}
	}

	return subtaskIds, nil
}


func (m *Mythic) GetAllTaskOutputByID(taskDisplayID int) ([]TaskOutputFragment, error) {
	var taskOutput TaskOutput
	
	variables := map[string]interface{}{
		"task_display_id": taskDisplayID,
	}
	err := m.GraphqlPost(&taskOutput, variables, "query")
	if err != nil {
		return nil, err
	}

	return taskOutput.Response, nil
}


// TODO:

// # ########## Operator Functions ##############

// # ########## File Functions ##############

// # ########## Operations Functions #############

// # ############ Process Functions ##############

// # ####### Analytic-based Functions ############

// # ####### Event Feed functions ############

// # ####### webhook ############

// # ####### C2 Functions #############

// # ####### Tag Functions ############





