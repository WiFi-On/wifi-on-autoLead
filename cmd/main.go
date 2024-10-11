package main

import (
	"log"
	"net/http"
	"os"
	controllers "wifionAutolead/internal/controllers/rest"
	"wifionAutolead/internal/repository"
	"wifionAutolead/internal/routes"
	"wifionAutolead/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../common/conf/.env"); err != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", err)
	}

	bitrixURLToken := os.Getenv("BITRIX_URL_TOKEN")
	dadataAPIKey := os.Getenv("DADATA_API_KEY")

	// Подключение к базе данных
	connStr := "user=admin password=admin dbname=eissd sslmode=disable"

	// Инициализация репозиториев
	eissdRepo, err := repository.NewDB(connStr)
	if err != nil {
		log.Fatalf("Не удалось инициализировать репозитории: %v", err)
	}

	// Инициализация логгера
	loggerService, err := services.NewLoggerService("../common/log/eissd.log")
	if err != nil {
		log.Fatalf("Не удалось инициализировать сервис логирования: %v", err)
	}
	defer func() {
		if err := loggerService.Close(); err != nil {
			log.Printf("Ошибка при закрытии сервиса логирования: %v", err)
		}
	}()

	// Инициализация сервисов
	dadataService := services.NewDadata(dadataAPIKey)
	eissdService := services.NewEISSD(eissdRepo, "../common/certs/krivoshein_dev.crt", "../common/certs/krivoshein_dev.key", "https://mpz-rc.rt.ru/xmlInteface", dadataService, loggerService)
	bitrixService := services.NewBitrix(bitrixURLToken)

	// Инициализация хэндлеров
	controllers := controllers.NewEISSDController(eissdService, bitrixService, dadataService)

	// Создание маршрутизатора
	mux := http.NewServeMux()

	// Регистрация маршрутов
	routes.RegisterRoutes(mux, controllers)

	log.Println("Запуск сервера на порту 3020...")
	loggerService.Log("Сервер запущен на порту 3020") // Логируем запуск сервера

	log.Fatal(http.ListenAndServe(":3020", mux))
}

// func main(){
// 	// bitrixClient := services.NewBitrix("https://on-wifi.bitrix24.ru/rest/11940/s1l6kfr94x1vuck4/")
// 	dadataClient := dadata.NewDadata("71378de14318d10009285e018aedbfe5a353bb5a")
// 	repo, err := repository.NewDB("user=admin password=admin dbname=eissd sslmode=disable")
// 	if err != nil {
// 		fmt.Errorf("не удалось инициализировать репозитории: %v", err)
// 	}
// 	eissdClient := services.NewEISSD(repo, "../common/certs/krivoshein_dev.crt", "../common/certs/krivoshein_dev.key", "https://mpz-rc.rt.ru/xmlInteface")

// 	// leads, err := bitrixClient.GetDealsOnProviders("52")
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// addressBitrix := leads[0].Address
// 	// addressArray := strings.Split(addressBitrix, "|")
// 	// fmt.Println(addressArray[0])

// 	address := "Тюмень Широтная 105 кв 14"

// 	dadataInfo, err := dadataClient.GetInfoOnAddress(address)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(dadataInfo.Suggestions[0].Data)
// 	data := dadataInfo.Suggestions[0].Data
// 	regionNumber := string([]rune(data.RegionKladrID)[0:2])
// 	thv, err := eissdClient.CheckTHV(regionNumber, data.Area, data.City, data.Settlement, data.Street, data.House, data.Block, data.Flat)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(thv)

// }

// func main() {
// 	clientDB, err := eissd.NewDB("user=admin password=admin dbname=eissd sslmode=disable")
// 	if err != nil {
// 		fmt.Errorf("не удалось инициализировать репозитории: %v", err)
// 	}
// 	defer clientDB.Close()

// 	// Открываем файл
// 	file, err := os.Open("tariffsOnRegion.txt")
// 	if err != nil {
// 		fmt.Println("Ошибка при открытии файла:", err)
// 		return
// 	}
// 	defer file.Close() // Закрываем файл в конце

// 	// Создаем сканер для построчного чтения файла
// 	scanner := bufio.NewScanner(file)

// 	// Читаем строки из файла
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		parts := strings.Split(line, "~")
// 		if parts[1] == "Технологический IPTV" {
// 			clientDB.AddTariffForRegion(parts[2], 0, 0, 0, 0)
// 		}
// 	}

// 	// Проверяем на наличие ошибок при чтении
// 	if err := scanner.Err(); err != nil {
// 		fmt.Println("Ошибка при чтении файла:", err)
// 	}
// }

// func main() {
// 	clientDB, err := eissd.NewDB("user=admin password=admin dbname=eissd sslmode=disable")
// 	clientEISSD := eissd.NewEISSDpars("../common/certs/krivoshein.crt", "../common/certs/krivoshein.key", "https://mpz.rt.ru/xmlInteface")
// 	if err != nil {
// 		fmt.Errorf("Не удалось инициализировать репозитории: %v", err)
// 	}
// 	defer clientDB.Close()

// 	tariffs, err := clientEISSD.GetTarrifsOnRegion("72")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(tariffs)
// }