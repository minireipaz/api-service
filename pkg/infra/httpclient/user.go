package httpclient

import (
	"minireipaz/pkg/domain/models"

)

type UserRepository struct {
	client HTTPClient
}

func NewUserClientHTTP(client HTTPClient) *UserRepository {
	return &UserRepository{client: client}
	// return &UserHTTPClient{
	//   client: &http.Client{Timeout: 10 * time.Second},
	// }
}

func (u *UserRepository) Create(_ *models.Users) (created bool, exist bool) {
	return false, false
}
