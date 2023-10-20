package vkapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const apiVersion = "5.100"

const requestURL string = "https://api.vk.com/method/%s?access_token=%s&v=" + apiVersion

func getRequestUrl(method, token string) string {
	return fmt.Sprintf(requestURL, method, token)
}

type logger func(object, text string)

var sleepLogger logger
var errorLogger logger

func SetSleepLogger(l logger) {
	sleepLogger = l
}

func SetErrorLogger(l logger) {
	errorLogger = l
}

func sleepLog(object, text string) {
	if sleepLogger != nil {
		sleepLogger(object, text)
	}
}

func errorLog(object, text string) {
	if errorLogger != nil {
		errorLogger(object, text)
	}
}

type Client struct {
	requestsLogger func(method, params string)
	token          string
	lastRequest    time.Time
	sleepTime      time.Duration

	workWithPool bool
	workers      *[]*Client
	workerIndex  int
}

type UserClientsPool struct {
	*UserClient
	clients []*Client
}

type ClientsPool struct {
	*Client
	clients []*Client
}

type ServiceClient struct {
	Client
}

type UserClient struct {
	Client
}

func NewPool(tokens []string) ClientsPool {
	if len(tokens) == 0 {
		panic("need at least one token")
	}

	pool := ClientsPool{
		Client:  nil,
		clients: make([]*Client, 0, 10),
	}

	for _, token := range tokens {
		c := NewClient(token)
		c.ActivatePool(&pool.clients)
		pool.clients = append(pool.clients, c)
	}

	pool.Client = pool.clients[0]
	return pool
}

func NewUserPool(tokens []string) UserClientsPool {
	if len(tokens) == 0 {
		panic("need at least one token")
	}

	pool := UserClientsPool{
		UserClient: nil,
		clients:    make([]*Client, 0, 10),
	}

	for _, token := range tokens {
		c := NewClient(token)
		c.ActivatePool(&pool.clients)
		pool.clients = append(pool.clients, c)
	}

	pool.UserClient = &UserClient{*pool.clients[0]}
	return pool
}

func (cp ClientsPool) SetRequestsLogger(a func(method, params string)) {
	for _, c := range cp.clients {
		c.requestsLogger = a
	}
}

func NewClient(token string) *Client {
	return &Client{token: token, sleepTime: time.Second / 20}
}

func NewUserClient(token string) *UserClient {
	uc := &UserClient{*NewClient(token)}
	uc.sleepTime = time.Second / 3
	return uc
}

func NewServiceClient(token string) *ServiceClient {
	sc := &ServiceClient{*NewClient(token)}
	sc.sleepTime = time.Second / 3
	return sc
}

func (c *Client) ActivatePool(workers *[]*Client) {
	c.workWithPool = true
	c.workers = workers
}

func (c *Client) DisablePool() {
	c.workWithPool = false
}

func (c *Client) request(method, params string) []byte {
	if c.requestsLogger != nil {
		c.requestsLogger(method, params)
	}

	for time.Since(c.lastRequest) < c.sleepTime {
		st := c.sleepTime - time.Since(c.lastRequest)
		if method != "groups.getLongPollServer" {
			sleepLog("vkapi", fmt.Sprintf("%d; %s; %s; %s", st, c.token, method, params))
		}
		time.Sleep(st)
	}
	c.lastRequest = time.Now()

	rURL := getRequestUrl(method, c.token)
	reader := strings.NewReader(params)
	r, err := http.Post(rURL, "application/x-www-form-urlencoded", reader)
	if err == nil {
		defer r.Body.Close()
	} else {
		log.Println("http post error", err)
		errorLog("http post", err.Error())
		return []byte{}
	}

	binAnswer, err := ioutil.ReadAll(r.Body)
	if strings.Contains(string(binAnswer), "\"error\"") {
		answer := string(binAnswer)
		if method == "appWidgets.update" ||
			strings.Contains(answer, "Can't send messages for users without permission") ||
			strings.Contains(answer, "user was kicked") ||
			strings.Contains(answer, "\"error_code\":15") ||
			strings.Contains(answer, "\"error_code\":917") {
			return nil
		}

		errorLog("vk answer", answer+"\nurl: "+rURL+"\nparams: "+params)
	}
	return binAnswer
}

var WorkersMutex = sync.Mutex{}

func (c *Client) Request(method, params string) []byte {
	if !c.workWithPool {
		return c.request(method, params)
	} else {
		WorkersMutex.Lock()
		c.workerIndex += 1
		if c.workerIndex == len(*c.workers) {
			c.workerIndex = 0
		}
		worker := (*c.workers)[c.workerIndex]
		WorkersMutex.Unlock()

		return worker.request(method, params)
	}
}
