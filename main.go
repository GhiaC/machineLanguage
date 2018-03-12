package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
	"log"
)

var transactions map[string]map[int]map[int]string
var nodes map[int]string
var start string
var end []string
var way []int

func main() {
	directoryNFA := "NFA/"
	directoryDFA := "DFA/"
	numbers := [5]string{"first", "second", "third", "fourth", "fifth"}
	for _, number := range numbers {
		fmt.Println("Result of " + number + " (NFA)")
		initial()
		readFile(directoryNFA + strings.Title(number) + "_NFA.txt")
		readTestFile(directoryNFA + "Strings_for_" + number + "_NFA.txt")
	}
	for _, number := range numbers {
		fmt.Println("Result of " + number + " (DFA)")
		initial()
		readFile(directoryDFA + strings.Title(number) + "_DFA.txt")
		readTestFile(directoryDFA + "Strings_for_" + number + "_DFA.txt")
	}

}
func initial() {
	nodes = make(map[int]string)
	transactions = make(map[string]map[int]map[int]string)
	way = make([]int, 0)
	end = make([]string, 0)
}
func parseTestCase(line string) {
	way = make([]int, 0)
	r := bufio.NewReader(strings.NewReader(line))
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			way = append(way, int(c)-48)
		}
	}
}

func trace(index int, current string) bool {
	result := false
	if index >= len(way) {
		for _, value := range end {
			if current == value {
				return true
			}
		}
		return false
	}
	tempSlice := transactions[current][way[index]]
	if len(tempSlice) == 0 {
		return false
	} else {
		for _, value := range tempSlice {
			result = result || trace(index+1, value)
		}
		return result
	}
}
func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	lineCounter, nodeCounter := 0, 0
	flag, flagI := false, false
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if lineCounter == 0 {
			readNode(string(line))
		} else if lineCounter == 1 {
			flag = true
		} else if len(nodes) > nodeCounter && flag {
			setTransaction(nodeCounter, string(line))
			nodeCounter++
		} else if flag && !flagI {
			setStartEnd(string(line), "start")
			flagI = true
		} else if flagI && flag {
			setStartEnd(string(line), "end")
		}
		lineCounter++
	}
}
func readTestFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else {
			parseTestCase(string(line))
			fmt.Println(trace(0, start))
		}
	}
}
func setStartEnd(line, state string) {
	r := bufio.NewReader(strings.NewReader(line))
	flag := false
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			if flag {
				if state == "start" {
					start = string(c)
				} else {
					if flag && strings.Compare(string(c), ",") != 0 {
						end = append(end, string(c))
					}
				}
			}
			if string(c) == ":" {
				flag = true
			}
		}
	}
}
func setTransaction(numberOfNode int, line string) {
	r := bufio.NewReader(strings.NewReader(line))
	flagOpen := false
	var tempSlice map[int]string
	indexCounter := 0
	tempMap := make(map[int]map[int]string)
	counterOfMapInMap := 0 // E
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			if flagOpen && string(c) != "," && string(c) != "}" {
				tempSlice[indexCounter] = string(c)
				indexCounter++
			}
			if string(c) == "{" {
				flagOpen = true
				tempSlice = make(map[int]string)
			} else if string(c) == "}" {
				flagOpen = false
				tempMap[counterOfMapInMap] = tempSlice
				counterOfMapInMap++
				indexCounter = 0
			}
		}
	}
	transactions[nodes[numberOfNode]] = tempMap
}
func readNode(line string) {
	r := bufio.NewReader(strings.NewReader(line))
	flag := false
	counter := 0
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			if flag && strings.Compare(string(c), ",") != 0 {
				nodes[counter] = string(c)
				counter++
			}
			if string(c) == ":" {
				flag = true
			}
		}
	}
}
