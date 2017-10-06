package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

var version string

func getPort() int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "9005"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	return port
}

func getDB() *Store {
	dbPath := os.Getenv("DB")

	if dbPath == "" {
		dbPath = "dev.db"
	}

	store, err := NewStore("dev.db")

	if err != nil {
		panic(err)
	}

	return store
}

func writePid(pidFile string) {
	pid := syscall.Getpid()
	err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		panic(err)
	}
}

func getIsDebug() bool {
	return os.Getenv("DEBUG") != ""
}

func main() {
	writePid("pagaraprata.pid")
	NewServer(getIsDebug(), getPort(), getDB()).Run()
}
