package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getMetrics() []string {
	startMetrics := []string{"Mean", "Median", "Mode", "SD"}
	var res []string
	if len(os.Args) != 1 {
		for _, val := range os.Args[1:] {
			m := ""
			for i, newM := range startMetrics {
				if val == newM && val != "" {
					m = newM
					startMetrics[i] = ""
					break
				}
			}
			if m == "" {
				fmt.Printf("ERROR: \"%s\" is not valid metric\n", val)
				continue
			}
			res = append(res, m)
		}
	}
	if len(res) == 0 {
		res = startMetrics
	}
	return res
}

func input() []int {
	var res []int
	reader := bufio.NewReader(os.Stdin)
	for true {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		str = strings.TrimSpace(str)
		val, err := strconv.Atoi(str)
		if err != nil || val < -100000 || val > 100000 {
			fmt.Printf("ERROR: \"%s\" is not valid value\n", str)
			continue
		}
		res = append(res, val)
	}
	return res
}

func fMean(numbers []int) string {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return fmt.Sprintf("Mean: %.2f", float64(sum)/float64(len(numbers)))
}

func fMedian(numbers []int) string {
	var num float64
	sort.Ints(numbers)
	l := len(numbers)
	if l%2 == 1 {
		num = float64(numbers[l/2])
	} else {
		num = (float64(numbers[l/2-1]) + float64(numbers[l/2])) / 2
	}
	return fmt.Sprintf("Median: %.2f", num)
}

func fMode(numbers []int) string {
	var res int
	maxRep := 0
	sort.Ints(numbers)
	if len(numbers) == 1 {
		res = numbers[0]
	}
	for i, rep := 0, 1; i < len(numbers)-1; i++ {
		if numbers[i] == numbers[i+1] && i+2 != len(numbers) {
			rep++
		} else if numbers[i] == numbers[i+1] && rep+2 > maxRep {
			res = numbers[i]
		} else if rep > maxRep {
			maxRep = rep
			res = numbers[i]
			rep = 1
		} else {
			rep = 1
		}
	}
	return fmt.Sprintf("Mode: %d", res)
}

func fSd(numbers []int) string {
	l := float64(len(numbers))
	var mean float64 = 0.0
	for _, num := range numbers {
		mean += float64(num)
	}
	mean /= l
	var sum float64 = 0.0
	for _, num := range numbers {
		sum += math.Pow(float64(num)-mean, 2)
	}
	res := 0.0
	if len(numbers) > 1 {
		res = math.Sqrt(sum / (l - 1))
	}
	return fmt.Sprintf("SD: %.2f", res)
}

func getAnswer(metrics []string, numbers []int) []string {
	var res []string
	if len(numbers) == 0 {
		fmt.Println("ERROR: Empty input")
		return res
	}
	for _, m := range metrics {
		if m == "Mean" {
			res = append(res, fMean(numbers))
		} else if m == "Median" {
			res = append(res, fMedian(numbers))
		} else if m == "Mode" {
			res = append(res, fMode(numbers))
		} else if m == "SD" {
			res = append(res, fSd(numbers))
		}
	}
	return res
}

func output(ans []string) {
	for _, s := range ans {
		fmt.Println(s)
	}
}

func main() {
	metrics := getMetrics()
	numbers := input()
	answer := getAnswer(metrics, numbers)
	output(answer)
}
