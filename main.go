package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
)

func main() {
	verbose := flag.Bool("verbose", false, "print information besides the info, including errors")
	flag.Parse()

	// The the percentage of CPU usage.
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		if *verbose {
			fmt.Println("reading /proc/loadavg:", err)
		}
		return
	}
	parts := bytes.Split(data, []byte(" "))
	f, err := strconv.ParseFloat(string(parts[0]), 64)
	if err != nil && *verbose {
		fmt.Printf("converting '%v': %v\n", string(parts[0]), err)
	}
	cpu := f / float64(runtime.NumCPU()) * 100

	// Get the percentage of memory usage.
	data, err = ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		if *verbose {
			fmt.Println("reading /proc/meminfo:", err)
		}
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	mem := (1.0 - parse(lines[2], verbose)/parse(lines[0], verbose)) * 100

	fmt.Printf("(C: %0.2f%% M: %0.2f%%)\n", cpu, mem)
}

func parse(line []byte, verbose *bool) float64 {
	parts := bytes.Split(line, []byte(" "))
	kb := parts[len(parts)-2]
	f, err := strconv.ParseFloat(string(kb), 64)
	if err != nil && *verbose {
		fmt.Printf("converting '%v': %v\n", string(kb), err)
	}
	return f
}
