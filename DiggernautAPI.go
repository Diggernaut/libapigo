package DiggernautAPIS

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Diggernaut/timestamp"
)

var client = &http.Client{}

type API struct {
	Key string
}

func (a API) String() string {
	return "Key: " + a.Key

}
func New(key string) *API {
	return &API{Key: key}
}
func NewProject(a *API) Project {
	return Project{API: a}
}
func NewDigger(a *API) Digger {
	return Digger{API: a}
}
func NewSession(a *API) Session {
	return Session{api: a}
}

// Project cointains single project
type Project struct {
	API         *API
	id          int
	Name        string
	Description string
}

func (p *Project) UnmarshalJSON(data []byte) error {
	type project struct {
		ID          int
		Name        string
		Description string
	}
	proj := project{}
	err := json.Unmarshal(data, &proj)
	if err != nil {
		return err
	}
	p.id = proj.ID
	p.Name = proj.Name
	p.Description = proj.Description
	return nil
}

func (p *Project) ID() int {
	return p.id
}

// Digger cointains single digger
type Digger struct {
	API          *API
	id           int
	Name         string
	URL          string
	Config       string `json:"config,omitempty"`
	status       string
	Schedulefrom timestamp.Timestamp `json:"schedule_from"`
	Scheduleto   timestamp.Timestamp `json:"schedule_to"`
	bandwidth    float64
	calls        int
	requests     int
}

func (d *Digger) UnmarshalJSON(data []byte) error {
	type digger struct {
		ID           int `json:"Id"`
		Name         string
		URL          string
		Config       string `json:"config,omitempty"`
		Status       string
		ScheduleFrom timestamp.Timestamp `json:"schedule_from"`
		ScheduleTo   timestamp.Timestamp `json:"schedule_to"`
		Bandwidth    float64
		Calls        int
		Requests     int
	}
	dig := digger{}
	err := json.Unmarshal(data, &dig)
	if err != nil {

		return err
	}
	d.id = dig.ID
	d.Name = dig.Name
	d.URL = dig.URL
	d.Config = dig.Config
	d.status = dig.Status
	d.Schedulefrom = dig.ScheduleFrom
	d.Scheduleto = dig.ScheduleTo
	d.bandwidth = dig.Bandwidth
	d.calls = dig.Calls
	d.requests = dig.Requests
	fmt.Println(d)
	return nil
}
func (d *Digger) ID() int {
	return d.id
}
func (d *Digger) Status() string {
	return d.status
}
func (d *Digger) Bandwdth() string {
	return fmt.Sprintf("%.0f\n", d.bandwidth)
}
func (d *Digger) Calls() int {
	return d.calls
}
func (d *Digger) Requests() int {
	return d.requests
}

// Session cointains single session
type Session struct {
	api        *API
	diggerID   int
	id         int
	startedAt  timestamp.Timestamp `json:"started_at,omitempty"`
	finishedAt timestamp.Timestamp `json:"finished_at,omitempty"`
	state      string
	runtime    int
	bandwidth  int64
	requests   int
	errors     int
	data       interface{}
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type session struct {
		ID         int `json:"Id"`
		DiggerID   int
		StartedAt  timestamp.Timestamp `json:"started_at,omitempty"`
		FinishedAt timestamp.Timestamp `json:"finished_at,omitempty"`
		State      string
		Runtime    int
		Bandwidth  int64
		Requests   int
		Errors     int
		Data       interface{}
	}
	ses := session{}
	err := json.Unmarshal(data, &ses)
	if err != nil {
		return err
	}
	s.id = ses.ID
	s.diggerID = ses.DiggerID
	s.startedAt = ses.StartedAt
	s.finishedAt = ses.FinishedAt
	s.state = ses.State
	s.runtime = ses.Runtime
	s.bandwidth = ses.Bandwidth
	s.requests = ses.Requests
	s.errors = ses.Errors
	s.data = ses.Data
	return nil
}

func (s *Session) DiggerID() int {
	return s.diggerID
}
func (s *Session) ID() int {
	return s.id
}
func (s *Session) State() string {
	return s.state
}
func (s *Session) Runtime() int {
	return s.runtime
}
func (s *Session) Bandwidth() int64 {
	return s.bandwidth
}
func (s *Session) Requests() int {
	return s.requests
}
func (s *Session) Errors() int {
	return s.errors
}
func (s *Session) Data() interface{} {
	return s.data
}
func (s *Session) StartedAt() timestamp.Timestamp {
	return s.startedAt
}
func (s *Session) FinishedAt() timestamp.Timestamp {
	return s.finishedAt
}
func (s *Session) API() *API {
	return s.api
}

