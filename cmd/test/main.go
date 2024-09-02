// Тестовый файл, где можно полностью проверить работу всех функций

package main

import (
	"fmt"
	"wifionAutolead/internal/services"
)

func main() {
	// Заполняем адрес
	var region int = 72
	var city string = "Тюмень"
	var street string = "Олимпийская"
	var house string = "47"
	var svcClassId int = 2 // Пример идентификатора класса услуги

	// Список для хранения итоговых результатов
	var finalResults []struct {
		RegionId int
		CityId   string
		StreetId string
		HouseId  string
	}

	// Получаем адреса районов/городов
	district_addresses, err := services.FetchAddressDirectory(region, 1, city)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Получаем адреса улиц
	street_addresses, err := services.FetchAddressDirectory(region, 0, street)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Поиск объектов, где ParentId из street_addresses совпадает с AddrObjectId из district_addresses
	for _, district := range district_addresses {
		for _, street := range street_addresses {
			if street.ParentId == district.AddrObjectId {
				// Для каждого найденного совпадения ищем дома
				houses, err := services.FetchAddressHouseInfo(region, street.AddrObjectId, house)
				if err != nil {
					fmt.Println("Ошибка при получении домов:", err)
					continue
				}

				// Вывод найденных домов и сохранение результата
				if len(houses) > 0 {
					fmt.Println("Найденные дома на улице", street.NameAddrObject)
					for _, house := range houses {
						fmt.Printf("Дом: %s, ID дома: %s\n", house.House, house.HouseId)

						// Сохранение результатов для вывода
						finalResults = append(finalResults, struct {
							RegionId int
							CityId   string
							StreetId string
							HouseId  string
						}{
							RegionId: house.RegionId,
							CityId:   district.AddrObjectId,
							StreetId: street.AddrObjectId,
							HouseId:  house.HouseId,
						})
					}
				} else {
					fmt.Println("Дома на улице", street.NameAddrObject, "не найдены.")
				}
			}
		}
	}

	// Вывод найденных объектов с RegionID, DistrictID, StreetID, HouseID
	if len(finalResults) > 0 {
		fmt.Println("Итоговые результаты:")
		for _, result := range finalResults {
			fmt.Printf("RegionID: %d, CityID: %s, StreetID: %s, HouseID: %s\n",
				result.RegionId, result.CityId, result.StreetId, result.HouseId)
		}

		// Проверка возможности подключения для каждого найденного результата
		for _, result := range finalResults {
			response, message, err := services.CheckConnectionPossibilityAgent(result.RegionId, result.CityId, result.StreetId, result.HouseId, svcClassId)
			if err != nil {
				fmt.Println("Ошибка при проверке возможности подключения:", err)
			} else if response == 0 {
				fmt.Printf("Подключение возможно для дома с HouseID: %s на улице с StreetID: %s.\n", result.HouseId, result.StreetId)
			} else {
				fmt.Printf("Подключение невозможно для дома с HouseID: %s на улице с StreetID: %s. Причина: %s\n", result.HouseId, result.StreetId, message)
			}
		}

	} else {
		fmt.Println("Не найдено объектов с совпадающими RegionID, CityID, StreetID и HouseID.")
	}
}
