package utils

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	//EventbriteEnvKey the key name for auth
	EventbriteEnvKey = "EVENTBRITE_AUTH"
	//EventbriteURL the eventbrite url for sfblock week
	//todo (dont forget to change the EVENTBRITE_AUTH token when switching these)
	//EventbriteURL = "https://www.eventbriteapi.com/v3/events/64449323662/attendees?page=%d"

	//EventbriteURL for fake group
	EventbriteURL = "https://www.eventbriteapi.com/v3/events/74857698391/attendees?page=%d"
	//HeaderPrefix the prefix to the value for auth with eventbrite
	HeaderPrefix = "Bearer %s"
)

//EventbriteData is the return data of a page call
type EventbriteData struct {
	PaginationInfo Pagination             `json:"pagination"`
	Attendees      longy.GenesisAttendees `json:"attendees"`
}

//Pagination is the pagination info of a page call
type Pagination struct {
	Count     int `json:"object_count"`
	Page      int `json:"page_number"`
	PageCount int `json:"page_count"`
}

//GetAttendees gets the attendee list from eventbrite while using the auth key found in an environmental var
func GetAttendees() (ga longy.GenesisAttendees, err sdk.Error) {
	authKey, isSet := os.LookupEnv(EventbriteEnvKey)
	if !isSet {
		err = types.ErrEventbriteEnvVariableNotSet("EVENTBRITE_AUTH must be set for generating the genesis file")
		return
	}
	ga, err = fetchAttendees(authKey)
	return
}

//fetchAttendees async gets and process the index of attendees from the paginated endpoint
func fetchAttendees(authKey string) (ga longy.GenesisAttendees, e sdk.Error) {
	client := http.Client{}

	headerAuth := fmt.Sprintf(HeaderPrefix, authKey)

	data, e := processPage(&client, 1, headerAuth)
	if e != nil {
		return
	}
	totalAttendees := data.PaginationInfo.Count
	aChan := make(chan longy.GenesisAttendee, totalAttendees)
	eChan := make(chan sdk.Error, totalAttendees)
	ga = data.Attendees

	var wg sync.WaitGroup

	for i := 2; i <= data.PaginationInfo.PageCount; i++ {
		wg.Add(1)
		go asyncGet(i, &wg, &client, headerAuth, aChan, eChan)
	}
	wg.Wait()

	if len(eChan) != 0 {
		e = <-eChan
		return
	}

	ga = mergeAttendees(aChan, ga)
	if len(ga) != totalAttendees {
		e = types.ErrAttendeeCountMismatch(
			"the total attendees should be %d, but we only got %d", totalAttendees, len(ga))
	}
	return ga, e
}

//mergeAttendees merges the first paginated call for attendees with the attendee channel populated from subsequent calls
func mergeAttendees(ac chan longy.GenesisAttendee, ga longy.GenesisAttendees) longy.GenesisAttendees {
	temp := make(longy.GenesisAttendees, len(ac))
	i := 0
	close(ac)
	for d := range ac {
		temp[i] = d
		i++
	}
	ga = append(ga, temp...)
	return ga
}

//asyncGet gets and writes the attendees from a request into the data channel from a go routine
//nolint: gocritic
func asyncGet(i int, wg *sync.WaitGroup, client *http.Client, headerAuth string, aChan chan longy.GenesisAttendee,
	eChan chan<- sdk.Error) {
	defer wg.Done()
	da, err := processPage(client, i, headerAuth)

	if err != nil {
		eChan <- err
		return
	}

	for _, a := range da.Attendees {
		aChan <- a
	}
}

//processPage gets and processes a single page returning its data
func processPage(client *http.Client, page int, headerAuth string) (data EventbriteData, e sdk.Error) {
	var res *http.Response
	res, e = getPage(client, page, headerAuth)
	if e != nil {
		return
	}
	data, e = processResp(res)
	return
}

//processResp processes a response for attendees into a struct
func processResp(res *http.Response) (data EventbriteData, e sdk.Error) {
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		e = types.ErrDefault(err.Error())
		return
	}
	return
}

//getPage gets the a paginated result off attendees
func getPage(client *http.Client, page int, headerAuth string) (res *http.Response, e sdk.Error) {
	var err error
	req := pageURL(page, headerAuth)
	res, err = client.Do(req)
	if err != nil {
		e = types.ErrNetworkResponseError(
			fmt.Sprintf("eventbrite call failed : %s", err.Error()))
		return
	}

	if res.StatusCode != http.StatusOK {
		e = types.ErrNetworkResponseError(
			fmt.Sprintf("eventbrite call returned with code : %d", res.StatusCode))
		return
	}
	return
}

//pageURL creates an authorized request for a paginated call to get attendees
func pageURL(page int, header string) *http.Request {
	url := fmt.Sprintf(EventbriteURL, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", header)
	return req
}
