package servicenow

import (
	"fmt"
	"os"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/senorprogrammer/wtf/logger"
	"github.com/senorprogrammer/wtf/wtf"
)

func (widget *Widget) GetChangeTasksAssignedToMe() ([]ChangeTask, error) {
	var changeTasks []ChangeTask
	baseUrl := wtf.Config.UString("wtf.mods.servicenow.baseUrl")

	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
	)

	if err := driver.Start(); err != nil {
		logger.Log(fmt.Sprintf("Failed to start driver: %s", err.Error()))
	}

	page, err := driver.NewPage()
	if err != nil {
		logger.Log(fmt.Sprint("Failed to open page: %s", err.Error()))
		return nil, err
	}

	if err := page.Navigate(baseUrl); err != nil {
		logger.Log(fmt.Sprintf("Failed to navigate: %s", err.Error()))
		return nil, err
	}

	// Login
	err = page.FindByID(`username`).SendKeys(wtf.Config.UString("wtf.mods.servicenow.username"))
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to submit username: %s", err.Error()))
		return nil, err
	}
	err = page.FindByID(`password`).SendKeys(widget.password())
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to submit password: %s", err.Error()))
		return nil, err
	}
	err = page.FindByID(`submit`).Click()
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to submit sign-in form: %s", err.Error()))
		return nil, err
	}

	if err := page.Navigate(baseUrl + changeTasksAssignedToMe); err != nil {
		logger.Log(fmt.Sprintf("Failed to navigate: %s", err.Error()))
	}

	rows, _ := page.AllByXPath(`//*[@id="change_task"]/table/tbody/tr`).Elements()
	for _, row := range rows {
		var changeTask ChangeTask

		columns, _ := row.GetElements(api.Selector{Using: "css selector", Value: "td"})
		if len(columns) == 1 {
			logger.Log("No change tasks assigned!")
			return changeTasks, nil
		}

		changeTask.taskNumber, err = columns[2].GetText()
		changeTask.url = baseUrl + changeTaskUrl + changeTask.taskNumber
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to get taskNumber: %s", err.Error()))
			return nil, err
		}
		changeTask.expectedStart, err = columns[5].GetText()
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to get expectedStart: %s", err.Error()))
			return nil, err
		}
		changeTask.expectedEnd, err = columns[6].GetText()
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to get expectedEnd: %s", err.Error()))
			return nil, err
		}
		changeTask.taskDescription, err = columns[7].GetText()
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to get taskDescription: %s", err.Error()))
			return nil, err
		}

		changeTasks = append(changeTasks, changeTask)
	}

	if err := driver.Stop(); err != nil {
		logger.Log(fmt.Sprintf("Failed to close pages and stop WebDriver: %s", err.Error()))
		return nil, err
	}

	return changeTasks, nil
}

/* -------------------- Unexported Functions -------------------- */

const (
	changeTasksAssignedToMe = `change_task_list.do?sysparm_query=assignment_group=javascript:getMyGroups()^active=true^change_request.stateIN30,40^EQ&sysparm_clear_stack=true`
	changeTaskUrl           = `nav_to.do?uri=task.do?sysparm_query=number=`
)

func (widget *Widget) password() string {
	return wtf.Config.UString(
		"wtf.mods.servicenow.password",
		os.Getenv("WTF_SERVICENOW_PASSWORD"),
	)
}
