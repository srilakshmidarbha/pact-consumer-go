package dps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type JsonError struct {
	Code string `json:"code" pact:"example=1234"`
	Msg  string `json:"msg" pact:"example=No customer with id 1234"`
}

type Response struct {
	Customer Customer `json:"customer" pact:"example=Customer{}"`
}
type Customer struct {
	Variant string `json:"variant" pact:"example=Original"`
}

func NewClient(host string) Client {
	return Client{host}
}

type Client struct {
	host string
}

func (c Client) PostCustomer(countrycode string) (Response, error) {
	client := http.Client{}

	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v2/fees/%s", c.host, countrycode),
		nil,
		)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println("Host:", req.Host)
	resp, err := client.Do(req)
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println("Error is", err)
	fmt.Println("Response is", bodyString)
	if resp.StatusCode == http.StatusNotFound {
		je := JsonError{}
		json.NewDecoder(resp.Body).Decode(&je)

		return Response{}, errors.New(fmt.Sprintf("%s - %s", je.Code, je.Msg))
	}

	d := Response{}
	json.Unmarshal(bodyBytes,&d)

	return d, nil

}
