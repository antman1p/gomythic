package gomythic

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

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

type MythicWebSocketHandler struct {
    Conn    *websocket.Conn
    timeout time.Duration
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


type MythicManager struct {
    instance *Mythic
}

func NewMythicManager() *MythicManager {
    return &MythicManager{}
}

func (mm *MythicManager) GetMythicInstance() *Mythic {
    if mm.instance == nil {
        mm.instance = &Mythic{} // create a new instance here
    }
    return mm.instance
}

func (mm *MythicManager) InvalidateMythicInstance() {
    mm.instance = nil // set the instance to nil
}

