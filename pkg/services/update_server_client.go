package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	libUrl "net/url"
	"os"
	"time"
	"update-service/pkg/lib"
	"update-service/pkg/models"
)

type VIPNetUpdateServerClient struct {
	url      string
	login    string
	password string
	client   *http.Client
}

func NewVIPNetUpdateServerClient(serverUrl, login, password string, client *http.Client) *VIPNetUpdateServerClient {
	return &VIPNetUpdateServerClient{
		url:      serverUrl,
		login:    login,
		password: password,
		client:   client,
	}
}

func (inst *VIPNetUpdateServerClient) Login() error {
	if err := inst.ping(); err != nil {
		return err
	}

	data := libUrl.Values{}
	data.Set("login", inst.login)
	data.Set("password", inst.password)
	request, err := http.NewRequest("POST", fmt.Sprintf("%v/login", inst.url), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// Perform the login request
	response, err := inst.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("login error, updatesStatus code: %v", response.StatusCode)
	}

	defer response.Body.Close()

	// Decode the login response
	mBody := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&mBody)
	if err != nil {
		return err
	}

	// Check if login was successful
	if mBody["success"] == false {
		return fmt.Errorf("failed to login")
	}

	return nil
}

func (inst *VIPNetUpdateServerClient) UpdateList(pkgType lib.PackageType) ([]models.RrUpdates, error) {
	// Create a request to get the updates list
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/service/updates_list?type=%s", inst.url, pkgType), nil)
	if err != nil {
		return nil, err
	}

	// Make the request to get the update list
	response, err := inst.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed request update_list: %v", response.StatusCode)
	}

	// Parse the response to extract package information
	updateList := &models.UpdateListResponse{}
	err = json.NewDecoder(response.Body).Decode(updateList)
	if err != nil {
		return nil, err
	}

	return updateList.RrUpdates, nil
}

// Download downloads a package from the update server.
func (inst *VIPNetUpdateServerClient) Download(pkgType lib.PackageType, pkgInfo *models.RrUpdates, dir4Save string) (string, error) {
	contextWithTimeout, cancelFunc := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancelFunc()
	// Create a download request
	request, err := http.NewRequestWithContext(
		contextWithTimeout,
		http.MethodGet,
		fmt.Sprintf("%vservice/download", inst.url),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Add query parameters to the request
	q := request.URL.Query()
	q.Add("type", pkgType.String())
	q.Add("path", pkgInfo.Link)
	q.Add("name", pkgInfo.Name)
	request.URL.RawQuery = q.Encode()

	// Perform the download
	response, err := inst.client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("unexpected: %v", response.StatusCode)
	}

	// Save the downloaded file to cache
	cache, err := os.Create(fmt.Sprintf("%v/%v", dir4Save, pkgInfo.Name))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(cache, response.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/%v", dir4Save, pkgInfo.Name), nil
}

func (inst *VIPNetUpdateServerClient) ping() error {
	// Create request to fetch cookies
	request, err := http.NewRequest("GET", inst.url, nil)
	if err != nil {
		return err
	}

	// Perform the request to fetch cookies
	response, err := inst.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("update server responded with %d", response.StatusCode)
	}
	defer response.Body.Close()

	return nil
}
