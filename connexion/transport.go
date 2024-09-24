package connexion

import (
	env "bucketool/environment"
	colorPrint "bucketool/utils/colorPrint"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-xmlfmt/xmlfmt"
)

type loggingRoundTripper struct {
    rt http.RoundTripper
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if env.IsDebugMode {
	   return lrt.LogRequest(req)
	}
	return lrt.rt.RoundTrip(req)
}

func newLoggingRoundTripper(rt http.RoundTripper) http.RoundTripper {
    if rt == nil {
        rt = http.DefaultTransport
    }
    return &loggingRoundTripper{rt: rt}
}

func (lrt *loggingRoundTripper) LogRequest(req *http.Request) (*http.Response, error) {
	// Afficher l'URL complète de la requête
	println(colorPrint.ColorPrint("Black", "      INTERCEPT HTTP REQUEST ACTIVE      ", &colorPrint.Options{
		Background : "White",
		Bold : true,
		Underline : true,
	})+ "\n")

	println("  "+colorPrint.BlueP("Method : "), req.Method)
	println("  "+colorPrint.BlueP("URL : "), req.URL.String())
	println("  "+colorPrint.BlueP("Proto : "), req.Proto)
	println("  "+colorPrint.BlueP("Host : "), req.Host)
	println("  "+colorPrint.BlueP("RemoteAddr : "), req.RemoteAddr)
	println("  "+colorPrint.BlueP("RequestURI : "), req.RequestURI)
	println("  "+colorPrint.BlueP("UserAgent : "), req.UserAgent() + "\n")

	println(colorPrint.ColorPrint("Black", " Request Header ", &colorPrint.Options{
		Background : "Blue",
		Bold : true,
		Underline : true,
	})+ "\n")


	// Afficher les en-têtes de la requête
	for name, values := range req.Header {
		for _, value := range values {
			println("  "+name+" : ", colorPrint.GreyP(value))
		}
	}

	println("\n"+colorPrint.ColorPrint("Black", " Request Body ", &colorPrint.Options{
		Background : "Blue",
		Bold : true,
		Underline : true,
	})+ "\n")

	// Afficher le corps de la requête, si applicable
	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			println(colorPrint.RedP("Error reading request body: %v"), err)
		} else {
		println(xmlfmt.FormatXML(string(bodyBytes), "", "  "),"\n")
			// Remettre le corps de la requête pour qu'il puisse être lu à nouveau
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	// Effectuer la requête
	resp, err := lrt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Afficher les en-têtes de la réponse
	println(colorPrint.ColorPrint("Black", " Response Header ", &colorPrint.Options{
		Background : "Green",
		Bold : true,
		Underline : true,
	})+ "\n")


	for name, values := range resp.Header {
		for _, value := range values {
			println("  "+name+" : ", colorPrint.GreyP(value))
		}
	}

	println("\n\n"+colorPrint.ColorPrint("Black", " Response Body ", &colorPrint.Options{
		Background : "Green",
		Bold : true,
		Underline : true,
	})+ "\n")

	// Afficher le corps de la réponse, si applicable
	if resp.Body != nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf(colorPrint.RedP("Error reading response body: %v"), err)
		} else {
			println(xmlfmt.FormatXML(string(bodyBytes), "", "  "))
			
			// Remettre le corps de la réponse pour qu'il puisse être lu à nouveau
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	println("\n\n" + colorPrint.ColorPrint("Black", "        END OF INTERCEPTION        ", &colorPrint.Options{
		Background : "White",
		Bold : true,
		Underline : true,
	})+ "\n")

	return resp, nil
}