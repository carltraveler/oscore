package common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type ReqParam struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func ParseGetParamByParamName(c *gin.Context, paramNames ...string) ([]string, error) {
	res := make([]string, len(paramNames))
	for k, paramName := range paramNames {
		res[k] = c.Param(paramName)
	}
	return res, nil
}

func ParseGetParam(c *gin.Context, paramStruct interface{}) error {
	tags := getTagName(paramStruct)
	res := make(map[string]interface{})
	for _, paramName := range tags {
		res[paramName] = c.Param(paramName)
		if res[paramName] == "" {
			return fmt.Errorf("param: %s is nil", paramName)
		}
	}
	bs, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, paramStruct)
}

func getTagName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

func ParsePostParam(c *gin.Context, paramStruct interface{}) error {
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	rp := &ReqParam{}
	err = json.Unmarshal(paramsBs, rp)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(rp.Params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, paramStruct)
	if err != nil {
		return err
	}
	return nil
}
