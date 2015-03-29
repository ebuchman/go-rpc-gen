package rpc

import (
	"bytes"
	"fmt"
	"github.com/tendermint/tendermint/binary"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Response struct {
	Status string
	Data   interface{}
	Error  string
}

//go:generate rpc-gen -interface Client -type ClientJSON,ClientHTTP -pkg core -excludes pipe.go -out .

type ClientJSON struct {
	addr string
}

type ClientHTTP struct {
	addr string
}

func NewClient(addr, typ string) Client {
	switch typ {
	case "HTTP":
		return &ClientHTTP{addr}
	case "JSONRPC":
		return &ClientJSON{addr}
	}
	return nil
}

func (c *ClientJSON) Call(method string, args ...interface{}) (*Response, error) {
	return nil, nil
}

func (c *ClientHTTP) Call(method string, args ...interface{}) (*Response, error) {
	fw, ok := funcMap[method]
	if !ok {
		return nil, fmt.Errorf("No known function method %s", method)
	}
	if len(args) != len(fw.args) {
		return nil, fmt.Errorf("Not enough arguments. Got %d, expected %d for method %s", len(args), len(fw.args), method)
	}
	values, err := argsToURLValues(fw.argNames, args...)
	if err != nil {
		return nil, err
	}
	return c.RequestResponse(method, values)
}

func argsToJson(args ...interface{}) ([][]string, error) {
	l := len(args)
	jsons := make([][]string, l)
	n, err := new(int64), new(error)
	for i, a := range args {
		//if its a slice, we serliaze separately and pack into a slice of strings
		// otherwise its a slice of length 1
		if v := reflect.ValueOf(a); v.Kind() == reflect.Slice {
			slice := []string{}
			for j := 0; j < v.Len(); j++ {
				buf := new(bytes.Buffer)
				binary.WriteJSON(v.Index(j).Interface(), buf, n, err)
				if *err != nil {
					return nil, *err
				}
				slice[j] = string(buf.Bytes())
			}
			jsons[i] = slice
		} else {
			buf := new(bytes.Buffer)
			binary.WriteJSON(a, buf, n, err)
			if *err != nil {
				return nil, *err
			}
			jsons[i] = []string{string(buf.Bytes())}
		}
	}
	return jsons, nil
}

func (c *ClientHTTP) RequestResponse(method string, values url.Values) (*Response, error) {
	resp, err := http.PostForm(c.addr+method, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	status := new(Response)
	//status.Data = ResponseStatus{}
	fmt.Println(string(body))
	binary.ReadJSON(status, body, &err)
	if err != nil {
		return nil, err
	}
	fmt.Println(status.Data)
	return status, nil
}

func (c *ClientJSON) requestResponse(s rpc.JsonRpc) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	resp, err = http.Post(c.addr, "text/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

/*
	What follows is used by `rpc-gen` when `go generate` is called
	to populate the rpc client methods
*/

// first we define the base interface, which rpc-gen will further populate with generated methods

//rpc-gen:define-interface Client
/*
type Client interface {
	Address() string // returns the remote address
}
*/

// A list of functions for encoding data to json strings
// first we declare that **packargsjson** should be replaced by a []string of **args**
// then we declare how to pack args of different types into the []string

//rpc-gen:define-set:packargsjson **args** []string

//rpc-gen:packargsjson []byte
func bytesToString(b []byte) (string, error) {
	return "0x" + hex.EncodeToString(b), nil
}

//rpc-gen:packargsjson int
func intToString(b int) (string, error) {
	return strconv.Itoa(b), nil
}

//rpc-gen:packargsjson default
func binaryWriter(b interface{}) (string, error) {
	buf, n, err := new(bytes.Buffer), new(int64), new(error)
	binary.WriteJSON(b, buf, n, err)
	return string(buf.Bytes()), *err
}

// for HTTP, we have a single function that converts args to values

//rpc-gen:define-func:packargshttp **args** url.Values
func argsToURLValues(argNames []string, args ...interface{}) (url.Values, error) {
	values := make(url.Values)
	slice, err := argsToJson(args...)
	if err != nil {
		return nil, err
	}
	for i, name := range argNames {
		s := slice[i]
		values.Set(name, s[0])
		for i := 1; i < len(s); i++ {
			values.Add(name, s[i])

		}
	}
	return values, nil
}

// Template functions to be filled in

/*rpc-gen:template:*ClientJSON func (c *ClientJSON) {{name}}({{args}}) ({{response}}, error) {
	params, err := {{packargsjson}}
	if err != nil{
		return nil, err
	}
	s := rpc.JsonRpc{
		JsonRpc: "2.0",
		Method:  {{lowername}},
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil{
		return nil, err
	}
	var status struct {
		Status string
		Data   {{response}}
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != ""{
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}*/

/*rpc-gen:template:*ClientHTTP func (c *ClientHTTP) {{name}}({{args}}) ({{response}}, error){
	values, err := {{packargshttp(argNames)}}
	if err != nil{
		return nil, err
	}
	resp, err := http.PostForm(requestAddr+{{lowername}}, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   {{response}}
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != ""{
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data.Account
}*/
