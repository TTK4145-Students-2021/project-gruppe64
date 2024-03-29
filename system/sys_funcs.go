package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var _mtx sync.Mutex // Mutex to protect system log

func SpawnBackup() {
	_mtx = sync.Mutex{}
	err := exec.Command("gnome-terminal", "-x", "go", "run", "main.go").Run()
	if err != nil {
		fmt.Println(err)
	}
}

func IsBackup() bool{
	if _, err := os.Stat("system/primary_doc"+strconv.Itoa(ElevatorID)+".txt"); os.IsNotExist(err){
		return false
	} else {
		return true
	}
}

func CheckPrimaryExistence(activateAsPrimary chan<- bool) {
	for {
		data1, _ := ioutil.ReadFile("system/primary_doc"+strconv.Itoa(ElevatorID)+".txt")
		num1, _ := strconv.Atoi(string(data1))
		time.Sleep(3*time.Second)
		data2, _ := ioutil.ReadFile("system/primary_doc"+strconv.Itoa(ElevatorID)+".txt")
		num2, _ := strconv.Atoi(string(data2))
		if num1 == num2 {
			activateAsPrimary <- true
			break
		}
	}
}

// GOROUTINE, main-initiated
// While alive, the primary will document its existence by counting and documenting
// the number to file.
func PrimaryDocumentation() {
	docNum := 0
	for {
		_ = ioutil.WriteFile("system/primary_doc"+strconv.Itoa(ElevatorID)+".txt", []byte(strconv.FormatInt(int64(docNum), 10)), 0644)
		time.Sleep(1*time.Second)
		docNum += 1
	}
}

func MakeBackupFile() {
	file, err := os.Create("system/primary_doc"+strconv.Itoa(ElevatorID)+".txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
}

func LogElevator(elevInfo Elevator) {
	_mtx.Lock()
	defer _mtx.Unlock()
	jsonFile, _ := json.MarshalIndent(elevInfo, "", " ")
	err := ioutil.WriteFile("system/sys_log"+strconv.Itoa(ElevatorID)+".json", jsonFile, 0644)
	if err !=nil {
		fmt.Println(err)
	}
}

func GetLoggedElevator() Elevator{
	_mtx.Lock()
	defer _mtx.Unlock()
	jsonFile, err := os.Open("system/sys_log"+strconv.Itoa(ElevatorID)+".json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var backupElevator Elevator
	_ = json.Unmarshal(byteValue, &backupElevator)
	return backupElevator
}