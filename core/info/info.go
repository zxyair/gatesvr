package info

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const logo = `
                    ____  __  ________
                   / __ \/ / / / ____/	
                  / / / / / / / __/
                 / /_/ / /_/ / /___
                /_____/\____/_____/
`

const (
	boxWidth          = 56
	verticalBorder    = "|"
	horizontalBorder  = "─"
	leftTopBorder     = "┌"
	rightTopBorder    = "┐"
	leftBottomBorder  = "└"
	rightBottomBorder = "┘"
)

func PrintBoxInfo(name string, infos ...string) {
	fmt.Println(buildTopBorder(name))
	for _, info := range infos {
		fmt.Println(buildRowInfo(info))
	}
	fmt.Println(buildBottomBorder())
}

func buildRowInfo(info string) string {
	str := fmt.Sprintf("%s %s", verticalBorder, info)
	str += strings.Repeat(" ", boxWidth-utf8.RuneCountInString(str)-1)
	str += verticalBorder
	return str
}

func buildTopBorder(name ...string) string {
	full := boxWidth - strLen(leftTopBorder) - strLen(rightTopBorder) - strLen(name...)
	half := full / 2
	str := leftTopBorder
	str += strings.Repeat(horizontalBorder, half)
	if len(name) > 0 {
		str += name[0]
	}
	str += strings.Repeat(horizontalBorder, full-half)
	str += rightTopBorder
	return str
}

func buildBottomBorder() string {
	full := boxWidth - strLen(leftBottomBorder) - strLen(rightBottomBorder)
	str := leftBottomBorder
	str += strings.Repeat(horizontalBorder, full)
	str += rightBottomBorder
	return str
}

func strLen(str ...string) int {
	if len(str) > 0 {
		return utf8.RuneCountInString(str[0])
	} else {
		return 0
	}
}
