package main

import "fmt"

func main() {
	explodeString()
}

//测试make和new的差异
func test1() {
	var arr = new([]int) // new返回一个指针
	//arr = append(arr, 1) 不行 因为需要接受切片
	*arr = append(*arr, 1)

	//或者这样
	arr2 := *arr
	arr2 = append(arr2, 2)

	// 或者这样
	arr3 := make([]int, 10)
	arr3 = append(arr3, 1)
}

//测试合并两个切片
func mergeSlice() {
	slice1 := []int{1, 2, 3}
	slice2 := make([]int, 2)
	slice2 = append(slice2, 4)
	slice2 = append(slice2, 5)
	slice2 = append(slice2, 5)

	//使用...展开切片作为参数
	slice1 = append(slice1, slice2...)
	fmt.Println(slice2)
	println(slice2)
}

//测试切割字符串
func explodeString() {
	s := "这是abcd123"
	fmt.Println(s[2:]) //切割中文了

}
