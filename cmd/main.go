package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"wifionAutolead/pkg/eissd"
)

// Айди ростелекома = 52
func main() {
	numbers := []string{
		"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
		"31", "32", "33", "34", "35", "36", "37", "38", "39", "40",
		"41", "42", "43", "44", "45", "46", "47", "48", "49", "50",
		"51", "52", "53", "54", "55", "56", "57", "58", "59", "60",
		"61", "62", "63", "64", "65", "66", "67", "68", "69", "70",
		"71", "72", "73", "74", "75", "76", "77", "78", "79", "80",
		"81", "82", "83", "84", "85", "86", "87", "88", "89", "90",
		"91", "92", "93", "94", "95", "96", "97", "98", "99",
	}
	

	godotenv.Load("../common/conf/.env")

	eissdClient := eissd.NewEISSDpars("https://eissd.roseltorg.ru", "../common/certs/krivoshein.crt", "../common/certs/krivoshein.key")

	districts, err := eissdClient.GetDistrictsOrAddresses("02", 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(districts)
	fmt.Println("Дистрикты получил")

	eissdDB, err := eissd.NewDB("../common/db/eissd.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Файл создан либо уже был создан")

	err = eissdDB.CreateDistrictsTable()
	if err != nil {
		panic(err)
	}
	fmt.Println("Таблица districts создана либо уже была создана")

	err = eissdDB.CreateStreetsTable()
	if err != nil {
		panic(err)
	}
	fmt.Println("Таблица streets создана либо уже была создана")

	err = eissdDB.CreateHousesTable()
	if err != nil {
		panic(err)
	}
	fmt.Println("Таблица houses создана либо уже была создана")

	for j := 0; j < len(numbers); j++ {
		
	}
	
	// for _, reg := range numbers {
	// 	districts, err := eissdClient.GetDistrictsOrAddresses(reg, 1)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(districts)
	// 	fmt.Println("Дистрикты получил")

	// 	for j, district := range districts {
	// 		err = eissdDB.AddDistrict(district.AddrObjectId, district.RegionId, district.NameAddrObject, district.AbbrNameObject)
	// 		if err != nil {
	// 			fmt.Println("Номер итерации = " , j)
	// 			fmt.Println("Ошибка = ", err)
	// 			fmt.Println("Дистрикт = ", district)
	// 			fmt.Println("Регион = ", reg)
	// 			panic("ОШИБКА")
	// 		}
	// 		fmt.Printf("Дистрикт #%d добавлен\n", j)
	// 	}
		
	// }

	for _, reg := range numbers {
		districts, err := eissdClient.GetDistrictsOrAddresses(reg, 0)
		if err != nil {
			panic(err)
		}
		fmt.Println(districts)
		fmt.Println("Улицы получил")

		for j, district := range districts {
			err = eissdDB.AddStreet(district.AddrObjectId, district.RegionId, district.NameAddrObject, district.AbbrNameObject, district.ParentId)
			if err != nil {
				fmt.Println("Номер итерации = " , j)
				fmt.Println("Ошибка = ", err)
				fmt.Println("Улица = ", district)
				fmt.Println("Регион = ", reg)
				panic("ОШИБКА")
			}
			fmt.Printf("Улица #%d добавлена\n", j)
		}
		
	}
}