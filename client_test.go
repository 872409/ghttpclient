package x_http_client

import (
	"fmt"
	logger "log"
	"os"
	"strings"
	"testing"

)

type ResponseJSON struct {
	Code int         `json:"code,int"`
	Msg  interface{} `json:"msg,string"`
}

func newClient() *Client {
	client, err := New("127.0.0.1:5001", "1234512345", "")
	client.Config.LogLevel = Debug
	client.Config.Logger = logger.New(os.Stdout, "", logger.LstdFlags)
	fmt.Println("newClient:", err)
	return client
}

func TestNew(t *testing.T) {
	client := newClient()
	headers := make(map[string]string)
	params := make(map[string]interface{})
	data := strings.NewReader("{\n  \"biz_code\": \"\",\n  \"biz_no\": \"\",\n  \"track_no\": \"\",\n  \"app_id\": 0,\n  \"request_no\": \"\",\n  \"request_ip\": \"\",\n  \"request_uid\": 0,\n  \"request_role\": \"\",\n  \"out_order_no\": \"\",\n  \"acct_uid\": 0,\n  \"title\": \"0\",\n  \"amount\": 0,\n  \"fee\": 0,\n  \"remark\": \"remark\"\n}")
	json := &ResponseJSON{}
	resp, err := client.Conn.DoJSONResponse("POST", "/account/trade/withdraw/pre", params, headers, data, json)
	// body := ""
	// if err == nil {
	// 	out, e := ioutil.ReadAll(resp.Body)
	// 	body = string(out)
	// 	err = e
	// }
	// log.Infof("body:%v", body)
	// log.Infof("Do:%v,%v,%v", resp.Headers, body, err)

	// jsonUnmarshal(resp.Body, json)
	fmt.Printf("body:%v", json)
	fmt.Printf("body:%v", json.Msg.(string))
	fmt.Printf("Do:%v,%v", resp, err)

}
