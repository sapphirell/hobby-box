package main

import "log"

func main() {
	a()
}

func a() {
	log.Println(10 / 4)   // 2
	log.Println(10.0 / 4) //2.5
	var f float64
	f = 10 / 4 //2
	log.Println(f)
	f = 10.0 / 4 //2.5
	log.Println(f)

	n1 := 10
	n2 := n1
	n2 = 11
	log.Println(n1, n2)

}
