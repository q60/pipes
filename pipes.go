package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	left = iota
	right
	down
	up
)
const (
	reset   = "\x1B[0m"
	red     = "\x1B[31m"
	green   = "\x1B[32m"
	yellow  = "\x1B[33m"
	blue    = "\x1B[34m"
	magenta = "\x1B[35m"
	cyan    = "\x1B[36m"
)

type pipesStruct struct {
	Init  []string
	Down  []string
	Up    []string
	Left  []string
	Right []string
}

type dimensions struct {
	X int
	Y int
}

// Prints out to x:y colored with c
func out(x int, y int, c string, out string) {
	fmt.Printf("\x1B[%d;%dH%s\x1B[1m%s\x1B[0m", y, x, c, out)
}

// Calculates maximum possible pipe length regarding terminal size (cols*rows)
func maxPipeLength(res int) (maxVal int) {
	if res <= 100 {
		maxVal = 200
	} else if res >= 5000 {
		maxVal = 80
	} else {
		maxVal = int(600000/(49*res) + 3800/49)
	}
	return
}

// Generates cryptographic-safe int64 in range [0, max)
func randomNumber64(max int64) (randInt64 int64) {
	number, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	randInt64 = number.Int64()
	return
}

// Generates cryptographic-safe int in range [min, max]
func randN(min int, max int) (randInt int) {
	randInt = int(randomNumber64(int64(max-min+1))) + min
	return
}

// Takes one random element from slice/list and returns it
func choice(slice interface{}) (element interface{}) {
	sliceVals := reflect.ValueOf(slice)
	element = sliceVals.Index(int(randomNumber64(int64(sliceVals.Len())))).Interface()
	return
}

