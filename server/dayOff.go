package main

import (
	"fmt"
	"strconv"
	"strings"
)

type dayOff struct {
	name      string
	surname   string
	startDate date
	endDate   date
}

type date struct {
	day   int
	month int
	year  int
}

func getDayOffRequest(input string) dayOff {
	var myDayOff dayOff
	str := strings.Split(input, " ")
	myDayOff.name = str[0]
	myDayOff.surname = str[1]
	myDayOff.startDate = getDate(str[2])
	myDayOff.endDate = getDate(str[3])
	return myDayOff
}
func getDate(input string) date {
	var myDate date
	str := strings.Split(input, ".")
	myDate.day, _ = strconv.Atoi(str[0])
	myDate.month, _ = strconv.Atoi(str[1])
	myDate.year, _ = strconv.Atoi(str[2])

	return myDate
}
func (d *date) toString() string {
	str := fmt.Sprintf("%d.%d.%d", d.day, d.month, d.year)
	return str
}
func (d *dayOff) toString() string {
	str := fmt.Sprintf(""+
		"| İsim  | Soyisim  | Izin baslangic  | Izin bitis |\n"+
		"| :-----------------| ----------------- | ----------------- | -----------------: |\n"+
		"|         %s        |         %s        |        %s         |         %s         |\n", d.name, d.surname, d.startDate.toString(), d.endDate.toString())
	return str
}
