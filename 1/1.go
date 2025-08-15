package main

import "fmt"

func main(){
	var w,s,i int
	fmt.Scanf("%d %d %d", &w, &s, &i)
	wsSum := 24 - (w + s) + i
	fmt.Println(wsSum)
}