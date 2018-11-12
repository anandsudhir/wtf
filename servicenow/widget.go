package servicenow

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
)

const HelpText = `
 Keyboard commands for ServiceNow:

   /: Show/hide this help window
   j: Select the next change task in the list
   k: Select the previous change task in the list
   r: Refresh the data

   arrow down: Select the next change task in the list
   arrow up:   Select the previous change task in the list

   return: Open the selected change task in a browser
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.TextWidget

	changeTasks []ChangeTask
	selected    int
}

func NewWidget(app *tview.Application, pages *tview.Pages) *Widget {
	widget := Widget{
		HelpfulWidget: wtf.NewHelpfulWidget(app, pages, HelpText),
		TextWidget:    wtf.NewTextWidget(app, "ServiceNow", "servicenow", true),
	}

	widget.HelpfulWidget.SetView(widget.View)
	widget.unselect()

	widget.View.SetScrollable(true)
	widget.View.SetRegions(true)
	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	changeTasks, err := widget.GetChangeTasksAssignedToMe()

	if err != nil {
		widget.View.SetWrap(true)
		widget.View.SetTitle(widget.Name)
		widget.View.SetText(err.Error())
	} else {
		widget.changeTasks = changeTasks
	}

	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	widget.View.SetWrap(false)

	widget.View.Clear()
	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - changetasks", widget.Name)))
	widget.View.SetText(widget.contentFrom(widget.changeTasks))
	widget.View.Highlight(strconv.Itoa(widget.selected)).ScrollToHighlight()
}

func (widget *Widget) contentFrom(changeTasks []ChangeTask) string {
	if widget.changeTasks == nil {
		return "No assigned change tasks to display"
	}

	var str string
	str = "[white::b] Change\t\t\tTask\t\t\tDue date\t\t\tShort description\n"
	for idx, changeTask := range changeTasks {
		str = str + fmt.Sprintf(
			`["%d"][""][%s] [greenyellow]%s`+"\t"+`[%s]%s`+"\t"+`%s`+"\t"+`%s`,
			idx,
			widget.rowColor(idx),
			changeTask.changeNumber,
			widget.rowColor(idx),
			changeTask.taskNumber,
			changeTask.expectedEnd,
			changeTask.taskDescription,
		)

		str = str + "\n"
	}

	return str
}

func (widget *Widget) rowColor(idx int) string {
	if widget.View.HasFocus() && (idx == widget.selected) {
		return wtf.DefaultFocussedRowColor()
	}

	return wtf.RowColor("servicenow", idx)
}

func (widget *Widget) next() {
	widget.selected++
	if widget.changeTasks != nil && widget.selected >= len(widget.changeTasks) {
		widget.selected = 0
	}

	widget.display()
}

func (widget *Widget) prev() {
	widget.selected--
	if widget.selected < 0 && widget.changeTasks != nil {
		widget.selected = len(widget.changeTasks) - 1
	}

	widget.display()
}

func (widget *Widget) openChangeTask() {
	sel := widget.selected
	if sel >= 0 && widget.changeTasks != nil && sel < len(widget.changeTasks) {
		changeTask := &widget.changeTasks[widget.selected]
		wtf.OpenFile(changeTask.url)
	}
}

func (widget *Widget) unselect() {
	widget.selected = -1
	widget.display()
}

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case "/":
		widget.ShowHelp()
	case "j":
		widget.next()
		return nil
	case "k":
		widget.prev()
		return nil
	case "r":
		widget.Refresh()
		return nil
	}

	switch event.Key() {
	case tcell.KeyDown:
		widget.next()
		return nil
	case tcell.KeyEnter:
		widget.openChangeTask()
		return nil
	case tcell.KeyEsc:
		widget.unselect()
		return event
	case tcell.KeyUp:
		widget.prev()
		return nil
	default:
		return event
	}
}