func main() {
	var (
		pipes = pipesStruct{
			[]string{"┃", "━"},
			[]string{"┃", "┛", "┗"},
			[]string{"┃", "┓", "┏"},
			[]string{"━", "┏", "┗"},
			[]string{"━", "┓", "┛"},
		}
		colors = [7]string{
			red,
			green,
			yellow,
			blue,
			magenta,
			cyan,
			reset,
		}
		direction      int
		pipe           string
		cords          dimensions
		term           dimensions
		resolution     []string
		entryPoint     string
		pipeColor      string
		startX, startY int
		maxPipes       int
		counter        int
		chance         int
		help           bool
		style          string
	)
	flag.StringVar(&style, "style", "pipe", "Pipes style")
	flag.StringVar(&style, "s", "pipe", "Pipes style")
	flag.BoolVar(&help, "help", false, "Help message")
	flag.BoolVar(&help, "h", false, "Help message")
	flag.Parse()
	if help {
		fmt.Printf("%spipes%s 1.5.1\n"+
			"Llathasa Veleth <llathasa@outlook.com>\nPipe generator.\n\n"+
			"%sUSAGE:%s\n"+
			"\tpipes [FLAGS] [OPTIONS]\n\n"+
			"%sFLAGS:%s\n"+
			"\t%s-h%s, %s--help%s\tPrints help information.\n\n"+
			"%sOPTIONS:%s\n"+
			"\t%s-s%s, %s--style <wire|thin|knob|double|cross>%s\n"+
			"\t\tChooses different pipes style.\n",
			green, reset, yellow, reset,
			yellow, reset, green, reset,
			green, reset, yellow, reset,
			green, reset, green, reset)
		os.Exit(0)
	}
	switch style {
	case "wire":
		pipes = pipesStruct{
			[]string{"│", "─"},
			[]string{"│", "╯", "╰"},
			[]string{"│", "╮", "╭"},
			[]string{"─", "╭", "╰"},
			[]string{"─", "╮", "╯"},
		}
	case "thin":
		pipes = pipesStruct{
			[]string{"│", "─"},
			[]string{"│", "┘", "└"},
			[]string{"│", "┐", "┌"},
			[]string{"─", "┌", "└"},
			[]string{"─", "┐", "┘"},
		}
	case "knob":
		pipes = pipesStruct{
			[]string{"╽", "╼"},
			[]string{"╿", "┙", "┕"},
			[]string{"╽", "┑", "┍"},
			[]string{"╼", "┍", "┕"},
			[]string{"╾", "┑", "┙"},
		}
	case "double":
		pipes = pipesStruct{
			[]string{"║", "═"},
			[]string{"║", "╝", "╚"},
			[]string{"║", "╗", "╔"},
			[]string{"═", "╔", "╚"},
			[]string{"═", "╗", "╝"},
		}
	case "cross":
		pipes = pipesStruct{
			[]string{"|", "-"},
			[]string{"|", "/", "\\"},
			[]string{"|", "\\", "/"},
			[]string{"-", "/", "\\"},
			[]string{"-", "\\", "/"},
		}
	}
	cmd := exec.Command("sh", // gets current cursor position (rows)
		"-c",
		"exec</dev/tty;ol=$(stty -g);stty raw -echo min 0;echo -en '\033[6n'>/dev/tty;IFS=';' read -r -d R -a pos;stty $ol;row=$((${pos[0]:2}-1));col=$((${pos[1]}-1));echo -ne $row")
	cmd.Stdin = os.Stdin
	savePos, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\x1B[%sS"+ // Scrolls down by savePos to save terminal history
		"\x1B[?25l", // Hides cursor
		string(savePos))
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\x1B[?25h"+ // Shows cursor
			"\x1B[1;1H"+ // Move cursor to [1;1]
			"\x1B[0J"+ // ANSI ED - Erase below
			"\x1B[1J"+ // ANSI ED - Erase above
			"%s[EXITED]%s\n",
			green, reset)
		os.Exit(0)
	}()
	for {
		cmd := exec.Command("stty", "size") // Gets current terminal size (cols*rows)
		cmd.Stdin = os.Stdin
		stdout, _ := cmd.Output()
		resolution = strings.Split(strings.Trim(string(stdout), "\n"), " ")
		termX, _ := strconv.ParseInt(resolution[1], 10, 0)
		termY, _ := strconv.ParseInt(resolution[0], 10, 0)
		term = dimensions{
			int(termX),
			int(termY),
		}
		entryPoint = choice(pipes.Init).(string)
		pipeColor = choice(colors).(string)
		startX, startY = randN(1, term.X), randN(1, term.Y)
		out(startX,
			startY,
			pipeColor,
			entryPoint)
		if entryPoint == pipes.Init[0] {
			direction = choice([2]int{down, up}).(int)
		} else {
			direction = choice([2]int{right, left}).(int)
		}
		cords = dimensions{
			startX,
			startY,
		}
		maxPipes = randN(20, maxPipeLength(term.X*term.Y))
		counter = 0
		for {
			chance = randN(0, 24)
			switch direction {
			case down:
				if cords.Y >= term.Y {
					cords.Y = 0
					counter++
					pipeColor = choice(colors).(string)
				}
				switch chance {
				case 0:
					pipe = pipes.Down[1]
				case 1:
					pipe = pipes.Down[2]
				default:
					pipe = pipes.Down[0]
				}
				cords.Y++
				out(cords.X,
					cords.Y,
					pipeColor,
					pipe)
				if pipe == pipes.Down[1] {
					direction = left
				} else if pipe == pipes.Down[2] {
					direction = right
				}
			case up:
				if cords.Y <= 1 {
					cords.Y = term.Y + 1
					counter++
					pipeColor = choice(colors).(string)
				}
				switch chance {
				case 0:
					pipe = pipes.Up[1]
				case 1:
					pipe = pipes.Up[2]
				default:
					pipe = pipes.Up[0]
				}
				cords.Y--
				out(cords.X,
					cords.Y,
					pipeColor,
					pipe)
				if pipe == pipes.Up[1] {
					direction = left
				} else if pipe == pipes.Up[2] {
					direction = right
				}
			case left:
				if cords.X <= 1 {
					cords.X = term.X + 1
					counter++
					pipeColor = choice(colors).(string)
				}
				switch chance {
				case 0:
					pipe = pipes.Left[1]
				case 1:
					pipe = pipes.Left[2]
				default:
					pipe = pipes.Left[0]
				}
				cords.X--
				out(cords.X,
					cords.Y,
					pipeColor,
					pipe)
				if pipe == pipes.Left[1] {
					direction = down
				} else if pipe == pipes.Left[2] {
					direction = up
				}
			case right:
				if cords.X >= term.X+1 {
					cords.X = 1
					counter++
					pipeColor = choice(colors).(string)
				}
				switch chance {
				case 0:
					pipe = pipes.Right[1]
				case 1:
					pipe = pipes.Right[2]
				default:
					pipe = pipes.Right[0]
				}
				cords.X++
				out(cords.X,
					cords.Y,
					pipeColor,
					pipe)
				if pipe == pipes.Right[1] {
					direction = down
				} else if pipe == pipes.Right[2] {
					direction = up
				}
			}
			time.Sleep(10 * time.Millisecond)
			if counter == maxPipes {
				fmt.Print("\x1B[0J" + // ANSI ED - Erase below
					"\x1B[1J" + // ANSI ED - Erase above
					"\x1B[?25l") // Hides cursor
				break
			}
		}
	}
}
