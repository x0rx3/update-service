package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"regexp"
	"time"
	"update-service/internal/model"
	"update-service/internal/utils"
)

type VIPNetIDSClient struct {
	client http.Client
}

func NewVIPNetIDSClient(client *http.Client) *VIPNetIDSClient {
	return &VIPNetIDSClient{
		client: *client,
	}
}

func (inst *VIPNetIDSClient) Login(idsUrl, login, password string) error {
	if err := inst.ping(idsUrl); err != nil {
		return err
	}

	data := url.Values{}
	data.Set("j_username", login)
	data.Set("j_password", password)

	// Create the login request
	req, err := http.NewRequest("POST", fmt.Sprintf("%vlogin", idsUrl), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; UTF-8")

	// Send the login request and check the response
	resp, err := inst.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the login was successful
	if resp.StatusCode != 200 {
		return fmt.Errorf("Login failed: %v", resp.StatusCode)
	}

	mBody := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&mBody)
	if err != nil {
		return err
	}

	if mBody["authorized"] == false {
		return fmt.Errorf("failed to login")
	}

	return nil
}

// SoftVersion retrieves the software version from the server's footer page.
func (inst *VIPNetIDSClient) SoftVersion(idsUrl string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%vfooter", idsUrl), nil)
	if err != nil {
		return "", err
	}

	resp, err := inst.client.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("footer request failed: unexpected status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	mBody := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&mBody)
	if err != nil {
		return "", err
	}
	version := regexp.MustCompile(`\d+\.\d`).FindString(fmt.Sprintf("%v", mBody["data"].(map[string]interface{})["soft_version"]))

	return version, nil
}

// Status retrieves the current status of the server from its dashboard.
func (inst *VIPNetIDSClient) Status(idsUrl string) ([]model.Status, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%vservice/dashboard", idsUrl), nil)
	if err != nil {
		return nil, err
	}

	resp, err := inst.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("dashboard request failed: unexpected status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	rule := model.Responce{}
	err = json.NewDecoder(resp.Body).Decode(&rule)
	if err != nil {
		return nil, err
	}

	return rule.Data.Status, nil
}

func (inst *VIPNetIDSClient) Upload(idsUrl, filePath string, pkgType utils.PackageType) error {
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file part
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Type", "application/x-compressed-tar")
	fileHeader.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s`, file.Name()),
	)

	fileWriter, err := writer.CreatePart(fileHeader)
	if err != nil {
		return err
	}

	// Copy file contents to multipart part
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	// Add package type field
	typePart, err := writer.CreatePart(textproto.MIMEHeader{"Content-Disposition": {"form-data; name=\"type\""}})
	if err != nil {
		return err
	}

	_, err = typePart.Write([]byte(pkgType))
	if err != nil {
		return err
	}

	// Add options field
	optionsPart, err := writer.CreatePart(textproto.MIMEHeader{"Content-Disposition": {"form-data; name=\"options\""}})
	if err != nil {
		return err
	}

	_, err = optionsPart.Write([]byte("{\"mode\":\"overwrite\",\"update_active\":false,\"update_save_payload\":false,\"retain_modified_text\":false}"))
	if err != nil {
		return err
	}

	// Add CSRF token field
	csrfTokenPart, err := writer.CreatePart(textproto.MIMEHeader{"Content-Disposition": {"form-data; name=\"X-Csrf-Token\""}})
	if err != nil {
		return err
	}

	csrfToken, err := inst.extractCsrfToken(idsUrl)
	if err != nil {
		return err
	}
	_, err = csrfTokenPart.Write([]byte(csrfToken))
	if err != nil {
		return err
	}

	writer.Close()

	contextWithTimeout, cancelFunc := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancelFunc()
	// Prepare the HTTP request
	request, err := http.NewRequestWithContext(contextWithTimeout, http.MethodPost, fmt.Sprintf("%vservice/update", idsUrl), &buf)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request
	resp, err := inst.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for HTTP success code
	if resp.StatusCode != 200 {
		return errors.New("non-200 response code")
	}

	// Parse server response
	updateMalwareResponse := &model.Responce{}

	err = json.NewDecoder(resp.Body).Decode(updateMalwareResponse)
	if err != nil {
		return err
	}

	// Check if response indicates success
	if !updateMalwareResponse.Success {
		return errors.New("update not successful")
	}

	return nil
}

func (inst *VIPNetIDSClient) ping(serverURL string) error {
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		return err
	}

	resp, err := inst.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed request: unexpected status code %d", resp.StatusCode)
	}

	return nil
}

// extractCsrfToken retieves the csrf token from client cookies.
func (c *VIPNetIDSClient) extractCsrfToken(idsURL string) (string, error) {
	parsedURL, err := url.Parse(idsURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	for _, cookie := range c.client.Jar.Cookies(parsedURL) {
		if cookie.Name == "X-Csrf-Token" {
			if cookie.Value != "" {
				return cookie.Value, nil
			}
			break
		}
	}

	return "", errors.New("missing or empty CSRF token")
}
