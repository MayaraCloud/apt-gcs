package gcs

import (
	"net/http"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	storage "google.golang.org/api/storage/v1"
	"log"
)

const (
	scope = storage.DevstorageReadOnlyScope
	accessTokenPath = "/etc/apt/gcs_access_token"
	serviceAccountJSONPath = "/etc/apt/gcs_sa_json"
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
	client, err := getClient()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
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

func getClient() (client *http.Client, err error) {

	switch {
	case fileExists(accessTokenPath):
		client, err = clientFromAccessToken(accessTokenPath)
		if err != nil {
			log.Fatalf("Unable to get client: %v", err)
		}
	case fileExists(serviceAccountJSONPath):
		client, err = clientFromServiceAccount(serviceAccountJSONPath)
		if err != nil {
			log.Fatalf("Unable to get client: %v", err)
		}
	default:
		client, err = google.DefaultClient(ctx, scope)
		if err != nil {
			log.Fatalf("Unable to get client: %v", err)
		}
	}
	return client, err
}

func clientFromAccessToken(accessTokenPath string) (client *http.Client, err error) {
	tokenBytes, err := ioutil.ReadFile(accessTokenPath)
	if err != nil {
		log.Fatalf("Error while reading access_token file: %v", err)
	}
	token := oauth2.Token{
		AccessToken: string(tokenBytes),
	}
	tokenSource := oauth2.StaticTokenSource(&token)
	return oauth2.NewClient(ctx, tokenSource), err
}

func clientFromServiceAccount(serviceAccountJSONPath string) (client *http.Client, err error) {
	JSONBytes, err := ioutil.ReadFile(serviceAccountJSONPath)
	if err != nil {
		log.Fatalf("Error while reading SA json file: %v", err)
	}
	credentials, err := google.CredentialsFromJSON(ctx, JSONBytes, scope)
	tokenSource := credentials.TokenSource
	return oauth2.NewClient(ctx, tokenSource), err
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}