package main

import "log"

func main() {
	s1 := []string{"1", "2", "3"}
	s2 := s1
	s2[0] = "-1"
	// 对s2的修改会对s1 造成修改
	//log.Println(s1)
	//log.Println(s2)

	// range 的时候修改v不会修改到原始
	for k, v := range s1 {
		v = "99"
		log.Println(k, v)
	}
	log.Println(s1)
}
