package discovery

import (
	"log"

	"github.com/ark-go/yeelight/pkg/jt"
)

// копируем данные от лампы в массив ламп и проверим что изменилось
func copyNewLampData(lamp *jt.Lamp) string {
	if lamp.ErrorResult {
		return ""
	}

	if lamp.Id == "" {
		return ""
	}
	if _, ok := jt.Lamps[lamp.Id]; !ok {
		jt.Lamps[lamp.Id] = lamp
		jt.Lamps[lamp.Id].IsChanges = true
		log.Println("Новая лампа !!!")
		return jt.Lamps[lamp.Id].Id
	}
	jt.Lamps[lamp.Id].NTS = ""
	if jt.Lamps[lamp.Id].Power != lamp.Power ||
		jt.Lamps[lamp.Id].Bright != lamp.Bright ||
		jt.Lamps[lamp.Id].ColorMode != lamp.ColorMode ||
		jt.Lamps[lamp.Id].Ct != lamp.Ct ||
		jt.Lamps[lamp.Id].Rgb != lamp.Rgb ||
		jt.Lamps[lamp.Id].Hue != lamp.Hue ||
		jt.Lamps[lamp.Id].Sat != lamp.Sat ||
		jt.Lamps[lamp.Id].NTS != lamp.NTS ||
		jt.Lamps[lamp.Id].Name != lamp.Name {
		//  копируем если было изменение
		jt.Lamps[lamp.Id] = lamp
		jt.Lamps[lamp.Id].IsChanges = true
	} else {
		jt.Lamps[lamp.Id].IsChanges = false
	}
	return ""
}
