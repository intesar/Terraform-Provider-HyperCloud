package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"log"
	//"io"
	"io/ioutil"
	"net/http"
	//"net/url"
	//"reflect"
	//"strconv"
	//"time"
	"bytes"
	"strings"
	"time"
)

type Auth struct {
	URL string
	KEY string
	SECRET string
}

func newAuth(url string, key string, secret string) (auth *Auth)  {
	return &Auth{URL: url, KEY: key, SECRET: secret}
}

func (a Auth) getVM(id string) (*HCApiResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", a.URL + "/api/dockerservers/" + id, nil)
	req.SetBasicAuth(a.KEY, a.SECRET)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	//fmt.Printf("String: [%s]", string(body))

	if err != nil {
		panic(err.Error())
	}

	s, err := a.getResult(body)

	return s, err
}

func (a Auth) getResult(body []byte) (*HCApiResponse, error) {
	var s = HCApiResponse{}
	err := json.Unmarshal(body, &s)
	if(err != nil){
		fmt.Println("whoops:", err)
	}
	return &s, err
}

func (a Auth) getBlueprint(blueprintId string) (*HCApiResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", a.URL + "/api/blueprints/" + blueprintId, nil)
	req.SetBasicAuth(a.KEY, a.SECRET)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	//fmt.Printf("String: [%s]", string(body))

	if err != nil {
		panic(err.Error())
	}

	s, err := a.getResult(body)

	return s, err
}

func (a Auth) newBlueprintClient() *http.Client {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", a.URL + "/api/blueprints/", nil)
	req.SetBasicAuth(a.KEY, a.SECRET)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return client
}

func (a Auth) waitForTask(id string) (*HCApiResponse) {
	s, _ := a.getVM(id)
	fmt.Printf("\nResource is under processing with status [%s]", s.Results.DockerServerStatus)
	if ( strings.HasSuffix(s.Results.DockerServerStatus, "ING")) {
		time.Sleep(3000 * time.Millisecond)
		return a.waitForTask(id)
	}

	return s

}

func printDetails(s *HCApiResponse)  {
	fmt.Println("\n HyperCloud resource...")
	fmt.Printf("\n ID: [%s]", s.Results.ID)
	fmt.Printf("\n Name: [%s]", s.Results.Name)
	fmt.Printf("\n Description: [%s]", s.Results.Description)
	fmt.Printf("\n Status: [%s]", s.Results.DockerServerStatus)
	fmt.Printf("\n IP: [%s]", s.Results.HostOrIp)

}

func (a Auth) create(blueprintId string) *HCApiResponse {

	log.Printf("[HC-INFO] Creating Compute from Blueprint %s...", blueprintId)

	vm := ComputeRequest{}

	vm.Blueprint = blueprintId
	result, err := a.getBlueprint(blueprintId)
	if err != nil {
		panic(err.Error())
	}
	vm.CloudProvider = result.Results.CloudProvider.ID

	jsonValue, _ := json.Marshal(vm)


	client := &http.Client{}
	req, err := http.NewRequest("POST", a.URL + "/api/dockerservers/sdi", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	req.SetBasicAuth(a.KEY, a.SECRET)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	body1, err := ioutil.ReadAll(res.Body)

	s, err := a.getResult(body1)

	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	log.Printf("[HC-INFO] Create Response %s", string(body1))
	//fmt.Printf("String: [%s]", string(body1))
	return a.waitForTask(s.Results.ID)

}

func (a Auth) delete(id string) (*HCApiResponse, error) {
	log.Printf("[HC-INFO] Deleting Compute ID %s...", id)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", a.URL + "/api/dockerservers/ " + id + "?force=true", nil)
	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}
	req.SetBasicAuth(a.KEY, a.SECRET)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	body1, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	s, err := a.getResult(body1)
	if err != nil {
		log.Printf("[HC-ERROR] %s", err)
		panic(err.Error())
	}

	//fmt.Printf("String: [%s]", string(body1))
	log.Printf("[HC-INFO] Delete Response %s", string(body1))
	a.waitForTask(s.Results.ID)

	return s, nil

}

type ComputeRequest struct {
	Blueprint string `json:"blueprint"`
	CloudProvider string `json:"cloudProvider"`
	Cluster string `json:"cluster"`
}

type HCAPIMessage struct {
	MessageType string `json:"messageType"`
	MessageKey string `json:"messageKey"`
	MessageText  string `json:"messageText"`
}

type CloudProvider struct {
	ID string `json:"id"`
}

type HCBlueprint struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Description string `json:"description"`
	BlueprintType string `json:"blueprintType"`
	Yml string `json:"yml"`
	TotalStars int `json:"totalStars"`
	TotalRun int `json:"totalRun"`
	HostOrIp string `json: "hostOrIp"`
	DockerServerStatus string `json:"dockerServerStatus"`
	CloudProvider CloudProvider `json:"cloudProvider"`
	CreateDate int `json:"created"`
}

type HCApiResponse struct {
	Errors bool `json:"errors"`
	//Messages []HCAPIMessage `json: "messages"`
	//RequestID     string `json:"requestId"`
	//TotalElements int `json:"totalElements"`
	//TotalPages    int `json:"totalPages"`
	Results HCBlueprint `json:"results"`
}
