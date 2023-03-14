package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"test-task-rit/internal/app/service"
	"test-task-rit/internal/app/types"
)

func main() {
	arg := os.Args
	if len(arg) == 1 {
		fmt.Println(errors.New("specify json file"))
		return
	}

	fileJSON, err := os.Open(arg[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileJSON.Close()

	byteFileJSON, err := ioutil.ReadAll(fileJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	actions := new(types.DataJSON)
	if err = json.Unmarshal(byteFileJSON, &actions); err != nil {
		fmt.Println(err)
		return
	}

	if err = service.NewActions(actions); err != nil {
		fmt.Println(err)
		return
	}

	newFileJSON, err := json.Marshal(actions)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = ioutil.WriteFile(fmt.Sprintf("resp-%s", arg[1]), newFileJSON, 0644); err != nil {
		fmt.Println(err)
		return
	}
}
