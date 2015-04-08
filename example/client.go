package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint2/binary"
	"github.com/tendermint/tendermint2/rpc"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type Response struct {
	Status string
	Data   interface{}
	Error  string
}

//go:generate go-rpc-gen -interface Client -pkg core -dir core -type *ClientHTTP,*ClientJSON -exclude pipe.go -out-pkg rpc -out client_methods.go

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

func (c *ClientJSON) requestResponse(s rpc.JSONRPC) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	resp, err := http.Post(c.addr, "text/json", buf)
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

/*rpc-gen:define-interface Client
type Client interface {
	Address() string // returns the remote address
}
*/

func bytesToString(b []byte) (string, error) {
	return "0x" + hex.EncodeToString(b), nil
}

func intToString(b int) (string, error) {
	return strconv.Itoa(b), nil
}

func binaryWriter(args ...interface{}) ([]interface{}, error) {
	list := []interface{}{}
	for _, a := range args {
		buf, n, err := new(bytes.Buffer), new(int64), new(error)
		binary.WriteJSON(a, buf, n, err)
		if *err != nil {
			return nil, *err
		}
		list = append(list, buf.Bytes())

	}
	return list, nil
}

// for HTTP, we have a single function that converts args to values

func argsToURLValues(argNames []string, args ...interface{}) (url.Values, error) {
	values := make(url.Values)
	if len(argNames) == 0 {
		return values, nil
	}
	if len(argNames) != len(args) {
		return nil, fmt.Errorf("argNames and args have different lengths: %d, %d", len(argNames), len(args))
	}
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

/*rpc-gen:imports:
github.com/tendermint/tendermint2/binary
github.com/tendermint/tendermint2/rpc
net/http
io/ioutil
fmt
*/

// Template functions to be filled in

/*rpc-gen:template:*ClientJSON func (c *ClientJSON) {{name}}({{args.def}}) ({{response}}) {
	params, err := binaryWriter({{args.ident}})
	if err != nil{
		return nil, err
	}
	s := rpc.JSONRPC{
		JSONRPC: "2.0",
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
		Data   {{response.0}}
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

/*rpc-gen:template:*ClientHTTP func (c *ClientHTTP) {{name}}({{args.def}}) ({{response}}){
	values, err := argsToURLValues({{args.name}}, {{args.ident}})
	if err != nil{
		return nil, err
	}
	resp, err := http.PostForm(c.addr+{{lowername}}, values)
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
		Data   {{response.0}}
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
