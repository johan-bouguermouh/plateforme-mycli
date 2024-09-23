package connexion

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Initialize the HTTP request
func (c *Connexion) initRequest() {
    req, err := http.NewRequest("GET", c.URL, nil)
    if err != nil {
        fmt.Println("Failed to create request:", err)
        return
    }
    c.Request = req
    c.addCommonHeaders()
}

// Add common headers like Authorization and Date to the request
func (c *Connexion) addCommonHeaders() {
    date := time.Now().UTC().Format(http.TimeFormat)
    c.Request.Header.Add("x-amz-date", date)
	c.Request.Header.Add("Host", c.Request.URL.Host) 

	c.TypeXML()
}

/** METHODS **/

func (c *Connexion) POST() *Connexion {
	c.Request.Method = "POST"
	return c
}

func (c *Connexion) GET() *Connexion {
	c.Request.Method = "GET"
	return c
}

func (c *Connexion) PUT() *Connexion {
	c.Request.Method = "PUT"
	return c
}

func (c *Connexion) DELETE() *Connexion {
	c.Request.Method = "DELETE"
	return c
}

func (c *Connexion) HEAD() *Connexion {
	c.Request.Method = "HEAD"
	return c
}

func (c *Connexion) OPTIONS() *Connexion {
	c.Request.Method = "OPTIONS"
	return c
}

func (c *Connexion) TRACE() *Connexion {
	c.Request.Method = "TRACE"
	return c
}

func (c *Connexion) PATCH() *Connexion {
	c.Request.Method = "PATCH"
	return c
}

func (c *Connexion) GetMethod() string {
	return c.Request.Method
}

/** HEADERS **/

	/** GENERAL **/

func (c *Connexion) AddHeader(key, value string) *Connexion {
	c.Request.Header.Add(key, value)
	return c
}

func (c *Connexion) SetHeader(key, value string) *Connexion {
	//on verfie si le header existe
	if _, ok := c.Request.Header[key]; ok {
		c.Request.Header.Set(key, value)
	} else {
		c.Request.Header.Add(key, value)
	}
	return c
}

func (c *Connexion) RemoveHeader(key string) *Connexion {
	//on verfie si le header existe
	if _, ok := c.Request.Header[key]; ok {
		c.Request.Header.Del(key)
	}
	return c
}

func (c *Connexion) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Connexion) GetHeaders() http.Header {
	return c.Request.Header
}

	/** SPECIFIC **/
	
func (c *Connexion) TypeXML() *Connexion {
	c.Request.Header.Add("Content-Type", "application/xml")
	return c
}

func (c *Connexion) TypeJSON() *Connexion {
	c.Request.Header.Add("Content-Type", "application/json")
	return c
}

func (c *Connexion) TypeText() *Connexion {
	c.Request.Header.Add("Content-Type", "text/plain")
	return c
}

func (c *Connexion) TypeHTML() *Connexion {
	c.Request.Header.Add("Content-Type", "text/html")
	return c
}

/** ROUTE **/

func (c *Connexion) Route(route string) *Connexion {
	c.Request.URL.Path = route
	return c
}


/** BODY **/

// SendXMLBody sets the request body to the provided XML string and updates the headers accordingly
func (c *Connexion) SendXMLBody(body string) *Connexion {
    c.Request.Header.Add("Content-Type", "application/xml")
    c.Request.Header.Add("Content-Length", fmt.Sprintf("%d", len(body)))
    c.Request.Body = io.NopCloser(strings.NewReader(body))
    return c
}



