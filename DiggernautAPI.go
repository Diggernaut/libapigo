package DiggernautAPIS

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

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
func NewProject(a *API) *Project {
	return &Project{API: a}
}
func NewDigger(a *API) *Digger {
	return &Digger{API: a}
}
func NewSession(a *API) *Session {
	return &Session{api: a}
}

type Project struct {
	API         *API
	id          int
	Name        string
	Description string
}

func (p *Project) String() string {
	return "\n\nProject{" + "\n    API: " + p.API.Key + "\n    ID: " + strconv.Itoa(p.ID()) + "\n    Name: " + p.Name + "\n    Descrition: " + p.Description + "}"
}

func (p *Project) ID() int {
	return p.id
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

type Digger struct {
	API          *API
	id           int
	project      int
	Name         string
	URL          string
	Config       string `json:"config,omitempty"`
	status       string
	Schedulefrom timestamp.Timestamp `json:"schedule_from"`
	Scheduleto   timestamp.Timestamp `json:"schedule_to"`
	bandwidth    float64
	calls        int
	requests     int
	LastSession  *Session
}

func (d *Digger) UnmarshalJSON(data []byte) error {
	type digger struct {
		ID           int `json:"Id"`
		Project      int
		Name         string
		URL          string
		Config       string `json:"config,omitempty"`
		Status       string
		ScheduleFrom timestamp.Timestamp `json:"schedule_from"`
		ScheduleTo   timestamp.Timestamp `json:"schedule_to"`
		Bandwidth    float64
		Calls        int
		Requests     int
		LastSession  *Session `json:"last_session"`
	}
	dig := digger{}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.UseNumber()
	err := decoder.Decode(&dig)
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
	d.LastSession = dig.LastSession
	return nil

}
func (d *Digger) String() string {
	return "\n\nDigger{ " + "\n" + "ID: " + strconv.Itoa(d.ID()) + "\n" + "Status: " + d.Status() +
		"\n" + "Bandwidth: " + d.BandwidthString() + "\n" + "Requests: " + strconv.Itoa(d.Requests()) + "\n" + "Schedule from: " +
		d.Schedulefrom.String() + "\n" + "Schedule to: " + d.Scheduleto.String() + "\nLastSession{" + d.LastSession.StringToDigger() + "}}"
}
func (d *Digger) SetID(id int) {
	d.id = id
}
func (d *Digger) ID() int {
	return d.id
}
func (d *Digger) Status() string {
	return d.status
}
func (d *Digger) BandwidthString() string {
	return fmt.Sprintf("%.2f", d.bandwidth/1024/1024) + " Mb"
}
func (d *Digger) Bandwidth() float64 {
	return d.bandwidth
}
func (d *Digger) Calls() int {
	return d.calls
}
func (d *Digger) Requests() int {
	return d.requests
}

type Session struct {
	api        *API
	digger     int
	id         int
	startedAt  timestamp.Timestamp `json:"started_at,omitempty"`
	finishedAt timestamp.Timestamp `json:"finished_at,omitempty"`
	state      string
	runtime    int
	bandwidth  float64
	requests   int
	errors     int
	data       interface{}
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type session struct {
		ID         int `json:"Id"`
		Digger     int
		StartedAt  timestamp.Timestamp `json:"started_at,omitempty"`
		FinishedAt timestamp.Timestamp `json:"finished_at,omitempty"`
		State      string
		Runtime    int
		Bandwidth  float64
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
	s.digger = ses.Digger
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

func (s *Session) Digger() int {
	return s.digger
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
func (s *Session) RuntimeString() string {
	t, _ := time.ParseDuration(fmt.Sprint(s.Runtime()) + "s")
	return t.String()
}
func (s *Session) BandwidthString() string {
	return fmt.Sprintf("%.2f", s.bandwidth/1024/1024) + " Mb"
}
func (s *Session) Bandwidth() float64 {
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
func (s *Session) SetID(id int) {
	s.id = id
}
func (s *Session) String() string {
	return "\nDiggerID: " + strconv.Itoa(s.Digger()) + "\n" + "ID: " + strconv.Itoa(s.ID()) + "\n" + "State: " + s.State() + "\n" + "Runtime: " + s.RuntimeString() + "\n" + "Bandwidth: " + s.BandwidthString() + "\n" + "Requests: " + strconv.Itoa(s.Requests()) + "\n" + "Errors: " + strconv.Itoa(s.Errors()) + "\n" + "Started At: " + s.StartedAt().String() + "\n" + "Finished At: " + s.FinishedAt().String()
}

func (s *Session) StringToDigger() string {
	return "\n    DiggerID: " + strconv.Itoa(s.Digger()) + "\n" + "    ID: " + strconv.Itoa(s.ID()) + "\n" + "    State: " + s.State() + "\n" + "    Runtime: " + s.RuntimeString() + "\n" + "    Bandwidth: " + s.BandwidthString() + "\n" + "    Requests: " + strconv.Itoa(s.Requests()) + "\n" + "    Errors: " + strconv.Itoa(s.Errors()) + "\n" + "    Started At: " + s.StartedAt().String() + "\n" + "    Finished At: " + s.FinishedAt().String()
}
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
	return ret, nil
}
func (a *API) GetDiggers() ([]*Digger, error) {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/diggers/", nil)
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
	var ret []*Digger
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	for _, value := range ret {
		value.API = a
		value.LastSession.api = a
	}
	return ret, nil
}

func (a *API) GetSessions() ([]*Session, error) {
	req, _ := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/session/", nil)
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
	var ret []*Session
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	for _, value := range ret {
		value.api = a
	}
	return ret, nil
}

func (a *API) CreateProject(params map[string]interface{}) (*Project, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return &Project{}, err
	}
	req, err := http.NewRequest("POST", "https://www.diggernaut.com/api/v1/projects/", bytes.NewReader(payload))
	if err != nil {
		return &Project{}, err
	}
	req.Header.Add("Authorization", "Token "+a.Key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return &Project{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Project{}, errors.New(string(body[:]))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Project{}, err
	}
	p := NewProject(a)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return &Project{}, err
	}
	return p, nil
}

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
		value.LastSession.api = p.API
	}
	return ret, nil

}

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
	return d, err
}

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

	d.LastSession.api = d.API
	if err != nil {
		return d, err
	}
	return d, err
}

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

func (s *Session) Get() (*Session, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/sessions/"+strconv.Itoa(s.ID()), nil)
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

func (s *Session) GetData() (interface{}, error) {
	req, err := http.NewRequest("GET", "https://www.diggernaut.com/api/v1/sessions/"+strconv.Itoa(s.ID())+"/data", nil)
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
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s.data)
	if err != nil {
		return nil, err
	}
	return s.Data(), nil
}
