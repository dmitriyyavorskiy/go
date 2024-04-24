package main

import (
	"fmt"
	"github.com/oxxo-labs/go/migration/services"
	"github.com/oxxo-labs/go/models"
	"log"
)

func main() {

	tokenResponse, err := services.Auth()

	if err != nil {
		log.Fatalf("Could not authenticate: %v", err)
	}

	fmt.Printf("Access Token Response %+v\n", tokenResponse.AccessToken)

	documentsResponse, err := services.GetDocumentList(tokenResponse.AccessToken, models.DocumentListRequest{
		DepartmentId: "10MAN",
		StoreId:      "504LX",
		Source:       "CT",
	})

	if err != nil {
		log.Fatalf("Could not get documents: %v", err)
	}

	fmt.Printf("Document Response %+v\n", documentsResponse)

}
