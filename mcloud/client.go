package mcloud

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HostURL - Default Mcloud URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := Client{
		HTTPClient: &http.Client{
			//Timeout: 600 * time.Second,
			Transport: tr},
		// Default Mcloud URL
		HostURL: HostURL,
		Auth: AuthStruct{
			Username: *username,
			Password: *password,
		},
	}

	if host != nil {
		c.HostURL = *host
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) waitForTaskToFinish(taskId int) error {
	client := c.HTTPClient

	var taskResponse TaskResponse

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/task/%s", strings.Trim(c.HostURL, "/"), strconv.Itoa(taskId)), nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", c.Token)

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		}

		err = json.Unmarshal(body, &taskResponse)
		//err = json.NewDecoder(resp.Body).Decode(sshKeyResponse)
		if err != nil {
			return err
		}

		debug(taskResponse.Task.Status)

		if taskResponse.Task.Status == "finished" {
			return nil
		} else if taskResponse.Task.Status == "failed" {
			return fmt.Errorf("mCloud task '%d' failed", taskId)
		} else if taskResponse.Task.Status != "running" {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}
