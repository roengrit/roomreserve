package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

//File64Encode -
func File64Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff), nil
}

//ParseDateTime -
func ParseDateTime(dateStr string, timeStr string) (time.Time, error) {
	fmt.Println(dateStr)
	fmt.Println(timeStr)
	sp := strings.Split(dateStr, "-")
	if dateStr != "" {
		retDate, errDate := time.Parse(time.RFC3339, sp[2]+"-"+sp[1]+"-"+sp[0]+"T"+timeStr+":00+00:00")
		if errDate != nil {
			if strings.Contains(errDate.Error(), "month out of range") {
				retDate, errDate := time.Parse(time.RFC3339, sp[2]+"-"+sp[0]+"-"+sp[1]+"T"+timeStr+":00+00:00")
				return retDate, errDate
			}
		}
		return retDate, errDate
	}
	return time.Now(), errors.New("วันที่ไม่ถูกต้อง")
}

//CreateDateTimeFromString -
func CreateDateTimeFromString(dateTimeStr string) (dateTime time.Time, err error) {
	sp := strings.Split(dateTimeStr, " ")
	if len(sp) != 2 {
		err = errors.New("รุปแบบวันที่/เวลาไม่ถูกต้อง")
		return
	}
	spDate := strings.Split(sp[0], "-")
	if len(spDate) != 3 {
		err = errors.New("รุปแบบวันที่ไม่ถูกต้อง")
		return
	}
	var year, month, day, hh, mm int
	year, err = strconv.Atoi(spDate[0])
	if err == nil {
		month, err = strconv.Atoi(spDate[1])
	}
	if err == nil {
		day, err = strconv.Atoi(spDate[2])
	}
	spTime := strings.Split(sp[1], ":")
	if len(spTime) < 2 {
		err = errors.New("รุปแบบเวลาไม่ถูกต้อง")
		return
	}
	if err == nil {
		hh, err = strconv.Atoi(spTime[0])
	}
	if err == nil {
		mm, err = strconv.Atoi(spTime[1])
	}
	return time.Date(int(year), time.Month(month), int(day), hh, mm, 0, 0, time.Now().Location()), nil
}

//AddToThaiYear -
func AddToThaiYear(dateTime time.Time) time.Time {
	return dateTime.AddDate(543, 0, 0)
}
