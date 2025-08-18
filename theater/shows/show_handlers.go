package shows

import (
	"log"

	"github.com/gsmayya/theater/utils"
)

/*
This will eventually will access the database and fetch new details, for now, it is dummy
*/
func GetShows() (map[string]*ShowData, error) {
	redis := utils.GetStoreAccess()
	allData, err := utils.GetAll(redis)
	if err != nil {
		return nil, err
	}
	shows := make(map[string]*ShowData)
	for key, value := range allData {
		show_info := &ShowData{}
		if _, err := show_info.JSONToShow(value); err != nil {
			return nil, err
		}
		shows[key] = show_info
	}
	return shows, nil
}

// GetShow retrieves a show from the cache using the provided UUID.
// It returns the retrieved show and an error if any.
func GetShow(uuid string) (*ShowData, error) {
	redis := utils.GetStoreAccess()
	showData, err := utils.GetFromCache(uuid, redis)
	if err != nil {
		return nil, err
	}
	show_info := &ShowData{}
	if _, err := show_info.JSONToShow(showData); err != nil {
		return nil, err
	}
	return show_info, nil
}

func PutShowDetails(name string, details string, price int32, totalTickets int32, location string) error {
	// This function will eventually update the show details in the database
	show_info := &ShowData{}
	show_info.NewShow(name, details, price, totalTickets, location)
	log.Println(show_info)
	return add(show_info)
}

func PutShow(info *ShowData) error {
	return add(info)
}

func add(showData *ShowData) error {
	jsonData, err := showData.ShowToJSON()
	if err != nil {
		return err
	}
	storeAccess := utils.GetStoreAccess()
	err = utils.AddToCache(showData.Show_Id.String(), jsonData, storeAccess)
	if err != nil {
		return err
	}
	return nil
}
