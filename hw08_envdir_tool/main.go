package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	dirName := args[1]
	env, err := ReadDir(dirName) // Читаем env-переменные из указанной директории
	if err != nil {
		fmt.Printf("unable to read env dir: %s", err)
		return
	}
	commandWithArgs := args[2:]           // Ожидаем, что 3 аргумент - имя команды, остальные - аргументы команды
	os.Exit(RunCmd(commandWithArgs, env)) // Завершаем команду с кодом, полученным из дочерней команды
}
