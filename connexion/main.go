package connexion

import (
	model "bucketool/model"
	color "bucketool/utils/colorPrint"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// var date string = time.Now().UTC().Format(http.TimeFormat)
 var Region string = "us-west-1"
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
	S3Client *s3.Client
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
	c.InitS3Service(context.Background())
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


func (c *Connexion) InitS3Service(ctx context.Context) error {
    // Assurez-vous que l'URL inclut le protocole
    url := c.URL
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        url = "http://" + url // ou "https://" selon votre configuration
    }

    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion("us-east-1"), // MinIO utilise généralement "us-east-1" comme région par défaut
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            c.Alias.KeyName,
            c.Alias.SecretKey,
            "",
        )),
		config.WithHTTPClient(&http.Client{
			Transport: newLoggingRoundTripper(nil),
		}),
    )
    if err != nil {
        println(color.RedP("unable to load SDK config:"))
		return err
    }

    // Create S3 service client with BaseEndpoint
    c.S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(url)
		o.UsePathStyle = true
    })

	return nil
}



func CreateURL(Alias model.Alias) string {
    // Vérifie si le dernier caractère de l'hôte est un "/"
    if Alias.HOST[len(Alias.HOST)-1] != '/' {
        // Insère le numéro de port avant le "/"
        return Alias.HOST + ":" + strconv.Itoa(Alias.Port) 
    } else {
        sanityzeHost := Alias.HOST[:len(Alias.HOST)-1]
        return sanityzeHost + ":" + strconv.Itoa(Alias.Port) 
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

