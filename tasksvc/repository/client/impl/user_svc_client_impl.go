package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tasksvc/config"
	"tasksvc/models"
	"tasksvc/repository/client"
	"tasksvc/utils"
)

type UserServiceClient struct {
	httpClient *utils.HttpClient
	config     *config.StartupConfig
}

func NewUserServiceClient(httpClient *utils.HttpClient, config *config.StartupConfig) client.UserServiceClient {
	return &UserServiceClient{
		httpClient,
		config,
	}
}

func (client UserServiceClient) CheckAdmin(ctx context.Context, username string) (bool, error) {
	credentials := username + ":"
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	headers := make(map[string][]string)
	headers["Authorization"] = []string{basicAuth}
	headers["Content-Type"] = []string{"application/json"}
	url := client.config.SvcBase.TaskSvc + "/check-admin"
	res, err := client.httpClient.Get(url, headers, nil)
	if err != nil {
		return false, err
	}
	if res != nil && res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		var response models.StatusMessageResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return false, err
		}
		adminResponse := response.Data.(map[string]interface{})
		return adminResponse["isAdmin"].(bool), nil

	}
	return false, errors.New(fmt.Sprintf("Could not check admin. Request failed with %d", res.StatusCode))
}
