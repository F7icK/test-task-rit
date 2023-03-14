package service

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"test-task-rit/internal/app/types"
)

func NewActions(actions *types.DataJSON) error {
	for idx, action := range actions.Actions {
		switch action.Type {
		case types.ActionType:
			switch action.Name {
			case types.CreateFile:
				if len(action.Parameters) < 1 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				fileName, err := serviceCreateFile(action.Parameters[0])
				if err != nil {
					return err
				}

				action.Result = fileName
				for i := 1; idx+i < len(actions.Actions); i++ {
					if actions.Actions[idx+i].Type != types.ConditionActionType {
						actions.Actions[idx+i].Parameters = append(actions.Actions[idx+i].Parameters, fileName)
						break
					}
				}
			case types.ChangeFileName:
				if len(action.Parameters) < 2 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				if err := serviceChangeFileName(action.Parameters[1], action.Parameters[0]); err != nil {
					return err
				}

				action.Result = action.Parameters[1]
				for i := 1; idx+i < len(actions.Actions); i++ {
					if actions.Actions[idx+i].Type != types.ConditionActionType {
						actions.Actions[idx+i].Parameters = append(actions.Actions[idx+i].Parameters, action.Parameters[0])
						break
					}
				}
			case types.RemoveFile:
				if len(action.Parameters) < 1 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				if err := serviceRemoveFile(action.Parameters[0]); err != nil {
					return err
				}

				action.Result = "File deleted"
				for i := 1; idx+i < len(actions.Actions); i++ {
					if actions.Actions[idx+i].Type != types.ConditionActionType {
						actions.Actions[idx+i].Parameters = append(actions.Actions[idx+i].Parameters, action.Parameters[0])
						break
					}
				}
			case types.GetFileCreationTime:
				if len(action.Parameters) < 1 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				timeCreationFile, err := serviceGetFileCreationTime(action.Parameters[0])
				if err != nil {
					return err
				}

				action.Result = timeCreationFile.Format("2006-01-02 15:04:05")
				for i := 1; idx+i < len(actions.Actions); i++ {
					if actions.Actions[idx+i].Type != types.ConditionActionType {
						actions.Actions[idx+i].Parameters = append(actions.Actions[idx+i].Parameters, action.Result, action.Parameters[0])
						break
					}
				}
			case types.WritingLineToFile:
				if len(action.Parameters) < 2 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				if err := serviceWritingLineToFile(action.Parameters[0], action.Parameters[1]); err != nil {
					return err
				}

				action.Result = fmt.Sprintf("line \"%s\" recorded in %s", action.Parameters[0], action.Parameters[1])
				for i := 1; idx+i < len(actions.Actions); i++ {
					if actions.Actions[idx+i].Type != types.ConditionActionType {
						actions.Actions[idx+i].Parameters = append(actions.Actions[idx+i].Parameters, action.Parameters[1])
						break
					}
				}
			}
		case types.ConditionType:
			switch action.Name {
			case types.CreatedAtTimeCondition:
				if idx+2 >= len(actions.Actions) {
					return errors.New(fmt.Sprintf("no implementation after condition with output to true and false in %s", action.Name))
				}

				if len(action.Parameters) < 3 {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				dateDefault, err := time.Parse("2006-01-02 15:04:05", action.Parameters[0])
				if err != nil {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				dateCreation, err := time.Parse("2006-01-02 15:04:05", action.Parameters[1])
				if err != nil {
					return errors.New(fmt.Sprintf("missing parameters in %s", action.Name))
				}

				// Предполагается что в JSON после условия первым будет действие если true, а вторым false
				if dateDefault.After(dateCreation) {
					action.Result = "true"
					actions.Actions[idx+1].Type = types.ActionType
					actions.Actions[idx+1].Parameters = append(actions.Actions[idx+1].Parameters, action.Parameters[2])
				} else {
					action.Result = "false"
					actions.Actions[idx+2].Type = types.ActionType
					actions.Actions[idx+2].Parameters = append(actions.Actions[idx+2].Parameters, action.Parameters[2])
				}
			}
		case types.ConditionActionType:
			continue
		default:
			return errors.New(fmt.Sprintf("failed to determine action type in %s", action.Name))
		}
	}

	return nil
}

func serviceCreateFile(fileName string) (string, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	return fileInfo.Name(), nil
}

func serviceChangeFileName(oldFileName, newFileName string) error {
	if err := os.Rename(oldFileName, newFileName); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func serviceRemoveFile(fileName string) error {
	if err := os.Remove(fileName); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func serviceGetFileCreationTime(fileName string) (*time.Time, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	ctime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)

	return &ctime, nil
}

func serviceWritingLineToFile(line, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(line); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
