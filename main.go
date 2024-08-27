package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
    "strconv"
)

const configFileName = ".mterm"
const scriptFileName = "mterm.sh"

func isHelpFlagPresent() bool {

    for _, v := range os.Args{
        if v == "-h" || v == "--h"{
            return true
        }
    }
    return false
}

func checkScript() bool {
    _ , err := os.OpenFile(scriptFileName,os.O_CREATE |  os.O_EXCL, 0755)
    if err != nil {
        //File already exists
        return true
    }
    return false
}

func createCdScript() {
    const scriptContent = `
    #!/usr/bin/env bash
    function mterm() {
        local path
        read -r path
        path=$(echo "$path")

        if [ -d "$path" ]; then
        cd "$path" || exit
        echo "Changed directory to $(pwd)"
        else
        echo "Error: $path is not a valid directory."
        exit 1
        fi
    }
    `
    script, wErr := os.OpenFile(scriptFileName, os.O_CREATE | os.O_RDWR, 0755)
    if wErr != nil {
        //File already exists
        fmt.Println("Error when creating script :", wErr)
        os.Exit(0)
    }
    defer script.Close()
    script.WriteString(scriptContent)
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
    if idx >= lineNumber {
        fmt.Println("Path with that number doesn't exist. Run -p or --print to see saved paths")
        os.Exit(0)
    }
    path := savedPaths[idx - 1]
    path = strings.Split(path, ".")[1]
    return path
}

func jumpToPath(idx int) {
    path := getSavedPath(idx)
    fmt.Println(path)
}

func main () {

    if !checkScript() {
        createCdScript()
        fmt.Println("Config created, run -h for list of commands")
    } else {
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
                os.Args[idx] = "-p"
                *readFlag = true
            } else if arg == "--p" {
                fmt.Println("Maybe you meant -p or --print ?")
                os.Exit(0)
            }

            if arg == "--jump" {
                os.Args[idx] = "-j"
                if idx + 1 < len(os.Args) {
                    parsedArg, convErr := strconv.Atoi(os.Args[idx + 1])
                    if convErr != nil {
                        fmt.Println("Error parsing argument : ", convErr)
                        os.Exit(0)
                    }
                    *jumpFlag = parsedArg
                } else {
                    fmt.Println("--jump expects an argument <index>")
                    os.Exit(0)
                }
            } else if arg == "-j" {
                if idx + 1 >= len(os.Args) {
                    fmt.Printf("%s expects an argument <index>\n", arg)
                    os.Exit(0)
                }
            } else if arg == "--j" {
                fmt.Println("Maybe you meant -j or --jump ?")
                os.Exit(0)
            }
        }
        flag.Parse()

        if isHelpFlagPresent() || len(os.Args) < 2 {
            //TODO: Update the command lists
            fmt.Println("Usage: mterm [-i | --insert] <path>")
            fmt.Printf("%6s", "mterm [-p || --print]\n")
            fmt.Printf("%6s", "mterm [-j || --jump] | mterm\n")
            os.Exit(0)
        }

        arg := os.Args[1]
        switch arg {
        case "-i" :
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
}
