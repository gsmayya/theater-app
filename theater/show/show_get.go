package show

import (
	"github.com/gsmayya/theater/utils"
)

/*
This will eventually will access the database and fetch new details, for now, it is dummy
*/
func GetShows() (map[string]*Show, error) {
	redis := utils.GetStoreAccess()
	allData, err := utils.GetAll(redis)
	if err != nil {
		return nil, err
	}
	shows := make(map[string]*Show)
	for key, value := range allData {
		show, err := JSONToShow(value)
		if err != nil {
			return nil, err
		}
		shows[key] = show
	}
	return shows, nil
}

func GetShow(uuid string) (*Show, error) {
	redis := utils.GetStoreAccess()
	str, err := utils.GetFromCache(uuid, redis)
	if err != nil {
		return nil, err
	}
	show, err := JSONToShow(str)
	if err != nil {
		return nil, err
	}
	return show, nil
}
