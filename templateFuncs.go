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
		var draftProp hyper.Property
		draftProp, _ = plr.Properties.Find("draftIndex")
		if draftProp.Value == index {
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
