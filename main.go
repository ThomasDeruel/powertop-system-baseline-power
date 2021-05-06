package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TIMER = 20
const DEFAULT_FILE = "powertop.csv"
const FIELDS_PER_RECORD = -1

type Powertop struct {
	system_baseline_power float64
}

func main() {
	// generate csv file
	// default name: powertop.csv
	generateFile()
	var timer = time.Now()
	readfile, err := readFile(DEFAULT_FILE)
	if err != nil {
		fmt.Println(err)
	}
	// find globales values with regexep: here "system baseline power"...
	re := regexp.MustCompile(`The [a-zA-Z ]+:[\s]*([0-9]*\.?[0-9]*)([\s\w]+)`)
	var values []float64
	for _, line := range readfile {
		find_sub_match := re.FindSubmatch([]byte(line[0]))
		if len(find_sub_match) > 0 {
			// [1] -> value (on string)
			// [2] -> format (exemple: mw, kw, uw, etc.)
			values = append(values, findAndConvertPrefix(string(find_sub_match[2]), string(find_sub_match[1])))
		}
	}
	powertop := Powertop{
		values[0],
	}
	fmt.Printf("system baseline power %v W\nat: %v.\n", powertop.system_baseline_power, timer.Format(time.RFC822))
}
func generateFile() {
	fmt.Printf("generate powertop in %v SEC.\n----------------------------\n", TIMER)
	// use exec.Command to run in 'sudo'
	out, err := exec.Command("/bin/sh", "-c", "sudo ./script.sh").Output()
	if err != nil {
		fmt.Print("error:%v", err)
	}
	fmt.Println(string(out))
}
func readFile(file string) ([][]string, error) {
	// open csv file
	csvFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	// read file
	r := csv.NewReader(csvFile)
	// https://golang.org/pkg/encoding/csv/
	r.FieldsPerRecord = FIELDS_PER_RECORD
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func findAndConvertPrefix(format string, val string) float64 {
	valconv, err := strconv.ParseFloat(val, 64)
	if err != nil {
		fmt.Println(err)
	}
	if strings.ContainsAny(format, "u") {
		valconv *= math.Pow10(-6)
	}
	if strings.ContainsAny(format, "m") {
		valconv *= math.Pow10(-3)
	}
	if strings.ContainsAny(format, "k") {
		valconv *= math.Pow10(3)
	}
	return valconv
}
