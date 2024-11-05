package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func dashesTopBot(w string) {
	l := len(w)
	var d string

	for i := 0; i < l; i++ {
		d += "-"
	}

	fmt.Println(d)
	fmt.Println(w)
	fmt.Println(d)
}

func dashesBot(w string) {
	l := len(w)
	var d string

	for i := 0; i < l; i++ {
		d += "-"
	}

	fmt.Println(w)
	fmt.Println(d)
}

func getTaskName() string {
	var task string
	fmt.Println("TASK:")
	for {
		fmt.Scanln(&task)
		task = strings.ToLower(task)
		if task != "bd" && task != "fso" &&
			task != "build" {
			fmt.Println("Enter correct task fuckboi")
			continue
		}
		break
	}
	return task
}

func getTaskHours() string {
	var hours uint8
	fmt.Println("HOURS:")
	for {
		fmt.Scanln(&hours)
		if hours == 0 || hours > 14 {
			fmt.Println("LYING FUCK")
			continue
		}
		break
	}
	return fmt.Sprintf("%d", hours)
}

type TaskData struct {
	Name  string
	Hours string
	Day   string
	Date  string
}

func createTaskData() TaskData {
	return TaskData{
		Name:  getTaskName(),
		Hours: getTaskHours(),
		Day:   fmt.Sprintf("%v", time.Now().Weekday()),
		Date:  time.Now().Format("01-02-2006"),
	}
}

func createCSV() {
	file, err := os.OpenFile("data.csv", os.O_CREATE, 0664)
	if err != nil {
		log.Fatal("Can't create data.csv bruh")
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatal("Couldn't get file stats from data.csv bruh")
	}

	if stats.Size() == 0 {
		file.WriteString("task,hours,weekday,date\n")
	}
}

func getLast7Days() {
	file, err := os.OpenFile("data.csv", os.O_RDONLY, 0664)
	if err != nil {
		log.Fatal("Can't read data.csv bruh")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read()

	dateHoursMap := make(map[string]int)
	datesSlice := make([]string, 0, 8)

	var lastDate string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Can't read record bruh")
		}

		hours, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal("Couldn't convert string to integer bruh")
		}

		dateHoursMap[record[3]] += hours
		if lastDate != record[3] {
			datesSlice = append(datesSlice, record[3])
		}
		lastDate = record[3]
		if len(datesSlice) == 8 {
			delete(dateHoursMap, datesSlice[0])
			copy(datesSlice[0:], datesSlice[1:])
			datesSlice = datesSlice[:len(datesSlice)-1]
		}
	}

	for _, date := range datesSlice {
		dashesBot(date + ": " + strconv.Itoa(dateHoursMap[date]))
	}
}

func getStats() {
	file, err := os.OpenFile("data.csv", os.O_RDONLY, 0664)
	if err != nil {
		log.Fatal("Can't read data.csv bruh")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read()

	var totalHours int
	taskHoursMap := make(map[string]int)
	dateHoursMap := make(map[string]int)
	var todaysHours int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Can't read record bruh")
		}

		hours, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal("Couldn't convert string to int bruh")
		}

		totalHours += hours
		taskHoursMap[record[0]] += hours
		dateHoursMap[record[3]] += hours

		if time.Now().Format("01-02-2006") == record[3] {
			todaysHours += hours
		}
	}

	dashesTopBot("Total Hours: " + strconv.Itoa(totalHours))

	dashesBot(fmt.Sprintf("Today's Hours: %d", todaysHours))

	dayTotal := 0
	days := 0
	for _, v := range dateHoursMap {
		days += 1
		dayTotal += v
	}
	var avg float64
	if len(dateHoursMap) != 0 {
		avg = float64(dayTotal / days)
	}
	dashesBot(fmt.Sprintf("Average Hours Per Day: %.2f", avg))

	for k, v := range taskHoursMap {
		dashesBot(k + ": " + strconv.Itoa(v))
	}
}

func addRecord() {
	file, err := os.OpenFile("data.csv", os.O_APPEND, 0664)
	if err != nil {
		log.Fatal("Can't open data.csv in write mode")
	}
	defer file.Close()

	getStats()
	var yn string
	for {
		fmt.Println("Add Record? n: exit --- d: last 7 days data")
		fmt.Scanln(&yn)
		yn = strings.ToLower(yn)
		if yn == "n" {
			break
		}
		if yn == "d" {
			getLast7Days()
			continue
		}
		rec := createTaskData()
		file.WriteString(fmt.Sprintf(
			"%s,%s,%s,%s\n", rec.Name, rec.Hours, rec.Day, rec.Date,
		))
	}
	getStats()
}

func main() {
	createCSV()
	addRecord()
}
