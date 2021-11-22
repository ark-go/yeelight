package discovery

import (
	"strconv"
	"strings"

	"github.com/ark-go/yeelight/pkg/jt"
)

// ответ от ламы в структуру
func parseStrToStruct(str string, lamp *jt.Lamp) error {
	res := strings.SplitN(str, ":", 2)
	if len(res) < 2 {
		return nil
	}
	res[0] = strings.Trim(res[0], " ")
	res[1] = strings.Trim(res[1], " ")

	switch res[0] {
	case "Cache-Control": //max - age = 3600
		lamp.CacheControl = res[1]
	case "Date":
	case "Ext":
	case "NTS": // only Notify (ssdp:alive)
		lamp.NTS = res[1]
	case "Location": //yeelight://172.16.172.44:55443
		res := strings.SplitN(res[1], ":", 2)
		if len(res) < 2 {
			lamp.Location = ""
			break
		}

		lamp.Location = strings.TrimLeft(res[1], "/")
	case "Server": //POSIX UPnP/1.0 YGLC/1
		lamp.Server = res[1]
	case "id": //0x0000000018430941
		lamp.Id = res[1]
	case "model": // ceila
		lamp.Model = res[1]
	case "fw_ver": //11
		lamp.FwVer = res[1]
	case "support": //get_prop set_default set_power toggle set_bright set_scene cron_add cron_get cron_del start_cf stop_cf set_adjust adjust_bright set_name set_ct_abx adjust_ct
		lamp.SupportAll = res[1]
	case "power": //off
		lamp.Power = res[1]
	case "bright": //100
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.Bright = -1
		} else {
			lamp.Bright = i
		}
	case "color_mode": //2
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.ColorMode = -1
		} else {
			lamp.ColorMode = i
		}
	case "ct": //2700
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.Ct = -1
		} else {
			lamp.Ct = i
		}
	case "rgb": //0
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.Rgb = -1
		} else {
			lamp.Rgb = i
		}
	case "hue": //0
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.Hue = -1
		} else {
			lamp.Hue = i
		}
	case "sat": //0
		if i, err := strconv.Atoi(res[1]); err != nil {
			lamp.Sat = -1
		} else {
			lamp.Sat = i
		}
	case "name":
		lamp.Name = res[1]
	}
	return nil
}
