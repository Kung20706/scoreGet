package main

import "fmt"

// binarySearch 使用二分搜索算法在已排序的切片中查找目標值
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	// 當左邊界小於等於右邊界時執行搜索
	for left <= right {
		// 計算中間索引
		mid := left + (right-left)/2

		// 如果中間值等於目標值，則返回中間索引
		if arr[mid] == target {
			return mid
		}

		// 如果中間值小於目標值，縮小搜索範圍到右半部分
		if arr[mid] < target {
			left = mid + 1
		} else { // 如果中間值大於目標值，縮小搜索範圍到左半部分
			right = mid - 1
		}
	}

	// 如果未找到目標值，返回 -1 表示未找到
	return -1
}

func main() {
	// 有序的整數切片
	sortedArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 目標值
	target := 7

	// 在切片中搜索目標值
	result := binarySearch(sortedArray, target)

	// 打印結果
	if result != -1 {
		fmt.Printf("目標值 %d 在切片中的索引為 %d\n", target, result)
	} else {
		fmt.Printf("目標值 %d 未在切片中找到\n", target)
	}
}
