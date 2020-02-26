package gcs

import (
	"net/http"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"io/ioutil"
	storage "google.golang.org/api/storage/v1"
	"log"
)

const (
	scope = storage.DevstorageReadOnlyScope
	accessTokenPath = "/etc/apt/gcs_access_token"
)

var (
	client   *context.Context
	service  *storage.Service
	oService *storage.ObjectsService
)

var ctx context.Context = context.Background()

//InitConfig creates the google storage client from the
//application default credentials or from an access_token
func InitConfig() {
	client, err := google.DefaultClient(ctx, scope)
	if err != nil {
		client, err = clientFromAccessToken(accessTokenPath)
		if err != nil {
			log.Fatalf("Unable to get client: %v", err)
		}
	}
	service, err = storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	oService = storage.NewObjectsService(service)
	if err != nil {
		log.Fatalf("Unable to create objects storage service: %v", err)
	}

}

func clientFromAccessToken(accessTokenPath string) (client *http.Client, err error) {
	tokenBytes, err := ioutil.ReadFile(accessTokenPath)
	token := oauth2.Token{
		AccessToken: string(tokenBytes),
	}
	tokenSource := oauth2.StaticTokenSource(&token)
	return oauth2.NewClient(ctx, tokenSource), err
}