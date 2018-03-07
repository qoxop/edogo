package c2

import "../t1"
import "fmt"

func init(){
	t1.T = 2;
}
func Echo(){
	fmt.Println(t1.T);
}