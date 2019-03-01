package main

func MapTo(arr []int, convert func(item, index int) string) (res []string) {
	for i := 0; i < len(arr); i++ {
		res = append(res, convert(arr[i], i))
	}
	return
}

func Convert(arr []int) []string {
	var result []string
	for _, v := range arr {
		var str string
		switch v {
		case 1:
			str = "one"
		case 2:
			str = "two"
		case 3:
			str = "three"
		case 4:
			str = "four"
		case 5:
			str = "five"
		case 6:
			str = "six"
		case 7:
			str = "seven"
		case 8:
			str = "eight"
		case 9:
			str = "nine"
		default:
			str = "unknown"
		}
		result = append(result, str)
	}
	return result
}

func main() {
}
