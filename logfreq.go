package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var firstTime time.Time
var lastTime time.Time

//var lastT, last2T *string
var tf *string
var delta *time.Duration

var empty string = ""

func main() {
	f := flag.String("f", "", "Time format - %T %C %H %M %S %m %d %y %b."+
		"\n\tMore precise format can be given via -re and -tf")
	regex := flag.String("re", "", "Regex to extract date and time")
	tf = flag.String("tf", "", "Time format")

	flag.Parse()

	if *f != "" {
		*regex = "(" + *f + ")"
		*tf = *f
	}

	if *regex == "" {
		fmt.Fprintf(os.Stderr, "Time regex must not be empty (-re <regexp)\n")
		os.Exit(1)
	}
	if *tf == "" {
		fmt.Fprintf(os.Stderr, "Time format must not be empty (-tf <format)\n")
		os.Exit(1)
	}
	tfmt := "2006/01/02 15:04"
	if strings.Index(*tf, "%T") != -1 || strings.Index(*tf, "%S") != -1 {
		tfmt = "2006/01/02 15:04:05"
	}

	*regex = strings.Replace(*regex, "%Z", `[A-Z][A-Z][A-Z][A-Z]?`, -1)
	*regex = strings.Replace(*regex, "%z", `[+-]\d\d\d\d`, -1)
	*regex = strings.Replace(*regex, "%b", `[A-Z][a-z][a-z]`, -1)
	*regex = strings.Replace(*regex, "%T", `\d\d:\d\d:\d\d`, -1)
	*regex = strings.Replace(*regex, "%C", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%H", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%M", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%S", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%m", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%d", `\d\d`, -1)
	*regex = strings.Replace(*regex, "%Y", `\d\d\d\d`, -1)

	*tf = strings.Replace(*tf, "%Z", "MST", -1)
	*tf = strings.Replace(*tf, "%z", "-0700", -1)
	*tf = strings.Replace(*tf, "%b", "Jan", -1)
	*tf = strings.Replace(*tf, "%T", "15:04:05", -1)
	*tf = strings.Replace(*tf, "%C", `06`, -1)
	*tf = strings.Replace(*tf, "%H", `15`, -1)
	*tf = strings.Replace(*tf, "%M", `04`, -1)
	*tf = strings.Replace(*tf, "%S", `05`, -1)
	*tf = strings.Replace(*tf, "%m", `01`, -1)
	*tf = strings.Replace(*tf, "%d", `02`, -1)
	*tf = strings.Replace(*tf, "%Y", `2006`, -1)

	re := regexp.MustCompile(*regex)

	ret := 0
	files := []string{"-"}
	if len(flag.Args()) > 0 {
		files = flag.Args()
	}

	var src io.Reader

	freq := map[string]int{}
	for _, file := range files {
		if file == "-" {
			src = os.Stdin
		} else {
			f, err := os.Open(file)
			if err != nil {
				log.Println(err)
				ret = 1
				continue
			}
			src = f
		}

		lnum := 0
		r := bufio.NewScanner(src)
		for r.Scan() {
			line := r.Text()
			lnum++
			m := re.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			tstamp := m[1]
			t, err := time.Parse(*tf, tstamp)
			if err != nil {
				log.Fatal(err)
			}
			normalizedTstamp := t.Format(tfmt)
			//fmt.Fprintf(os.Stderr, "tstamp %v => %v\n", tstamp, normalizedTstamp)
			freq[normalizedTstamp] = freq[normalizedTstamp] + 1
		}
	}
	keys := []string{}
	for k, _ := range freq {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("date time frequency cumulative")
	cumulative := 0
	for _, k := range keys {
		cumulative += freq[k]
		fmt.Println(k, freq[k], cumulative)
	}
	os.Exit(ret)
}
