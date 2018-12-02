package go_site

import (
    "fmt"
    "github.com/zhanben/go_site/tool/config"
)

func main(){
    //Read config file
    err := config.ParseConfig()
    if err != nil {
        panic(fmt.Errorf("Failed to read config file: %s \n", err))
    }
    //Init log

}