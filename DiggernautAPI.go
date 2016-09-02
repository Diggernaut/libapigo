package DiggernautAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}
var apikey string

// API struct contains slice of projects.
type API struct {
	Projects []Project `json:"projects,omitempty"`
}

// Project cointains single project
type Project struct {
	ID          int
	Name        string
	Description string
	Diggers     []Digger `json:"diggers,omitempty"`
}

// Digger cointains single digger
type Digger struct {
	ID           int
	Name         string
	URL          string
	Config       string `json:"config,omitempty"`
	Status       string
	Schedulefrom time.Time `json:"schedule_from,omitempty"`
	Scheduleto   time.Time `json:"schedule_to,omitempty"`
	Bandwidth    string
	Calls        int
	Requests     int
	Sessions     []Session `json:"sessions,omitempty"`
}

// Session cointains single session
type Session struct {
	DiggerID   int
	ID         int
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	State      string
	Runtime    int
	Bandwidth  int
	Requests   int
	Errors     int
	Data       interface{}
}

// SetUpAPIKey it`s startpoint for using our API
// arg must be you API key
func SetUpAPIKey(key string) {
	apikey = key
}

// GetProjects returns list of projects linked with
// authenticated user account and push it in API.Projects slice
func (a *API) GetProjects() error {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/", nil)
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &a.Projects)
	if err != nil {
		return err
	}

	return nil
}

// CreateProject creates new project for authenticated user account
// and push it in API.Projects slice
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
	defer resp.Body.Close()
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

// Get returns project parameters and rewrite Project in API.Projects slice
func (p *Project) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return err
	}
	return nil
}

// Put updates project parameters
// and rewrite Project in API.Projects slice,
// all required fields will be updated with sent parameters.
func (p *Project) Put(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

// Patch updates project parameters partially
// and rewrite Project in API.Projects slice,
// only sent fields will be updated.
func (p *Project) Patch(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

// Delete deletes project
// Note! you must call API.GetProjects()
// to update Projects slice
func (p *Project) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}

	return nil
}

// GetDiggers returns list of diggers from specified project
// and push it in API.Projects[p.ID].Diggers slice
func (p *Project) GetDiggers() error {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID)+"/diggers", nil)
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p.Diggers)
	if err != nil {
		return err
	}
	return nil
}

// CreateDigger creates new digger for authenticated user account
// and push it in API.Projects[p.ID].Diggers slice
func (p *Project) CreateDigger(params map[string]interface{}) error {
	params["project"] = p.ID
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
	defer resp.Body.Close()
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

// Get gets parameters for digger
// and rewrite Digger in API.Projects[p.ID].Digger[d.ID] slice
func (d *Digger) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID), nil)
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
		return err
	}
	return nil
}

// Put updates digger parameters
// and rewrite Digger in API.Projects[p.ID].Digger[d.ID] slice,
// all required fields will be updated with sent parameters
func (d *Digger) Put(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
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

// Patch updates digger parameters partially
// and rewrite Digger in API.Projects[p.ID].Digger[d.ID] slice,
// only sent fields will be updated.
func (d *Digger) Patch(params map[string]interface{}) error {
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
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

// Delete deletes digger
// Note! you must call Project[p.ID].GetDiggers()
// to update Diggers slice
func (d *Digger) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	return nil
}

// GetSessions gets list of sessions for digger
// and push it in API.Projects[p.ID].Diggers[d.ID].Sessions slice
func (d *Digger) GetSessions() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID)+"/sessions", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d.Sessions)
	if err != nil {
		return err
	}
	return nil
}

// Get gets session parameters
// and rewrite it in API.Projects[p.ID].Diggers[d.ID].Sessions[s.ID] slice
func (s *Session) Get() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerID)+"/sessions/"+strconv.Itoa(s.ID), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return err
	}
	return nil
}

// GetData gets data scraped in given session
// and push it in API.Projects[p.ID].Diggers[d.ID].Sessions[s.ID].Data
func (s *Session) GetData() error {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerID)+"/sessions/"+strconv.Itoa(s.ID)+"/data", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+apikey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s.Data)
	if err != nil {
		return err
	}
	return nil
}
