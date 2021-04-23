package main

import (
	"context"
	"github.com/ctreminiom/go-atlassian/jira"
	"log"
	"os"
)

func main() {

	var (
		host  = os.Getenv("HOST")
		mail  = os.Getenv("MAIL")
		token = os.Getenv("TOKEN")
	)

	atlassian, err := jira.New(nil, host)
	if err != nil {
		return
	}

	atlassian.Auth.SetBasicAuth(mail, token)

	var (
		fieldID   = "customfield_10038"
		contextID = 10180

		payload = &jira.FieldContextOptionListScheme{
			Options: []*jira.CustomFieldContextOptionScheme{

				// Single/Multiple Choice example
				{
					ID:       "10064",
					Value:    "Option 3 - Updated",
					Disabled: false,
				},
				{
					ID:       "10065",
					Value:    "Option 4 - Updated",
					Disabled: true,
				},

				///////////////////////////////////////////
				/*
					// Cascading Choice example
					{
						OptionID: "1027",
						Value:    "Argentina",
						Disabled: false,
					},
					{
						OptionID: "1027",
						Value:    "Uruguay",
						Disabled: false,
					},
				*/

			}}
	)

	fieldOptions, response, err := atlassian.Issue.Field.Context.Option.Update(context.Background(), fieldID, contextID, payload)
	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", string(response.BodyAsBytes))
		}
		return
	}

	log.Println("Response HTTP Code", response.StatusCode)
	log.Println("HTTP Endpoint Used", response.Endpoint)

	for _, option := range fieldOptions.Options {
		log.Println(option)
	}

}