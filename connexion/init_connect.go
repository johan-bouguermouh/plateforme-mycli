package connexion

import (
	model "bucketool/model"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var date string = time.Now().UTC().Format(http.TimeFormat)

// Connexion is a struct that contains the information needed to connect to a API S3/ minio

type Connexion struct {
	// Token is the token to connect to the API
	URL string
	Auth string
	Alias model.Alias
}

// NewConnexion creates a new Connexion struct
func Use(Alias model.Alias) *Connexion {
	url := CreateURL(Alias)
	c := &Connexion{
		URL : url,
		Auth: "",
		Alias: Alias,
	}
	c.setHMACSignature()
	return c
}

// Connect connects to the API
func (c *Connexion) Connect() (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Create the request
	request, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return nil,err
	}

	// Add the headers
	c.addCommonHeaders(request)


	// Send the request
	response, err := client.Do(request)
	if err != nil {
		return nil,err
	}

	return response, nil
}


func(c *Connexion)setHMACSignature() {
	secretKey := c.Alias.SecretKey
	stringToSign := fmt.Sprintf("GET\n\n\n%s\n/", date)
	// Use protocole HMAC-SHA256 to generate the signature of the request, for v4 signature of AWS
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write([]byte(stringToSign))

	//encoded signature on base 64
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	authHeader := fmt.Sprintf("AWS %s:%s", c.Alias.KeyName, signature)
	c.Auth = authHeader
}

// Adds common headers like Authorization and Date to the request
func (c *Connexion) addCommonHeaders(req *http.Request) {
    date := time.Now().UTC().Format(http.TimeFormat)
    req.Header.Add("Authorization", c.Auth)
    req.Header.Add("Date", date)
}

func CreateURL(Alias model.Alias) string {
    // Vérifie si le dernier caractère de l'hôte est un "/"
    if Alias.HOST[len(Alias.HOST)-1] != '/' {
        // Insère le numéro de port avant le "/"
        return Alias.HOST + ":" + strconv.Itoa(Alias.Port) + "/"
    } else {
        sanityzeHost := Alias.HOST[:len(Alias.HOST)-1]
        return sanityzeHost + ":" + strconv.Itoa(Alias.Port) + "/"
    }
}

