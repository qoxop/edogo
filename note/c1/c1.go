package c1

import "../t1"
import "fmt"

func init(){
	t1.T = 1;
}
func Echo(){
	fmt.Println(t1.T);
}