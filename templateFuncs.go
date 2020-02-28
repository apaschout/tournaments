package tournaments

import (
	"fmt"
	"html/template"

	"github.com/cognicraft/hyper"
)

func createSeatForIndex(plrs hyper.Items, index int) template.HTML {
	var name string
	var href string
	for _, plr := range plrs {
		var seatProp hyper.Property
		seatProp, _ = plr.Properties.Find("seatIndex")
		if seatProp.Value == index {
			for _, link := range plr.Links {
				if link.Rel == "details" {
					href = link.Href
				}
			}
			prop, _ := plr.Properties.Find("name")
			name = prop.Value.(string)
		}
	}
	return template.HTML(fmt.Sprintf(`<a class="seat flex-container" title="%s" href="%s" target="_blank" style="text-decoration: none;">%d</a>`, name, href, index+1))
}

func getDraftIndex(plr hyper.Item) int {
	prop, _ := plr.Properties.Find("draftIndex")
	return prop.Value.(int)
}

func getDetails(item hyper.Item) string {
	var res string
	for _, link := range item.Links {
		if link.Rel == "details" {
			res = link.Href
		}
	}
	return res
}

func getStart(trn hyper.Item) string {
	prop, _ := trn.Properties.Find("start")
	return prop.Value.(string)
}

func getName(item hyper.Item) string {
	prop, _ := item.Properties.Find("name")
	return prop.Value.(string)
}

func getFormat(trn hyper.Item) string {
	prop, _ := trn.Properties.Find("format")
	return prop.Value.(string)
}

func matches(trn hyper.Item) []Match {
	prop, _ := trn.Properties.Find("matches")
	return prop.Value.([]Match)
}

func actionByRel(trn hyper.Item, rel string) hyper.Action {
	var res hyper.Action
	for _, ac := range trn.Actions {
		if ac.Rel == rel {
			res = ac
			break
		}
	}
	return res
}

func participantNameByID(trn hyper.Item, ID PlayerID) string {
	var res string
	parts := trn.Items[0]
	for _, part := range parts.Items {
		if part.ID == string(ID) {
			prop, _ := part.Properties.Find("name")
			res = prop.Value.(string)
			break
		}
	}
	return res
}
