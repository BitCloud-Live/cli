package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	mainColor      = color.BgBlue
	secColor       = color.BgYellow
	mainTxt        = color.New(color.FgHiWhite, color.BgBlack, color.Bold)
	mainTxtSprint  = mainTxt.SprintFunc()
	mainTxtColor   = mainTxt.SprintfFunc()
	secTxtColor    = color.New(color.FgWhite, color.BgBlack, color.Faint).SprintfFunc()
	mainTxtBlink   = color.New(color.FgHiWhite, color.BgBlack, color.Bold, color.BlinkSlow).SprintfFunc()
	mainTitle      = color.New(color.FgHiWhite, mainColor, color.Bold)
	mainTitleColor = mainTitle.SprintfFunc()
	mainTitlePrint = mainTitle.SprintFunc()
	secTitleColor  = color.New(color.FgBlack, secColor, color.Bold).SprintfFunc()
	secTitleBlink  = color.New(color.FgBlack, secColor, color.Bold, color.BlinkSlow).SprintfFunc()
	whiteSpace     = "    "
	whiteSpaceDash = "  - "
)

func mainTxtPrintln(body ...interface{})  { fmt.Println(mainTxtSprint(body...)) }
func colorfulPrintln(body ...interface{}) { fmt.Println(mainTitlePrint(body...)) }
func colorfulPrint(body ...interface{})   { fmt.Print(mainTitlePrint(body...)) }

func printKeyVal(space, k, v string) {
	fmt.Printf("%s:%s  \r\n",
		mainTxtColor(" %s%s ", space, k),
		secTxtColor(" %s ", v))
}

func printTitleByStatus(space, title, status string) {
	if strings.ToLower(status) == "down" {
		fmt.Printf("%s%s:%s\r\n",
			space,
			secTitleColor(" %s ", title),
			secTitleBlink(" # [Unreachable]   "))
	} else {
		fmt.Printf("%s%s:\r\n",
			space,
			mainTitleColor(" %s ", title))
	}
}

func printTitle(space, title string) {
	fmt.Printf("%s%s:\r\n",
		space,
		mainTitleColor(" %s ", title))
}
