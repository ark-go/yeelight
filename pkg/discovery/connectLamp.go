package discovery

import (
	"log"
	"net"

	"github.com/ark-go/yeelight/pkg/jt"
)

// func connectLamps() {
// 	for _, lamp := range jt.Lamps {
// 		ConnectLamp(lamp)
// 	}
// }

func ConnectLamp(newId string) {
	for _, lamp := range jt.Lamps {
		if lamp.ErrorResult {
			continue
		}
		if lamp.Location == "" {
			continue
		}
		if lamp.Conn != nil {
			continue // todo test connect

		}
		conn, err := net.Dial("tcp", lamp.Location)
		if err != nil {
			log.Println(err)
		}
		log.Println("Подключились:", lamp.Location)
		lamp.Conn = conn

	}
}
