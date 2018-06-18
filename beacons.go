package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	maidenhead "github.com/pd0mz/go-maidenhead"
)

var bands = []int{
	144,
	222,
	420,
	902,
	1240,
	2300,
	2390,
	3300,
	5650,
	10000,
	24000,
	47000,
	76000,
	122250,
	134000,
	241000,
	275000,
}

func main() {
	resp, err := http.Get("http://www.newsvhf.com/beacons2.html")
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"band", "freq", "call", "grid", "lat", "lon", "state", "city", "comments"})
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 {
			continue
		}
		if !unicode.IsNumber(rune(t[0])) {
			continue
		}
		f := strings.Fields(t)
		freq, call, grid := f[0], f[1], f[2]
		f = f[3:len(f)]
		if len(grid) < 2 {
			call += " " + grid
			grid = f[0]
			f = f[1:len(f)]
		}
		var band string
		fnum, err := strconv.Atoi(freq[0:strings.Index(freq, ".")])
		if err == nil {
			for _, i := range bands {
				if fnum < i {
					break
				}
				band = fmt.Sprintf("%d", i)
			}
		}
		p, err := maidenhead.ParseLocator(grid)
		if err != nil {
			log.Printf("parsing %q: %v", grid, err)
		}
		var state string
		if len(f[0]) == 2 {
			state = f[0]
			f = f[1:len(f)]
		}
		var city string
		if len(t) > 30 && strings.TrimSpace(t[27:30]) != "" {
			if t[26] == ' ' {
				city = strings.TrimSpace(t[27:40])
				if len(t) > 40 {
					f = []string{t[40:len(t)]}
				} else {
					f = nil
				}
			} else {
				city = f[0]
				f = f[1:len(f)]
			}
		}
		comments := strings.Join(f, " ")
		w.Write([]string{band, freq, call, grid, fmt.Sprintf("%.6f", p.Latitude), fmt.Sprintf("%.6f", p.Longitude), state, city, comments})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
