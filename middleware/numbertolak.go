package middleware

import (
	"fmt"
	"strings"
)

var units = []string{"ສູນ", "ໜຶ່ງ", "ສອງ", "ສາມ", "ສີ່", "ຫ້າ", "ຫົກ", "ເຈັດ", "ແປດ", "ເກົ້າ", "ສິບ"}
var positions = []string{"", "ສິບ", "ຮ້ອຍ", "ພັນ", "ໝື່ນ", "ແສນ", "ລ້ານ", "ໂກດ", "ກື", "ຕື້", "ຕິ່ວ", "ຕົງ"}

func cleanInput(amount string) string {
	amount = strings.ReplaceAll(amount, "ກີບ", "")
	amount = strings.ReplaceAll(amount, ",", "")
	amount = strings.ReplaceAll(amount, " ", "")

	return amount
}

func convertIntegerPart(intPart string) string {
	var result string

	for i := 0; i < len(intPart); i++ {
		digit := int(intPart[i] - '0')
		fmt.Println(digit, len(intPart)-2)
		if digit != 0 {
			if i == len(intPart)-1 && digit == 1 {
				result += "ເອັດ"
			} else if i == len(intPart)-2 && digit == 2 {
				result += "ຊາວ"
			} else if i == len(intPart)-2 && digit == 1 {
				result += ""
			} else {
				result += units[digit]
			}

			if i == len(intPart)-2 && digit == 2 {

			} else {
				result += positions[len(intPart)-i-1]
			}
			// if len(intPart) >= 8 {
			// 	result += "ລ້ານ"
			// }

		}

	}

	if result == "ເອັດ" {
		return "ໜຶ່ງ"
	}

	return result
}

func convertDecimalPart(decimalPart string) string {
	var result string
	for i := 0; i < len(decimalPart); i++ {
		digit := int(decimalPart[i] - '0')
		if digit != 0 {
			if i == len(decimalPart)-1 && digit == 1 {
				result += "ເອັດ"
			} else if i == len(decimalPart)-2 && digit == 2 {
				result += "ຊາວ"
			} else if i == len(decimalPart)-2 && digit == 1 {
				result += ""
			} else {
				result += units[digit]
			}
			// result += positions[len(decimalPart)-i-1]
			if i == len(decimalPart)-2 && digit == 2 {

			} else {
				result += positions[len(decimalPart)-i-1]
			}
		}
	}

	if result == "" {
		return "ຖ້ວນ"
	}

	if result == "ເອັດ" {
		return "ໜຶ່ງກີບ"
	}

	return result + "ກີບ"
}

func ConvertToText(amount string) string {

	parts := strings.Split(cleanInput(amount), ".")
	if len(parts) > 2 {
		return "ຂໍ້ມູນບໍ່ຖືກຕ້ອງ,ຈຳນວນທີ່ຖືກຕ້ອງເປັນແບບຕົວຢ່າງເຊັ່ນ 1000.012 ຫຼື 10000"
	}
	if len(parts) <= 1 {
		amount = amount + ".00"
		parts = strings.Split(cleanInput(amount), ".")
	}
	// fmt.Println("here")
	if len(parts[0]) > 12 || len(parts[1]) > 12 {
		return "ບໍສາມາດນັບໄດ້ ຈຳນວນທີ່ນັບໄດ້ສູງສຸດແມ່ນ 99 ຕົງ ຫຼືເລກສູງສຸດ 12 ຫຼັກເທົ່ານັ້ນ"
	}
	intText := convertIntegerPart(parts[0])
	decText := convertDecimalPart(parts[1])

	if intText == "" && decText == "ຖ້ວນ" {
		return "ສູນກີບຖ້ວນ"
	}

	if intText == "" {
		intText = "ສູນ"
	}
	wordRepresentation := ""
	if len(parts) == 1 {
		wordRepresentation = intText + "ກີບ" + decText
	} else {
		if decText == "ຖ້ວນ" {
			wordRepresentation = intText + "ກີບ" + decText
		} else {
			wordRepresentation = intText + "ຈຸດ" + decText
		}

	}

	return wordRepresentation
}
