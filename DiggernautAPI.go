package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}
var apikey string

type API struct {
	Projects []Project `json:"projects,omitempty"`
}
type Project struct {
	Id          int
	Name        string
	Description string
	Diggers     []Digger `json:"diggers,omitempty"`
}
type Digger struct {
	Id           int
	Name         string
	Url          string
	Config       string `json:"config,omitempty"`
	Status       string
	Schedulefrom time.Time `json:"schedule_from,omitempty"`
	Scheduleto   time.Time `json:"schedule_to,omitempty"`
	Bandwidth    string
	Calls        int
	Requests     int
	Sessions     []Session `json:"sessions,omitempty"`
}
type Session struct {
	DiggerId   int
	Id         int
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	State      string
	Runtime    int
	Bandwidth  int
	Requests   int
	Errors     int
	Data       interface{}
}

func SetUpAPIKey(key string) {
	apikey = key
}

func (a *API) GetProjects() error {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/", nil)
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &a.Projects)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(a.Projects)
	return nil
}
func (a *API) CreateProject(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://www.diggernaut.com/api/v1/projects/", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	p := Project{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		return err
	}
	a.Projects = append(a.Projects, p)
	return nil
}

func (p *Project) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.Id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *Project) Put(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.Id), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) Patch(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.Id), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}
	return nil
}
func (p *Project) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.Id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	return nil
}

func (p *Project) GetDiggers() error {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.Id)+"/diggers", nil)
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p.Diggers)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *Project) CreateDigger(params map[string]interface{}) error {
	params["project"] = p.Id
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://www.diggernaut.com/api/v1/diggers/", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	d := Digger{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return err
	}
	p.Diggers = append(p.Diggers, d)
	return nil
}

func (d *Digger) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.Id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *Digger) Put(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.Id), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}
	return nil
}

func (d *Digger) Patch(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.Id), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}
	return nil
}
func (d *Digger) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.Id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	return nil
}

func (d *Digger) GetSessions() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.Id)+"/sessions", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d.Sessions)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Session) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerId)+"/sessions/"+strconv.Itoa(s.Id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Session) GetData() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerId)+"/sessions/"+strconv.Itoa(s.Id)+"/data", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s.Data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}


