package httpRequest

import (
        "crypto/tls"
        "io"
        "io/ioutil"
        "net/http"
)

//List of available HTTP Methods
const (
        GET     = "GET"
        POST    = "POST"
        PUT     = "PUT"
        PATCH   = "PATCH"
        HEAD    = "HEAD"
        OPTIONS = "OPTIONS"
        DELETE  = "DELETE"
)

//Request Method's response format
type Response struct {
        Header http.Header //Header of the response request
        Code   int         //StatusCode of the response request
        Body   []byte      //Body of the response request
}

//Request Method's optionals options
type ReqOptions struct {
        PostBody  io.Reader         //Body of the requisition
        Headers   map[string]string //Header of the requisition
        Transport http.RoundTripper //Transport
        Username  string            //Username
        Password  string            //Password
}

//Generic Request for all available HTTP Methods.
//It will return type Response
func Request(method, url string, p ReqOptions) (Response, error) {
        var output Response
        var client http.Client

        request, err := http.NewRequest(method, url, p.PostBody)
        client = http.Client{Transport: p.Transport}

        if err != nil {
                return output, err
        }

        request.Close = true

        if p.Headers != nil {
                for key, value := range p.Headers {
                        request.Header.Set(key, value)
                }
        }

        if len(p.Username) > 0 && len(p.Password) > 0 {
                request.SetBasicAuth(p.Username, p.Password)
        }

        response, errDo := client.Do(request)

        if errDo != nil {
                return output, errDo
        }

        defer response.Body.Close()

        content, errCon := ioutil.ReadAll(response.Body)

        if errCon != nil {
                return output, errCon
        }
        output = Response{Header: response.Header, Code: response.StatusCode, Body: content}
        return output, nil
}

func SetOptions(op ReqOptions) ReqOptions {
        transport := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient},
        }

        options := ReqOptions{
                PostBody:  op.PostBody,
                Headers:   op.Headers,
                Transport: transport,
        }

        if len(op.Username) > 0 && len(op.Password) > 0 {
                options.Username = op.Username
                options.Password = op.Password
        }

        return options
}
