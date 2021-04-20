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


var _mtx sync.Mutex


func SpawnBackup() {
	_mtx = sync.Mutex{}
	err := exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go").Run()
	//LINUX:
	// err := exec.Command("gnome-terminal", "-x", "go", "run", "main.go").Run()
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
	json.Unmarshal(byteValue, &backupElevator)

	return backupElevator
}


