package connexion

import (
	model "bucketool/model"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	color "bucketool/utils/colorPrint"
)

// var date string = time.Now().UTC().Format(http.TimeFormat)
 var region string = "us-west-1"
// var typeServer string = "s3"
//var requestType string = "aws4_request"

// Connexion is a struct that contains the information needed to connect to a API S3/ minio

type Connexion struct {
	// Token is the token to connect to the API
	URL string
	Auth string
	Alias model.Alias
	Request *http.Request
	Response *http.Response
	S3Service *s3.S3
}

// NewConnexion creates a new Connexion struct
func Use(Alias model.Alias) *Connexion {
	url := CreateURL(Alias)
	c := &Connexion{
		URL : url,
		Auth: "",
		Alias: Alias,
	}
	//c.setHMACSignature()
	c.InitS3Service()
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
	c.addCommonHeaders()


	// Send the request
	response, err := client.Do(request)
	if err != nil {
		return nil,err
	}

	return response, nil
}


func(c *Connexion) InitS3Service() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Endpoint: aws.String(c.URL),
		Credentials: credentials.NewStaticCredentials(c.Alias.KeyName, c.Alias.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Utiliser le style de chemin pour MinIO
    })

	if err != nil {
		println(color.RedP("Error creating session"), err)
	}

	// Create S3 service client
	svc := s3.New(sess)
	c.S3Service = svc
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


// Send sends the request and returns the response
func (c *Connexion) Send() (*http.Response, error) {
	// Configurer le transport pour forcer l'utilisation de HTTP/1.1
	transport := &http.Transport{
		ForceAttemptHTTP2: false, // Désactiver HTTP/2
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: transport,
	}
	response, err := client.Do(c.Request)
	if err != nil {
		return nil, err
	}
	c.Response = response
	return response, nil
}

