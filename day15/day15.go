package main

// callback should treat its argument as immutable
// If mutation is needed use slices.CLone

func iteratePartitions(sum, k int, accumulated []int, callback func([]int)) {
	if k == 1 {
		callback(append(accumulated, sum))
		return
	}

	for i := range sum + 1 {
		iteratePartitions(sum-i, k-1, append(accumulated, i), callback)
	}
}
