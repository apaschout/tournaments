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

func createMatches(trn hyper.Item) template.HTML {
	res := ""
	prop, _ := trn.Properties.Find("matches")
	for _, match := range prop.Value.([]Match) {
		res += fmt.Sprintf("<p>%s vs %s</p>", match.Player1, match.Player2)
	}
	return template.HTML(res)
}
