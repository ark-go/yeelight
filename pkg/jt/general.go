package jt

import (
	"net"
	"reflect"
	"strings"
)

type LampRecv struct {
	Id     int
	Result []string
	Method string
	Params map[string]interface{} `json:"params"`
}

type LampStatus struct {
	Power  string //	on / off
	Bright string //	Процент яркости. Range 1 ~ 100
	// Ct          string //	Цветовая температура. Range 1700 ~ 6500(k)
	// Rgb         string //	Цвет. Range 1 ~ 16777215
	// Hue         string //	Оттенок. Range 0 ~ 359
	// Sat         string //	Насыщенность. Range 0 ~ 100
	// Color_mode  string //	1: rgb mode / 2: цветовая температура mode / 3: hsv mode
	// Flowing     string //	0: нет потока no flow is running / 1:color flow is running цветовой поток идет
	// Delayoff    string //	Оставшееся время таймера сна. Range 1 ~ 60 (minutes)
	// Flow_params string // Текущие параметры потока (имеет значение, только если значение параметра flowing равно 1)
	// Music_on    string //	1: Музыкальный режим включен / 0: Музыкальный режим выключен
	// Name        string //	Имя устройства, заданное командой «set_name».
}

func (l LampStatus) GetFields() []string {
	v := reflect.ValueOf(l)
	typeOfS := v.Type()
	var otv []string
	for i := 0; i < v.NumField(); i++ {
		otv = append(otv, strings.ToLower(typeOfS.Field(i).Name))
		//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	//log.Println("vvv", otv)
	return otv
}

//type lamp map[string]string
type LampSupport struct {
	get_prop      bool
	set_default   bool
	set_power     bool
	toggle        bool
	set_bright    bool
	set_scene     bool
	cron_add      bool
	cron_get      bool
	cron_del      bool
	start_cf      bool
	stop_cf       bool
	set_adjust    bool
	adjust_bright bool
	set_name      bool
	set_ct_abx    bool
	adjust_ct     bool
}

type Lamp struct {
	// что вернул запрос, 200 или при ошибке, весь ответ
	Result string // HTTP/1.1 200 OK
	//	Запрос не вернул 200,
	//	то что вернул смотреть в Result
	ErrorResult bool
	// не интересно, по умолчанию 1 час, время когда лампа пошлет данные
	CacheControl string // todo : max-age=3600
	// не используется
	Date int // skip
	// не используется
	Ext int // skip
	// адрес лампы из ответа
	Location string // yeelight://172.16.172.44:55443
	// IP адрес
	IPlamp string
	NTS    string // только при рекламной рассылке, лампа подключилась к сети..
	// сервер лампы
	Server string // POSIX UPnP/1.0 YGLC/1
	// уникальныи ИД лампы
	Id    string // хз что тут может быть,  0x0230000018430941 .
	Model string // ceila
	// версия прошивки
	FwVer string // 11  string/int?
	// поддерживаемые методы
	SupportAll string // get_prop set_default set_power toggle set_bright set_scene cron_add cron_get cron_del start_cf stop_cf set_adjust adjust_bright set_name set_ct_abx adjust_ct
	Power      string // off
	Bright     int    // 100
	ColorMode  int    // 2
	Ct         int    // 2700
	Rgb        int    // 0
	Hue        int    // 0
	Sat        int    // 0
	Name       string //
	//	mutex      sync.Mutex
	IsChanges bool
	// TCP Socket
	Conn net.Conn
}

// проверяет поддерживается ли метод
func (l *Lamp) IsSupport(method string) bool {
	return strings.Contains(l.SupportAll, method)
}

var Lamps map[string]*Lamp

func init() {
	Lamps = make(map[string]*Lamp)
}
