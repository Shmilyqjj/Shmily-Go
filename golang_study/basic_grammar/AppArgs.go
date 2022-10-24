package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	file_abs_path := args[0]
	if len(args) == 4 {
		arg1 := args[1]
		arg2 := args[2]
		arg3 := args[3]
		fmt.Println(fmt.Sprintf("file_abs_path=%s,arg1=%s,arg2=%s,arg3=%s", file_abs_path, arg1, arg2, arg3))
	} else {
		fmt.Println("file_abs_path=" + file_abs_path)
	}

}
