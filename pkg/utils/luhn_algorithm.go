package utils

func IsValidLuhn(number int64) bool {
	if number <= 0 {
		return false
	}

	sum := 0
	alt := false

	for number > 0 {
		digit := int(number % 10)

		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alt = !alt
		number /= 10
	}

	return sum%10 == 0
}
