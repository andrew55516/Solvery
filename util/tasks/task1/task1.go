package task1

func IsValidInput(input []int) bool {
	if len(input) == 0 {
		return false
	}
	min := input[0]
	max := input[0]

	for _, v := range input {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min >= 1 && max <= len(input)
}

func FindMissingNums(input []int) []int {
	result := make([]int, 0)
	exists := make([]bool, len(input))
	for _, v := range input {
		exists[v-1] = true
	}

	for i := 0; i < len(exists); i++ {
		if !exists[i] {
			result = append(result, i+1)
		}
	}
	return result
}
