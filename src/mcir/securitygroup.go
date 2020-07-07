package mcir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloud-barista/cb-tumblebug/src/common"
)

// 2020-04-13 https://github.com/cloud-barista/cb-spider/blob/master/cloud-control-manager/cloud-driver/interfaces/resources/SecurityHandler.go

type SpiderSecurityReqInfoWrapper struct { // Spider
	ConnectionName string
	ReqInfo        SpiderSecurityInfo
}

/*
type SpiderSecurityReqInfo struct { // Spider
	Name          string
	VPCName       string
	SecurityRules *[]SpiderSecurityRuleInfo
	//Direction     string // @todo used??
}
*/

type SpiderSecurityRuleInfo struct { // Spider
	FromPort   string `json:"fromPort"`
	ToPort     string `json:"toPort"`
	IPProtocol string `json:"ipProtocol"`
	Direction  string `json:"direction"`
}

type SpiderSecurityInfo struct { // Spider
	// Fields for request
	Name    string
	VPCName string

	// Fields for both request and response
	SecurityRules *[]SpiderSecurityRuleInfo

	// Fields for response
	IId          common.IID // {NameId, SystemId}
	VpcIID       common.IID // {NameId, SystemId}
	Direction    string     // @todo userd??
	KeyValueList []common.KeyValue
}

type TbSecurityGroupReq struct { // Tumblebug
	Name           string                    `json:"name"`
	ConnectionName string                    `json:"connectionName"`
	VNetId         string                    `json:"vNetId"`
	Description    string                    `json:"description"`
	FirewallRules  *[]SpiderSecurityRuleInfo `json:"firewallRules"`
}

type TbSecurityGroupInfo struct { // Tumblebug
	Id                   string                    `json:"id"`
	Name                 string                    `json:"name"`
	ConnectionName       string                    `json:"connectionName"`
	VNetId               string                    `json:"vNetId"`
	Description          string                    `json:"description"`
	FirewallRules        *[]SpiderSecurityRuleInfo `json:"firewallRules"`
	CspSecurityGroupId   string                    `json:"cspSecurityGroupId"`
	CspSecurityGroupName string                    `json:"cspSecurityGroupName"`
	KeyValueList         []common.KeyValue         `json:"keyValueList"`

	// Disabled for now
	//ResourceGroupName  string `json:"resourceGroupName"`
}

func CreateSecurityGroup(nsId string, u *TbSecurityGroupReq) (TbSecurityGroupInfo, error) {
	check, _ := CheckResource(nsId, "securityGroup", u.Name)

	if check {
		temp := TbSecurityGroupInfo{}
		err := fmt.Errorf("The securityGroup " + u.Name + " already exists.")
		//return temp, http.StatusConflict, nil, err
		return temp, err
	}

	//url := common.SPIDER_URL + "/securitygroup?connection_name=" + u.ConnectionName
	url := common.SPIDER_URL + "/securitygroup"

	method := "POST"

	//payload := strings.NewReader("{ \"Name\": \"" + u.CspSecurityGroupName + "\"}")
	tempReq := SpiderSecurityReqInfoWrapper{}
	tempReq.ConnectionName = u.ConnectionName
	tempReq.ReqInfo.Name = u.Name
	tempReq.ReqInfo.VPCName = u.VNetId
	tempReq.ReqInfo.SecurityRules = u.FirewallRules

	payload, _ := json.Marshal(tempReq)
	fmt.Println("payload: " + string(payload)) // for debug

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		common.CBLog.Error(err)
		content := TbSecurityGroupInfo{}
		//return content, res.StatusCode, nil, err
		return content, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	if err != nil {
		common.CBLog.Error(err)
		content := TbSecurityGroupInfo{}
		//return content, res.StatusCode, body, err
		return content, err
	}

	fmt.Println("HTTP Status code " + strconv.Itoa(res.StatusCode))
	switch {
	case res.StatusCode >= 400 || res.StatusCode < 200:
		err := fmt.Errorf(string(body))
		common.CBLog.Error(err)
		content := TbSecurityGroupInfo{}
		//return content, res.StatusCode, body, err
		return content, err
	}

	temp := SpiderSecurityInfo{}
	err2 := json.Unmarshal(body, &temp)
	if err2 != nil {
		fmt.Println("whoops:", err2)
	}

	content := TbSecurityGroupInfo{}
	//content.Id = common.GenUuid()
	content.Id = common.GenId(u.Name)
	content.Name = u.Name
	content.ConnectionName = u.ConnectionName
	content.VNetId = temp.VpcIID.NameId
	content.CspSecurityGroupId = temp.IId.SystemId
	content.CspSecurityGroupName = temp.IId.NameId
	content.Description = u.Description
	content.FirewallRules = temp.SecurityRules
	content.KeyValueList = temp.KeyValueList

	// cb-store
	fmt.Println("=========================== PUT CreateSecurityGroup")
	Key := common.GenResourceKey(nsId, "securityGroup", content.Id)
	Val, _ := json.Marshal(content)
	err = common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		//return content, res.StatusCode, body, err
		return content, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")
	//return content, res.StatusCode, body, nil
	return content, nil
}
