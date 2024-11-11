package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"

	"github.com/danielmesquitta/stay-active/internal/app/cli"
)

func main() {
	startedAt := time.Now()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	interval, timeout, err := cli.Run()
	if err != nil {
		panic(err)
	}

	if err := validate(timeout, interval); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Program will run for %s pressing one random letter every %s\n",
		timeout,
		interval,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				quit <- syscall.SIGINT
				return
			default:
				time.Sleep(interval)
				randomLetter := string(rune(rand.Intn(26) + 'a'))
				robotgo.TypeStr(randomLetter)
			}
		}
	}()

	<-quit
	fmt.Println()
	fmt.Println("Program was running for", time.Since(startedAt))
}

func validate(timeout, interval time.Duration) error {
	if interval == 0 && timeout == 0 {
		return errors.New("No interval or timeout was set")
	}
	if timeout == 0 {
		return errors.New("Timeout must be greater than 0")
	}
	if interval == 0 {
		return errors.New("Interval must be greater than 0")
	}
	if interval > timeout {
		return errors.New("Interval must be less than timeout")
	}
	return nil
}
