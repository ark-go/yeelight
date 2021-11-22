package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ark-go/yeelight/pkg/jt"
)

type readResult struct {
	readBytes int
	resByte   []byte
	err       error
}

func SendCommand(id string) {
	if lamp, ok := jt.Lamps[id]; ok {
		ls := jt.LampStatus{}

		i, _ := strconv.ParseInt(strings.TrimLeft(id, "0x"), 16, 64) // 64/32
		is := strconv.Itoa(int(i))
		fmt.Printf("%s    %d\n", id, i)
		//	httpRequest := `{"id":1,"method":"get_prop","params":["power","bright","ct","rgb","hue","sat","color_mode","flowing","delayoff","flow_params","music_on","name"]}` + "\r\n"
		// Внимательно!! собираем строку
		httpRequest := `{"id":` + is + `,"method":"get_prop","params":["` + strings.Join(ls.GetFields(), "\",\"") + "\"]}" + "\r\n"
		log.Println("vvv", httpRequest)
		if n, err := lamp.Conn.Write([]byte(httpRequest)); err != nil {
			log.Println(err)
			return
		} else {
			log.Println("Отправлена команда..", n, "   ", lamp.Location)
			//io.Copy(os.Stdout, lamp.Conn)
			go func() {
				for {
					//p := make([]byte, 2048)
					log.Println("TCP Ждем данные...")
					//	n, err := lamp.Conn.Read(p)
					ctx, cancel := context.WithCancel(context.Background())
					//cancel()
					_ = cancel
					x, err := readReader(ctx, lamp.Conn.Read) // todo тестирование контекста - убрать
					if err != nil {
						log.Println("Ошибка readTCP...", err.Error())
						break
					}
					n = x.readBytes
					p := x.resByte
					if x.err != nil {
						log.Println("ошибка", x.err.Error())
					}
					log.Println("прочитали :", n, string(p))
					if n > 0 {
						formatStatus(p[:n])
					}
				}
			}()
		}

	}
}

func formatStatus(p []byte) {
	status := jt.LampRecv{}
	err := json.Unmarshal(p, &status)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	if status.Params != nil {
		log.Printf("Пришло TCP: %+v", status.Params)
	} else if status.Result != nil {
		log.Printf("Ответ TCP: %+v", status.Result)
	}
	//strings.Join([]string{"1", "ll2", "kk3", "mm4"}, ", ")
}

// принимает функцию Read,
func readReader(ctx context.Context, read func(b []byte) (n int, err error)) (*readResult, error) {
	resultCh := make(chan *readResult)
	go func() {
		if err := ctx.Err(); err != nil { // контекст уже отменен, не надо ничего читать
			log.Println("Отменен контекст.")
			return
		}
		b := make([]byte, 2048)
		n, err := read(b) // ожидающая функция, здесь поток ждет данные
		//hex.DecodeString
		//time.Sleep(5 * time.Second)
		select {
		case <-ctx.Done(): // возможно нам уже не нужен результат
		case resultCh <- &readResult{readBytes: n, resByte: b[:n], err: err}: // получим данные Readers
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // выход с ошибкой
	case result := <-resultCh: // тут уже есть результат
		return result, nil
	}
}

/*
power		on / off
bright		Процент яркости. Range 1 ~ 100
ct 			Цветовая температура. Range 1700 ~ 6500(k)
rgb			Цвет. Range 1 ~ 16777215
hue 		Оттенок. Range 0 ~ 359
sat			Насыщенность. Range 0 ~ 100
color_mode	1: rgb mode / 2: цветовая температура mode / 3: hsv mode
flowing 	0: нет потока no flow is running / 1:color flow is running цветовой поток идет
delayoff	Оставшееся время таймера сна. Range 1 ~ 60 (minutes)
flow_params Текущие параметры потока (имеет значение, только если значение параметра flowing равно 1)
music_on 	1: Музыкальный режим включен / 0: Музыкальный режим выключен
name 		Имя устройства, заданное командой «set_name».
*/
//"power","bright","ct","rgb","hue","sat","color_mode","flowing","delayoff","flow_params","music_on","name"
