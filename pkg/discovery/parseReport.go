package discovery

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/ark-go/yeelight/pkg/jt"
)

func ParseReport(p []byte, ip string) error {

	r := bytes.NewReader(p)
	scanner := bufio.NewScanner(r) // Будем читать по строкам
	lamp := &jt.Lamp{
		IPlamp: ip,
	}
	var ii int = 0
	for scanner.Scan() {
		ii++
		txt := scanner.Text()
		if ii == 1 {
			if strings.Contains(txt, "NOTIFY") {
				continue // это рекламное сообщение, наример лампа вошла в  сеть
			}
			if !strings.Contains(txt, "200") { // в первой строке должно быть 200
				lamp.ErrorResult = true
				lamp.Result = string(p) // весь ошибочный ответ сохраним для анализа
				return fmt.Errorf("запрос к лампе вернул ошибку")
			} else {
				continue // пропускаяем первую строку
			}

		}
		//log.Println(strconv.Itoa(ii) + " : " + txt)
		parseStrToStruct(txt, lamp)

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if newId := copyNewLampData(lamp); newId != "" {
		ConnectLamp(newId)
		SendCommand(newId)
	}
	// todo удалить,   просмотр всех ламп
	for key, value := range jt.Lamps {
		log.Printf("%s Значение: %+v\n", key, value)
	}
	return nil
}
