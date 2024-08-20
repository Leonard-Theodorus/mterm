package main

import (
	"flag"
	"fmt"
	"os"
    "os/user"
)

const configFileName = ".mterm"

func isHelpPresent() bool {

    for _, v := range os.Args{
        if v == "-h" || v == "--h"{
            return true
        }
    }
    return false
}

func insertNewPath(newPath string) {
    //check for file; create if none
    cfgFilePath := getConfigFilePath()
    dirErr := os.Chdir(cfgFilePath)
    if dirErr != nil {
        //user home dir doesn't exists
        fmt.Println("Can't access directory: ", dirErr)
        os.Exit(1)
    }
    //TODO: number for pahts
    contentToAppend := fmt.Sprintf("%d.%s", 1, newPath)

    configFile, err := os.OpenFile(configFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    fmt.Println(configFile)
    if err != nil {
        //Opening file error
        fmt.Println("Error opening file: ", err)
        os.Exit(1)
    }
    defer configFile.Close()

    configFile.WriteString(contentToAppend)
    fmt.Println("Sucessfully saved to terminal history")
}

func getConfigFilePath() string {
    currentUser, err := user.Current()
    if err != nil {
        fmt.Println("Can't get current user : ", err)
        os.Exit(1) // some error happened
    }
    configFilePath := currentUser.HomeDir
    return configFilePath
}

func main () {
    flag.BoolFunc("h", "Help commands for memoterm", func(s string) error {
        return nil // no-error
    })
    pathFlag := flag.String("i", "Esc", "Insert new path")
    flag.Parse()
    if isHelpPresent() {
        fmt.Println("Usage: mterm [-i | --insert] <path>")
        os.Exit(0)
    }
    fmt.Println(*pathFlag)
    arg := os.Args[1]
    switch arg {
    case "-i" :
        //function save
        insertNewPath(*pathFlag)
    case "--insert":
        //function save
        insertNewPath(*pathFlag)

    default:
        fmt.Println("No such command, run mterm -h for help")
        os.Exit(0)
    }

}
