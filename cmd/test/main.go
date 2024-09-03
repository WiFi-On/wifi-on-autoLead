// Тестовый файл, где можно полностью проверить работу всех функций

package main

import (
	"fmt"
	"strings"
	"time"
	"wifionAutolead/internal/services"
)

type Address struct {
	RegionID int
	CityID   string
	StreetID string
	HouseID  string
}

func main() {
	// Для примера используем захардкоженные значения
	regionID := 72
	cityName := "Тюмень"
	streetName := "Олимпийская"
	houseName := "47"

	// Замер времени для выполнения функции справочника города
	startCity := time.Now()
	cityMap, err := services.FetchAddressDirectory(regionID, 1, cityName)
	if err != nil {
		fmt.Println("Ошибка при получении данных для города:", err)
		return
	}
	fmt.Printf("Время выполнения справочника города: %v\n", time.Since(startCity))

	// Замер времени для выполнения функции справочника улицы
	startStreet := time.Now()
	streetMap, err := services.FetchAddressDirectory(regionID, 0, streetName)
	if err != nil {
		fmt.Println("Ошибка при получении данных для улицы:", err)
		return
	}
	fmt.Printf("Время выполнения справочника улицы: %v\n", time.Since(startStreet))

	// Замер времени на пересечение города и улиц
	startIntersection := time.Now()
	addresses := findIntersections(regionID, cityMap, streetMap)
	fmt.Printf("Время выполнения пересечения города и улиц: %v\n", time.Since(startIntersection))

	// Замер времени для получения справочника домов и их пересечения
	startHouseDirectory := time.Now()
	houseMap := fetchHouseDirectory(regionID, addresses, houseName)
	fmt.Printf("Время выполнения получения справочника домов: %v\n", time.Since(startHouseDirectory))

	// Замер времени на пересечение домов с найденными адресами
	startHouseIntersection := time.Now()
	addresses = intersectHousesWithAddresses(addresses, houseMap)
	fmt.Printf("Время выполнения пересечения домов с адресами: %v\n", time.Since(startHouseIntersection))

	// Вывод содержимого среза пересечений для проверки
	fmt.Println("Список пересечений:")
	for _, intersection := range addresses {
		fmt.Printf("RegionID: %d, CityID: %s, StreetID: %s, HouseID: %v\n",
			intersection.RegionID, intersection.CityID, intersection.StreetID, intersection.HouseID)
	}
}

// findIntersections ищет пересечения между городами и улицами
func findIntersections(regionID int, cityMap, streetMap map[string][]string) []Address {
	var addresses []Address
	for _, cityValues := range cityMap {
		for _, cityValue := range cityValues {
			cityAddrObjectId := strings.Split(cityValue, ";")[0]

			for _, streetValues := range streetMap {
				for _, streetValue := range streetValues {
					streetAddrObjectId := strings.Split(streetValue, ";")[0]
					streetParentId := strings.Split(streetValue, ";")[1]

					if cityAddrObjectId == streetParentId {
						address := Address{
							RegionID: regionID,
							CityID:   cityAddrObjectId,
							StreetID: streetAddrObjectId,
							HouseID:  "", // Пока пустое значение
						}
						addresses = append(addresses, address)
					}
				}
			}
		}
	}
	return addresses
}

// fetchHouseDirectory получает справочник домов для каждого найденного адреса
func fetchHouseDirectory(regionID int, addresses []Address, houseName string) map[string][]string {
	houseMap := make(map[string][]string)
	for _, address := range addresses {
		currentHouseMap, err := services.FetchAddressHouseInfo(regionID, address.StreetID, houseName)
		if err != nil {
			fmt.Println("Ошибка при получении данных для дома:", err)
			continue
		}

		for key, values := range currentHouseMap {
			houseMap[key] = values
		}
	}
	return houseMap
}

// intersectHousesWithAddresses сопоставляет дома с найденными адресами
func intersectHousesWithAddresses(addresses []Address, houseMap map[string][]string) []Address {
	for i, address := range addresses {
		houseKey := fmt.Sprintf("%d;%s", address.RegionID, address.StreetID)
		if houses, exists := houseMap[houseKey]; exists {
			for _, houseValue := range houses {
				addresses[i].HouseID = strings.Split(houseValue, ";")[0]
				break
			}
		}
	}
	return addresses
}
