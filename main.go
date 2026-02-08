package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/digitalstraw/retry-go/example"
	"github.com/digitalstraw/retry-go/re"
)

var attempts = 0

func okFunc() (int, error) {
	if attempts++; attempts < 3 {
		fmt.Println("okFunc attempt", attempts, "failed at", time.Now())
		return 0, errors.New("overridden error")
	}
	fmt.Println("okFunc succeeded on attempt", attempts, ".")
	return 42, nil
}

func failFunc() (int, error) {
	fmt.Println("failFunc failed. Will wait 2s from", time.Now())
	return 0, errors.New("expected error")
}

func main() {
	ctx := context.Background()

	// Try default policy (Const) with ok result
	fmt.Println("--- okFunc() on default ConstPolicy --- ")
	i, err := re.Try(ctx, okFunc)
	if err != nil {
		panic("Unexpected error: " + err.Error())
	}
	fmt.Println("okFunc() returned", i)
	fmt.Println("")

	// Try Const policy with fail
	fmt.Println("--- failFunc() on ConstPolicy --- ")
	attempts = 0
	start := time.Now()
	fmt.Println("failFunc() on Const() policy with 2s interval and 5s timeout...")
	fmt.Println("Start time=", start)
	fmt.Println("Timout time=", start.Add(5*time.Second))

	i, err = re.Try(ctx, failFunc, re.Const().WithInterval(2*time.Second).WithTimeout(5*time.Second))
	if err != nil {
		fmt.Println("End time=", time.Now())
		fmt.Println("failFunc() failed with", err.Error(), "after", time.Since(start))
	}
	fmt.Println("")

	// Try Custom policy with ok result
	fmt.Println("--- okFunc() on CustomPolicy --- ")
	attempts = 0
	start = time.Now()

	fmt.Println("okFunc() on Custom() with increasing interval by 100ms on each attempt and 5s timeout...")
	fmt.Println("Start time=", start)
	fmt.Println("Timout time=", start.Add(5*time.Second))

	customPolicy := example.Custom()

	expectedSleeps := []time.Duration{
		customPolicy.SleepDuration(1, 0),
		customPolicy.SleepDuration(2, 0),
	}

	fmt.Println("Expected sleeps:", expectedSleeps)

	i, err = re.Try(ctx, okFunc, customPolicy)
	if err != nil {
		panic("Unexpected error: " + err.Error())
	}
	fmt.Println("End time=", time.Now())
	fmt.Println(
		"okFunc() succeeded with result",
		i,
		"after",
		time.Since(start),
		"and",
		customPolicy.Failures(),
		"failures.",
	)
}
