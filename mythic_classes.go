package Mythic_Go_Scripting

import "encoding/json"

type Mythic struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	APIToken           string `json:"apitoken"`
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	ServerIP           string `json:"server_ip"`
	ServerPort         int    `json:"server_port"`
	SSL                bool   `json:"ssl"`
	HTTP               string `json:"-"`
	WS                 string `json:"-"`
	GlobalTimeout      int    `json:"global_timeout"`
	ScriptingVersion   string `json:"scripting_version"`
	CurrentOperationID int    `json:"current_operation_id"`
	Schema             string `json:"schema"`
}

/* CreateTaskMutation = `
		mutation createTasking(
			$callback_id: Int!,
			$command: String!,
			$params: String!,
			$token_id: Int,
			$tasking_location: String,
			$original_params: String,
			$parameter_group_name: String
		) {
			createTask(
				callback_id: $callback_id,
				command: $command,
				params: $params,
				token_id: $token_id,
				tasking_location: $tasking_location,
				original_params: $original_params,
				parameter_group_name: $parameter_group_name
			) {
				status
				id
				display_id
				error
			}
		}
	` */

type CreateTaskMutation struct {
	CreateTask struct {
		Status    MythicStatus `graphql:"status"`
		ID        int `graphql:"id"`
		DisplayID int `graphql:"display_id"`
		Error     string `graphql:"error"`
	} `graphql:"createTask(callback_id: $callback_id, command: $command, params: $params, token_id: $token_id, tasking_location: $tasking_location, original_params: $original_params, parameter_group_name: $parameter_group_name)"`
}


// Defining struct for Callback
type Callback struct {
	Architecture    string          `graphql:"architecture"`
	Description     string          `graphql:"description"`
	Domain          string          `graphql:"domain"`
	ExternalIP      string          `graphql:"external_ip"`
	Host            string          `graphql:"host"`
	ID              int             `graphql:"id"`
	DisplayID       int             `graphql:"display_id"`
	IntegrityLevel  int              `graphql:"integrity_level"`
	IP              string          `graphql:"ip"`
	ExtraInfo       string          `graphql:"extra_info"`
	SleepInfo       string          `graphql:"sleep_info"`
	PID             int             `graphql:"pid"`
	OS              string          `graphql:"os"`
	User            string          `graphql:"user"`
	AgentCallbackID string          `graphql:"agent_callback_id"`
	OperationID     int             `graphql:"operation_id"`
	ProcessName     string          `graphql:"process_name"`
	Payload         CallbackPayload `graphql:"payload"`
}

// Defining struct for Payload
type CallbackPayload struct {
	OS          string             `graphql:"os"`
	PayloadType CallbackPayloadType `graphql:"payloadtype"`
	Description string             `graphql:"description"`
	UUID        string             `graphql:"uuid"`
}

// Defining struct for PayloadType
type CallbackPayloadType struct {
	Name string `graphql:"name"`
}

type CallbackQuery struct {
	Callback []Callback `graphql:"callback(order_by: {id: asc})"`
}

type ActiveCallbackQuery struct {
	Callback []Callback `graphql:"callback(where: {active: {_eq: true}}, order_by: {id: asc})"`
}

/* TaskFragment = `
		fragment task_fragment on task {
			callback {
				id
				display_id
			}
			id
			display_id
			operator {
				username
			}
			status
			completed
			original_params
			display_params
			timestamp
			command_name
			tasks {
				id
			}
			token {
				token_id
			}
		}
	` */

type TaskFragment struct {
	Callback struct {
		ID         int    `graphql:"id"`
		DisplayID  int    `graphql:"display_id"`
	} `graphql:"callback"`
	ID             int    `graphql:"id"`
	DisplayID      int    `graphql:"display_id"`
	Operator       struct {
		Username string   `graphql:"username"`
	} `graphql:"operator"`
	Status         MythicStatus     `graphql:"status"`
	Completed      bool   `graphql:"completed"`
	OriginalParams string `graphql:"original_params"`
	DisplayParams  string `graphql:"display_params"`
	Timestamp      string `graphql:"timestamp"`
	CommandName    string `graphql:"command_name"`
	Tasks          []struct {
		ID int `graphql:"id"`
	} `graphql:"tasks"`
	Token struct {
		TokenID string `graphql:"token_id"`
	} `graphql:"token"`
}


