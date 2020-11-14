package array

func newBitsWords(nums []int32) (int32, []uint64) {

	n := int32(0)
	if len(nums) > 0 {
		n = nums[len(nums)-1] + 1
	}

	nWords := (n + 63) >> 6
	words := make([]uint64, nWords)

	for _, i := range nums {
		iWord := i >> 6
		i = i & 63
		words[iWord] |= 1 << uint(i)
	}
	return n, words
}