// GetProjects returns list of projects linked with
// authenticated user account and push it in APIS.Projects slice
func (a *API) GetProjects() ([]*Project, error) {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/", nil)
	req.Header.Add("Authorization", "Token "+a.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var ret []*Project
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	for _, value := range ret {
		value.API = a
	}
	fmt.Printf("%+v\n", ret[0])
	return ret, nil
}

// CreateProject creates new project for authenticated user account
// and push it in APIS.Projects slice
func (a *API) CreateProject(params map[string]interface{}) (Project, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return Project{}, err
	}
	req, err := http.NewRequest("POST", "https://www.diggernaut.com/api/v1/projects/", bytes.NewReader(payload))
	if err != nil {
		return Project{}, err
	}
	req.Header.Add("Authorization", "Token "+a.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return Project{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return Project{}, errors.New(string(body[:]))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Project{}, err
	}
	p := NewProject(a)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return Project{}, err
	}
	return p, nil
}

// Get returns project parameters and rewrite Project
func (p Project) Get() (Project, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID()), nil)
	if err != nil {
		return p, err
	}
	req.Header.Add("Authorization", "Token "+p.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return p, errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Put updates project parameters
// and rewrite Project,
// all required fields will be updated with sent parameters.
func (p Project) Put() (Project, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return p, err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID()), bytes.NewReader(payload))
	if err != nil {
		return p, err
	}
	req.Header.Add("Authorization", "Token "+p.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return p, errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(body, p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Patch updates project parameters partially
// and rewrite Project,
// only sent fields will be updated.
func (p Project) Patch() (Project, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return p, err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID()), bytes.NewReader(payload))
	if err != nil {
		return p, err
	}
	req.Header.Add("Authorization", "Token "+p.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return p, errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Delete deletes project
func (p Project) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID()), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+p.API.Key)
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
// and push it in Project.Diggers slice
func (p *Project) GetDiggers() ([]*Digger, error) {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/projects/"+strconv.Itoa(p.ID())+"/diggers", nil)
	req.Header.Add("Authorization", "Token "+p.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body[:]))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var ret []*Digger
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	for _, value := range ret {
		value.API = p.API
	}
	return ret, nil

}

// CreateDigger creates new digger for authenticated user account
// and push it in Projects.Diggers slice
func (p *Project) CreateDigger(params map[string]interface{}) (*Digger, error) {
	params["project"] = p.ID
	payload, err := json.Marshal(params)
	if err != nil {
		return &Digger{}, err
	}
	req, err := http.NewRequest("POST", "https://www.diggernaut.com/api/v1/diggers/", bytes.NewReader(payload))
	if err != nil {
		return &Digger{}, err
	}
	req.Header.Add("Authorization", "Token "+p.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return &Digger{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Digger{}, errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Digger{}, err
	}
	d := NewDigger(p.API)
	err = json.Unmarshal(body, &d)
	if err != nil {
		return &Digger{}, err
	}
	return &d, err
}

// Get gets parameters for digger
// and rewrite Digger
func (d *Digger) Get() (*Digger, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID()), nil)
	if err != nil {
		return d, err
	}
	req.Header.Add("Authorization", "Token "+d.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return d, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return d, errors.New(string(body[:]))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d)
	if err != nil {
		return d, err
	}
	return d, err
}

// Put updates digger parameters
// and rewrite Digger,
// all required fields will be updated with sent parameters
func (d *Digger) Put() (*Digger, error) {
	payload, err := json.Marshal(d)
	if err != nil {
		return d, err
	}
	req, err := http.NewRequest("PUT", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID()), bytes.NewReader(payload))
	if err != nil {
		return d, err
	}
	req.Header.Add("Authorization", "Token "+d.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return d, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return d, errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(body, d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// Patch updates digger parameters partially
// and rewrite Digger,
// only sent fields will be updated.
func (d *Digger) Patch() (*Digger, error) {
	payload, err := json.Marshal(d)
	if err != nil {
		return d, err
	}
	req, err := http.NewRequest("PATCH", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID()), bytes.NewReader(payload))
	if err != nil {
		return d, err
	}
	req.Header.Add("Authorization", "Token "+d.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return d, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return d, errors.New(string(body[:]))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// Delete deletes digger
func (d *Digger) Delete() error {
	req, err := http.NewRequest("DELETE", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID()), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Token "+d.API.Key)
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
	d = &Digger{}
	return nil
}

// GetSessions gets list of sessions for digger
// and push it in Diggers.Sessions slice
func (d *Digger) GetSessions() ([]*Session, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(d.ID())+"/sessions", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+d.API.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var ret []*Session
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	for _, value := range ret {
		value.api = d.API
	}
	return ret, nil
}

// Get gets session parameters
// and rewrite Session
func (s *Session) Get() (*Session, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerID())+"/sessions/"+strconv.Itoa(s.ID()), nil)
	if err != nil {
		return s, err
	}
	req.Header.Add("Authorization", "Token "+s.api.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return s, errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetData gets data scraped in given session
// and push it in Session.Data
func (s *Session) GetData() (interface{}, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/"+strconv.Itoa(s.DiggerID())+"/sessions/"+strconv.Itoa(s.ID())+"/data", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+s.api.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body[:]))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s.data)
	if err != nil {
		return nil, err
	}
	return s.Data(), nil
}