/* TaskOutputFragment = `
		fragment task_output_fragment on response {
			id
			timestamp
			response_text
			task {
				id
				display_id
				status
				completed
				agent_task_id
				command_name
			}
		}
	` */
	
type TaskOutputFragment struct {
	ID            int    `graphql:"id"`
	Timestamp     string `graphql:"timestamp"`
	ResponseText  string `graphql:"response_text"`
	Task struct {
		ID		    int                 `graphql:"id"`
		DisplayID   int		            `graphql:"display_id"`
		Status	    MythicStatus	    `graphql:"status"`
		Completed   bool		        `graphql:"completed"`
		AgentTaskID	int	                `graphql:"agent_task_id"`
		CommandName	string	            `graphql:"command_name"`
	} `graphql:"task"`
}

type TaskQuery struct {
	Task []TaskFragment `graphql:"task(order_by: {id: desc})"`
}

type TaskQueryWithCallback struct {
	Task []TaskFragment `graphql:"task(where: {callback: {display_id: {_eq: $callbackDisplayID}}}, order_by: {id: asc})"`
}

type TaskWaitForStatusSubscription struct {
    TaskStream []TaskFragment `graphql:"task_stream(cursor: {initial_value: {timestamp: \"1970-01-01\"}}, batch_size: 1, where: {display_id: {_eq: $DisplayID}})"`
}


type TaskWaitForStatusSubscriptionVariables struct {
    DisplayID int `graphql:"display_id"`
}

	/* CreateAPITokenMutation = `
		mutation createAPITokenMutation {
			createAPIToken(token_type: "User") {
				id
				token_value
				status
				error
				operator_id
			}
		}
	`

	GetAPITokensQuery = `
		query GetAPITokens {
			apitokens(where: {active: {_eq: true}}) {
				token_value
				active
				id
			}
		}
	` */


type CreateAPITokenMutation struct {
	CreateAPIToken CreateAPIToken `graphql:"createAPIToken(token_type: \"User\")"`
}

type CreateAPIToken struct {
	ID         int    `graphql:"id"`
	TokenValue string `graphql:"token_value"`
	Status     MythicStatus     `graphql:"status"`
	Error      string `graphql:"error"`
	OperatorID int    `graphql:"operator_id"`
}








func (m *Mythic) String() string {
	data, _ := json.MarshalIndent(m, "", "    ")
	return string(data)
}

type MythicStatus string

const (
	Error         MythicStatus = "error"
	Completed     MythicStatus = "completed"
	Processed     MythicStatus = "processed"
	Processing    MythicStatus = "processing"
	Preprocessing MythicStatus = "preprocessing"
	Delegating    MythicStatus = "delegating"
	Submitted     MythicStatus = "submitted"
)

func (s MythicStatus) String() string {
	return string(s)
}

func (s MythicStatus) Equals(obj interface{}) bool {
	if status, ok := obj.(string); ok {
		return s.String() == status
	} else if status, ok := obj.(MythicStatus); ok {
		return s.String() == status.String()
	}
	return false
}

func (s MythicStatus) GreaterThanOrEqual(obj interface{}) bool {
	targetObj := ""
	if status, ok := obj.(string); ok {
		targetObj = status
	} else if status, ok := obj.(MythicStatus); ok {
		targetObj = status.String()
	}
	if targetObj == "" {
		return false
	}
	if targetObj == "delegating" {
		targetObj = "delegating"
	} else if targetObj == "error" {
		targetObj = "error"
	}
	selfObj := s.String()
	if selfObj == "delegating" {
		selfObj = "delegating"
	}
	if selfObj == "error" {
		return true
	} else if selfObj == "completed" {
		return true
	}
	enumMapping := map[string]int{
		"preprocessing": 0,
		"submitted":     1,
		"delegating":    2,
		"processing":    3,
		"processed":     4,
		"completed":     5,
		"error":         6,
	}
	if _, ok := enumMapping[targetObj]; !ok {
		panic("Can't compare status of type: " + targetObj)
	} else if _, ok := enumMapping[selfObj]; !ok {
		panic("Can't compare status of type: " + selfObj)
	}
	return enumMapping[selfObj] >= enumMapping[targetObj]
}
