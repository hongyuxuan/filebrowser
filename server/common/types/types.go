package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/minio/minio-go/v7"
)

type TraceIDKey struct{}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "[]"
	}
	return str, err
}

func (s *StringList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

/*********** map[string]string ***********/
type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "{}"
	}
	return str, err
}

func (s *StringMap) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

/*********** map[string]interface{} ***********/
type InterfaceMap map[string]interface{}

func (s InterfaceMap) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "{}"
	}
	return str, err
}

func (s *InterfaceMap) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

type GetDataReq struct {
	Tablename string `path:"tablename"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
	Search    string `form:"search,optional"`
	Filter    string `form:"filter,optional"`
	Range     string `form:"range,optional"`
	Preload   bool   `form:"preload,optional"`
	Sort      string `form:"sort,optional"`
}

type LoginRes struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
	RefreshAfter int64  `json:"refresh_after"`
}

type ListResult struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}

type S3Conn struct {
	S3Endpoint string        `json:"s3_endpoint"`
	S3Region   string        `json:"s3_region"`
	Client     *minio.Client `json:"client"`
}
