package httpclient

type CredentialFacebookHTTPRepository struct {
	client HTTPClient
}

func NewCredentialFacebookRepository(httpCli HTTPClient) *CredentialFacebookHTTPRepository {
	return &CredentialFacebookHTTPRepository{
		client: httpCli,
	}
}
