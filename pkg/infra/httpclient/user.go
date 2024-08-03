package httpclient

import (
	"minireipaz/pkg/domain/models"
  // "minireipaz/pkg/domain/repos"
	// "net/http"
	// "time"
)

type UserRepository struct {
	client HttpClient
}


func NewUserClientHTTP(client HttpClient) *UserRepository {
  return &UserRepository{client: client}
	// return &UserHTTPClient{
  //   client: &http.Client{Timeout: 10 * time.Second},
  // }
}

func (u *UserRepository) Create(user *models.Users) (created bool, exist bool) {

	return false, false
}
