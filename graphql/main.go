package main

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"log"
)

type HubById struct {
	hubById Hub
}
type Hub struct {
	id                      string `json:"id"`
	dakiPlusHubId           string `json:"dakiPlusHubId"`
	fulfillmentHubIdForPlus string `json:"fulfillmentHubIdForPlus"`
	state                   string `json:"state"`
	timeZone                string `json:"timeZone"`
}

type ResponseData struct {
	hubById struct {
		id                      string `json:"id"`
		dakiPlusHubId           string `json:"dakiPlusHubId"`
		fulfillmentHubIdForPlus string `json:"fulfillmentHubIdForPlus"`
		state                   string `json:"state"`
		timeZone                string `json:"timeZone"`
	} `json:"hubById"`
}

type Hubs struct {
	hubs []Hub
}

func main() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://api-sbx-mx.gcp-mioxxo.com")

	// make a request
	//	request := graphql.NewRequest(`
	//    	query hubs {
	//    		hubs {
	//        		id
	//        		dakiPlusHubId
	//        		fulfillmentHubIdForPlus
	//        		state
	//        		timeZone
	//    		}
	//		}
	//`)
	request := graphql.NewRequest(`
	query hubById ($id: ID!) {
		hubById (id: $id) {
			id
			dakiPlusHubId
			fulfillmentHubIdForPlus
			state
			timeZone
		}
	}
`)

	request.Var("id", "50645")

	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("apollo-require-preflight", "false")

	ctx := context.Background()

	var respData ResponseData
	if err := client.Run(ctx, request, &respData); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", respData)
}
