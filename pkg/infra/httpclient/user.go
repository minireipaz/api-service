package httpclient

type UserHTTPRepository struct {
	client HTTPClient
}

func NewUserClientHTTP(client HTTPClient) *UserHTTPRepository {
	return &UserHTTPRepository{client: client}
}
