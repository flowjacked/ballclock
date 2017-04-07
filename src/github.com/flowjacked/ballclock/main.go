package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/flowjacked/ballclock/queue"
	"github.com/flowjacked/ballclock/stack"
)

var (
	// CLI flags
	ballCount = flag.Int("ballCount", 0, "Number of balls in the ball queue. Valid values are between 27 and 127")
	runTime   = flag.Int("runTime", 0, "The number of minutes to run the clock")

	// For ball control
	clock = make(chan int)
	quit  = make(chan int)

	// Set stacks
	oneMinStack  = stack.NewStack(4)  // Holds 4 balls (ints)
	fiveMinStack = stack.NewStack(11) // Holds 11 balls (ints)
	oneHourStack = stack.NewStack(11) // Holds 11 balls (ints)
	ballQueue    *queue.Queue

	// Set a counter for the total number of minutes
	timeInMinutes = 0
)

const (
	MinBalls = 27
	MaxBalls = 127
)

func main() {
	// Parse flags
	flag.Usage = func() {
		version := "0.001"
		app := "ballclock"
		fmt.Fprintf(os.Stderr, "Usage of %s v%s:\n", app, version)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "type help for help\n")
	}
	flag.Parse()
	if flag.Arg(0) == "help" {
		flag.Usage()
		return
	}

	// Check for the proper range
	if *ballCount < MinBalls || *ballCount > MaxBalls {
		fmt.Printf("ERROR: ballCount must be between %d and %d\n", MinBalls, MaxBalls)
		flag.Usage()
		return
	}

	// Set channels
	oneMinTrack := make(chan int)
	fiveMinTrack := make(chan int)
	oneHourTrack := make(chan int)

	// Create and populate ball clock queue
	qLength := *ballCount // Should be pulled from command line
	ballQueue = queue.NewQueue(qLength)
	for i := 1; i <= qLength; i++ {
		_ = ballQueue.Push(i)
	}
	// Save the queue's state for later comparison
	ballQueue.SaveState()

	// Call our track handlers
	go clockWatcher(oneMinTrack, *runTime)
	go oneMinWatcher(oneMinTrack, fiveMinTrack)
	go fiveMinWatcher(fiveMinTrack, oneHourTrack)
	go oneHourWatcher(oneHourTrack)

	// Go routines need to tell us when to stop waiting
	<-quit

	// Printed info
	runIndefinitely := 0
	if *runTime == runIndefinitely {
		fmt.Printf("%d balls cycle after %d days\n", *ballCount, timeInMinutes/60/24)
	} else {
		state := map[string]interface{}{}
		state["min"] = oneMinStack.GetStack()
		state["fivemin"] = fiveMinStack.GetStack()
		state["hour"] = oneHourStack.GetStack()
		state["main"] = ballQueue.GetQueue()
		s, err := json.Marshal(state)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(s))
	}
}

// responsible for taking a ball from the queue and pushing it on the minute channel every minute
func clockWatcher(oneMin chan int, runTime int) {
	// Infinitely pop balls off. If there is an error, then the queue is temporarily empty, just keep going
	for {
		if v, err := ballQueue.Pop(); err == nil {
			oneMin <- v

			// Wait to put more balls on until watchers are ready
			<-clock
		}

		// Only run for a specified amount of time. If it's 0, run indefinitely
		if runTime > 0 && timeInMinutes >= runTime {
			quit <- 0
			return
		}
	}
}

/**
 * responsible for listening on the oneMinTrack chan and pushing each value onto the oneMinStack.
 * if an error is encountered, then pop all values and push them onto the queue.
 **/
func oneMinWatcher(oneMin chan int, fiveMin chan int) {
	// Infinitely watch oneMin channel and do stuff to stuff
	for {
		v := <-oneMin
		timeInMinutes++

		// Keep pushing values onto the stack until it's full, at which point, push all values to queue
		if err := oneMinStack.Push(v); err != nil {
			for {
				// If we're empty, then discontinue popping
				if used, err := oneMinStack.Pop(); err != nil {
					break
				} else {
					// Not empty? Keep pushing values onto queue in reverse order
					_ = ballQueue.Push(used)
				}
			}
			fiveMin <- v
		} else {
			clock <- 0
		}
	}
}

/**
 * responsible for listening on the fiveMin chan and pushing each value onto the fiveMinStack.
 * if an error is encountered, then pop all values and push them onto the queue.
 **/
func fiveMinWatcher(fiveMin chan int, oneHour chan int) {
	// Infinitely watch oneMin channel and do stuff to stuff
	for {
		v := <-fiveMin

		// Keep pushing values onto the stack until it's full, at which point, push all values to queue
		if err := fiveMinStack.Push(v); err != nil {
			for {
				// If we're empty, then discontinue popping
				if used, err := fiveMinStack.Pop(); err != nil {
					break
				} else {
					// Not empty? Keep pushing values onto queue in reverse order
					_ = ballQueue.Push(used)
				}
			}
			// 12th 5 minute ball encountered, push onto oneHour channel
			oneHour <- v
		} else {
			clock <- 0
		}
	}
}

/**
 * responsible for listening on the oneHour chan. If an error is encountered, then
 * pop the values off the stack onto the queue. The 12th ball is then pushed onto
 * the queue
 **/
func oneHourWatcher(oneHour chan int) {
	// Infinitely watch oneHour channel and do stuff to stuff
	for {
		v := <-oneHour

		// Keep pushing values onto the stack until it's full, at which point, push all values to queue
		if err := oneHourStack.Push(v); err != nil {
			for {
				// If we're empty, then discontinue popping
				if used, err := oneHourStack.Pop(); err != nil {
					break
				} else {
					// Not empty? Keep pushing values onto queue in reverse order
					_ = ballQueue.Push(used)
				}
			}
			// 12th 1 hour ball encountered and all other balls have returned to the queue. Now push the 12th ball onto the queue
			_ = ballQueue.Push(v)
			if ballQueue.EqualsOrigin() {
				quit <- 0
				return
			}
		}
		clock <- 0
	}
}
