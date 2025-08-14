package show

import (
	"log"

	"github.com/gsmayya/theater/utils"
)

func PutShowDetails(name string, details string, price int32, totalTickets int32, location string) error {
	// This function will eventually update the show details in the database
	show_info := NewShow(name, details, price, totalTickets, location)
	log.Println(show_info)
	return add(show_info)
}

func PutShow(show_info *Show) error {
	// This function will eventually update the show details in the database
	log.Println(show_info)
	return add(show_info)
}

func add(show_info *Show) error {
	log.Println(show_info)
	redis := utils.GetStoreAccess()
	show_info_str, err := ShowToJSON(show_info)
	if err != nil {
		return err
	}
	utils.AddToCache(show_info.Show_Id.String(), show_info_str, redis)
	return nil
}
