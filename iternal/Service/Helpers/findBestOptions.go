package Helpers

func FindBest(size int64) (int, int) {

	switch {
	case size >= 100000000:

		fileResult := size / 1000000

		x := 50
		ResultPart := int(fileResult) / x

		NumOfGoroutine := ResultPart + 1
		return x, NumOfGoroutine

	default:
		return 5, 20

	}
}
