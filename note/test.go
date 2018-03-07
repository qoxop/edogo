package main
import (
    "fmt"
    "flag"
)

const  repInfo = "the url of your git repository, make sure you can clone the repository";
const tInfo = "Time interval of git fetch";

var rep string = "";
var t int = 10;
func autoFetch(rep string,interval int){
	if rep != ""{
		
	}
}
func clone(url string) string{
	if rep != "" {

	}else{
		return "clone faild";
	}
}
func main(){
	flag.StringVar(&rep,"rep","",repInfo);
	flag.IntVar(&t,"T",10,tInfo);
    flag.Parse()
	fmt.Printf("you git repository is  : %s\n",rep);
	fmt.Printf("Time interval of git fetch is  : %d\n",t);
	fmt.Println("cloning...");
	clone(rep);
	fmt.Println("...");
	autoFetch(rep,t)
}