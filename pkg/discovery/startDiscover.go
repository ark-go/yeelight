package discovery

import (
	"log"
	"net"
	"regexp"
	"time"
)

// func StartDiscover1() string {
// 	pc, err := net.ListenPacket("udp4", ":8829") 1982
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer pc.Close()
// 	c := make(chan string) // Делает канал для связи
// 	go listenReports(pc)
// 	go sendRequest(pc)
// 	ticker := time.NewTicker(60 * time.Second)
// 	for {
// 		select { // Оператор select
// 		case gopherID := <-c: // Ждет, когда проснется гофер
// 			log.Println("gopher ", gopherID, " has finished sleeping")
// 		case <-ticker.C: // Ждет окончания времени
// 			log.Println("тик так")
// 			go sendRequest(pc)
// 			//return // Сдается и возвращается
// 		}
// 	}
// 	//	return "string(p)"
// }

// возвращает IP адрес из строки
func findIP(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + `\.` + numBlock + `\.` + numBlock + `\.` + numBlock

	regEx := regexp.MustCompile(regexPattern)
	//mask := net.CIDRMask(20, 32)
	//return strings.Trim(regEx.FindString(input), " ")
	return regEx.FindString(input)
}

func listenReports(pc *net.UDPConn) error {
	log.Println("Cтарт сервер UDP")
	for {
		p := make([]byte, 2048)
		log.Println("UDP Ждем данные...")
		//time.Sleep(3 * time.Second)
		n, addr, err := pc.ReadFromUDP(p)
		if err != nil {
			log.Println("Ошибка listen...", err)
			return err
		}
		log.Println("Пришло", n, " от: ", addr.String())

		ip := findIP(addr.String())
		if err := ParseReport(p[:n], ip); err != nil {
			return err
		}

		//	ConnectLamp() // идем на TCP подключения
	}
	//	return nil
}
func sendRequest(pc *net.UDPConn) {
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1982")
	if err != nil {
		panic(err)
	}
	payload := ""
	payload += "M-SEARCH * HTTP/1.1\r\n"
	payload += "HOST: \"239.255.255.250:1982\"\r\n"
	payload += "MAN: \"ssdp:discover\"\r\n"
	payload += "ST: wifi_bulb\r\n"
	//pc.RemoteAddr()
	//_, err = pc.WriteTo([]byte(payload), addr)
	_, err = pc.WriteToUDP([]byte(payload), addr)
	if err != nil {
		panic(err)
	}
	log.Println("Отослали запрос/данные ...")
}

func StartDiscover() string {
	// порт 1982 нужен для ожидания рекламных сообщений
	// а для поимки ответа на multicast порт как бы не важен?
	addr1, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1982") //"239.255.255.250:8829")
	if err != nil {
		panic(err)
	}
	//pc, err := net.ListenUDP("udp4", addr1)
	pc, err := net.ListenMulticastUDP("udp4", nil, addr1)
	//pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}

	defer pc.Close()
	c := make(chan string) // Делает канал для связи
	go listenReports(pc)
	go sendRequest(pc)
	//go ConnectLamp()

	ticker := time.NewTicker(10 * time.Minute)
	for {
		select { // Оператор select
		case gopherID := <-c: // Ждет, когда проснется гофер
			log.Println("gopher ", gopherID, " has finished sleeping")
		case <-ticker.C: // Ждет окончания времени
			log.Println("тик так")
			go sendRequest(pc)
			//return // Сдается и возвращается
		}
	}
	//	return "string(p)"
}
