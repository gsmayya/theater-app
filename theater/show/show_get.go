package show

import (
	"github.com/gsmayya/theater/utils"
)

/*
This will eventually will access the database and fetch new details, for now, it is dummy
*/
func GetShows() (map[string]string, error) {
	// mock -- Need to get from redis
	show := NewShow("show1", "Movie 1", 100, 50, "Location 1")
	return ShowToMap(show), nil
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
