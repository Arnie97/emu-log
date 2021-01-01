package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
)

// railMapHandler redirects the user to the page for a given railway station
// if an exact match for the station name was found on the railway map, or to
// the home page of the map website otherwise.
func railMapHandler(w http.ResponseWriter, r *http.Request) {
	const site = "http://cnrail.geogv.org"
	stationID := ""
	stationName := chi.URLParam(r, "stationName")
	defer func() {
		http.Redirect(w, r, fmt.Sprintf(
			"%s/zhcn/station/%s?useMapboxGl=true", site, stationID,
		), http.StatusSeeOther)
	}()

	keyword := stationName
	if len(stationName) > 2 {
		keyword = strings.TrimSuffix(stationName, "所")
	}
	resp, err := common.HTTPClient(http.DefaultTransport).Get(fmt.Sprintf(
		"%s/api/v1/match_feature/%s?locale=zhcn", site, keyword,
	))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	matches := struct {
		Success bool
		Data    [][3]string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil || !matches.Success {
		return
	}

	for _, m := range matches.Data {
		itemID, itemType, itemName := m[0], m[1], m[2]
		if itemType != "STATION" {
			continue
		} else if strings.Replace(itemName, "线路所", "所", 1) == stationName {
			stationID = itemID
			return
		}
	}
}
