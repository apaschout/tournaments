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

func getDetails(item hyper.Item) string {
	var res string
	for _, link := range item.Links {
		if link.Rel == "details" {
			res = link.Href
		}
	}
	return res
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

func propertyByName(trn hyper.Item, name string) interface{} {
	prop, _ := trn.Properties.Find(name)
	return prop.Value
}

func wins(m Match) string {
	if m.P1Count <= m.P2Count {
		return fmt.Sprintf("(%d : %d)", m.P2Count, m.P1Count)
	} else {
		return fmt.Sprintf("(%d : %d)", m.P1Count, m.P2Count)
	}
}

func getParticipants(trn hyper.Item) []hyper.Item {
	res := []hyper.Item{}
	partItem := hyper.Item{}
	for _, v := range trn.Items {
		if v.Type == "participants" {
			partItem = v
			break
		}
	}
	for _, v := range partItem.Items {
		if v.Type == "participant" {
			res = append(res, v)
		}
	}
	return res
}

func sortParticipants(items []hyper.Item) []hyper.Item {
	if len(items) == 1 {
		return items
	}
	middle := int(len(items) / 2)
	var (
		l = make([]hyper.Item, middle)
		r = make([]hyper.Item, len(items)-middle)
	)
	for i := 0; i < len(items); i++ {
		if i < middle {
			l[i] = items[i]
		} else {
			r[i-middle] = items[i]
		}
	}
	return merge(sortParticipants(l), sortParticipants(r))
}

func merge(l, r []hyper.Item) []hyper.Item {
	res := make([]hyper.Item, len(l)+len(r))
	i := 0
	for len(l) > 0 && len(r) > 0 {
		matchesL, _ := l[0].Properties.Find("matchWins")
		matchesR, _ := r[0].Properties.Find("matchWins")

		gamesL, _ := l[0].Properties.Find("gameWins")
		gamesR, _ := r[0].Properties.Find("gameWins")
		if matchesL.Value.(int) > matchesR.Value.(int) {
			res[i] = l[0]
			l = l[1:]
		} else if matchesL.Value.(int) < matchesR.Value.(int) {
			res[i] = r[0]
			r = r[1:]
		} else if gamesL.Value.(int) > gamesR.Value.(int) {
			res[i] = l[0]
			l = l[1:]
		} else {
			res[i] = r[0]
			r = r[1:]
		}
		i++
	}

	for j := 0; j < len(l); j++ {
		res[i] = l[j]
		i++
	}
	for j := 0; j < len(r); j++ {
		res[i] = r[j]
		i++
	}
	return res
}
