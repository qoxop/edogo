
package main
import (
	"os/exec"
	"log"
	"fmt"
)

func main(){
	cmd := exec.Command("sleep","1");
	out,err := cmd.CombinedOutput();
	cmd.Run();
	

	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(out,cmd.Path,cmd.Env)
}