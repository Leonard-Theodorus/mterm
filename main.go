package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const configFileName = ".mterm"

func isHelpFlagPresent() bool {

    for _, v := range os.Args{
        if v == "-h" || v == "--h"{
            return true
        }
    }
    return false
}
func getLineNumber() int {
    dirErr := os.Chdir(getConfigFilePath())
    if dirErr != nil {
        //user home dir doesn't exists
        fmt.Println("Can't access directory: ", dirErr)
        os.Exit(1)
    }

    configFile, err := os.OpenFile(configFileName, os.O_RDONLY, 0400)
    //0400 -> linux file permission bits 
    if err != nil {
        //Opening file error || file doesn't exists
        fmt.Println("Config doesn't exist yet, run -i or --insert to create one")
        os.Exit(0)
    }
    line := 1
    scanner := bufio.NewScanner(configFile)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
        line++
    }
    return line
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

    configFile, err := os.OpenFile(configFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
    if err != nil {
        //Opening file error
        fmt.Println("Error opening file: ", err)
        os.Exit(1)
    }
    defer configFile.Close()
    currLineNumber := getLineNumber()
    contentToAppend := fmt.Sprintf("%d.%s\n", currLineNumber, newPath)

    configFile.WriteString(contentToAppend)
    fmt.Println("Sucessfully saved to terminal history")
}

func readConfigFile() []string{
    dirErr := os.Chdir(getConfigFilePath())
    if dirErr != nil {
        //user home dir doesn't exists
        fmt.Println("Can't access directory: ", dirErr)
        os.Exit(1)
    }

    configFile, err := os.OpenFile(configFileName, os.O_RDONLY, 0400)
    //0400 -> linux file permission bits 
    if err != nil {
        //Opening file error || file doesn't exists
        fmt.Println("Config doesn't exist yet, run -i or --insert to create one")
        os.Exit(0)
    }
    defer configFile.Close()
    scanner := bufio.NewScanner(configFile)
    var cfgFileContent []string
    for scanner.Scan() {
        cfgFileContent = append(cfgFileContent, string(scanner.Text()))
    }
    return cfgFileContent
}

func printSavedPaths() {
    savedPaths := readConfigFile()
    for _, content := range savedPaths{
        fmt.Println(content)
    }
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

func getSavedPath(idx int) string {
    savedPaths := readConfigFile()
    lineNumber := getLineNumber()
    if idx > lineNumber {
        fmt.Println("Path with that number doesn't exist. Run -p or --print to see saved paths")
        os.Exit(0)
    }
    path := savedPaths[idx - 1]
    path = strings.Split(path, ".")[1]
    return path
}

func jumpToPath(idx int) {
    //TODO: Find ways to engineer changing paths
    path := getSavedPath(idx)
    cmd := exec.Command("cd", path)
    cmdOut, err := cmd.Output()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    } else {
        fmt.Printf("%s\n", cmdOut)
    }
}

func main () {

    flag.BoolFunc("h", "Help commands for memoterm", func(s string) error {
        return nil // no-error
    })
    pathFlag := flag.String("i", "Esc", "Insert new path")
    readFlag := flag.Bool("p", false, "Print out saved paths")
    jumpFlag := flag.Int("j", 0, "Jump to saved path")
    for idx, arg := range os.Args {

        if arg == "--insert"{
            os.Args[idx] = "-i"
            if idx + 1 < len(os.Args) {
                *pathFlag = os.Args[idx + 1]
            } else {
                fmt.Println("--insert expects an argument <path>")
                os.Exit(0)
            }
        } else if arg == "-i" {
            if idx + 1 >= len(os.Args) {
                fmt.Printf("%s expects an argument <path>\n", arg)
                os.Exit(0)
            }
        } else if arg == "--i" {
            fmt.Println("Maybe you meant -i or --insert ?")
            os.Exit(0)
        }

        if arg == "--print" {
            *readFlag = true
        } else if arg == "--p" {
            fmt.Println("Maybe you meant -p or --print ?")
            os.Exit(0)
        }
    }
    flag.Parse()

    if isHelpFlagPresent() {
        //TODO: Update the command lists
        fmt.Println("Usage: mterm [-i | --insert] <path>")
        os.Exit(0)
    }

    arg := os.Args[1]
    switch arg {
    case "-i" :
        //function save
        insertNewPath(*pathFlag)
    case "-p":
        printSavedPaths()
    case "-j":
        jumpToPath(*jumpFlag)

    default:
        fmt.Println("No such command, run mterm -h for help")
        os.Exit(0)
    }

}
