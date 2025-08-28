package aptosindexergrpcgo

// Example demonstrates how to use filterFast function
//
//	numbers := []int{1, 2, 3, 4, 5, 6}
//	filtered := filterFast(numbers, func(n int) bool {
//		return n%2 == 0
//	})
//	// filtered: [2, 4, 6]
//
// Note: This modifies the original slice. If you need to preserve
// the original, make a copy first:
//
//	original := []int{1, 2, 3, 4, 5, 6}
//	copy := make([]int, len(original))
//	copy(copy, original)
//	filtered := filterFast(copy, func(n int) bool {
//		return n%2 == 0
//	})
func FilterFast[T any](slice []T, predicate func(T) bool) []T {
	j := 0
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			slice[j] = slice[i]
			j++
		}
	}
	return slice[:j]
}
