package re

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RetryTestSuite struct {
	suite.Suite
}

func (s *RetryTestSuite) SetupSuite() {
}

func (s *RetryTestSuite) TearDownSuite() {
}

func TestRetryTestSuite(t *testing.T) {
	suite.Run(t, new(RetryTestSuite))
}

func generate(count int, p TryPolicy) []time.Duration {
	var sleeps []time.Duration
	sleep := p.SleepDuration(1, 0)
	sleeps = append(sleeps, sleep)
	previous := sleep

	for i := 2; i <= count; i++ {
		current := p.SleepDuration(i, previous)
		sleeps = append(sleeps, current)
		previous = current
	}
	return sleeps
}

func (s *RetryTestSuite) TestBackoffPolicy() {
	p := Backoff()
	sleeps := generate(12, p)

	s.Require().Equal(100*time.Millisecond, sleeps[0])
	s.Require().Equal(200*time.Millisecond, sleeps[1])
	s.Require().Equal(400*time.Millisecond, sleeps[2])
	s.Require().Equal(800*time.Millisecond, sleeps[3])
	s.Require().Equal(1_600*time.Millisecond, sleeps[4])
	s.Require().Equal(3_200*time.Millisecond, sleeps[5])
	s.Require().Equal(6_400*time.Millisecond, sleeps[6])
	s.Require().Equal(12_800*time.Millisecond, sleeps[7])
	s.Require().Equal(25_600*time.Millisecond, sleeps[8])
	s.Require().Equal(51_200*time.Millisecond, sleeps[9])
	s.Require().Equal(60_000*time.Millisecond, sleeps[10])
	s.Require().Equal(60_000*time.Millisecond, sleeps[11])
}

func (s *RetryTestSuite) TestBackoffPolicyWithCustomParameters() {
	p := Backoff().WithFactor(3).WithInterval(50 * time.Millisecond).WithMaxInterval(10 * time.Second)
	sleeps := generate(8, p)

	s.Require().Equal(50*time.Millisecond, sleeps[0])
	s.Require().Equal(150*time.Millisecond, sleeps[1])
	s.Require().Equal(450*time.Millisecond, sleeps[2])
	s.Require().Equal(1_350*time.Millisecond, sleeps[3])
	s.Require().Equal(4_050*time.Millisecond, sleeps[4])
	s.Require().Equal(10_000*time.Millisecond, sleeps[5])
	s.Require().Equal(10_000*time.Millisecond, sleeps[6])
	s.Require().Equal(10_000*time.Millisecond, sleeps[7])
}

func (s *RetryTestSuite) TestConstantPolicy() {
	p := Const()
	sleeps := generate(5, p)

	s.Require().Equal(1*time.Second, sleeps[0])
	s.Require().Equal(1*time.Second, sleeps[1])
	s.Require().Equal(1*time.Second, sleeps[2])
	s.Require().Equal(1*time.Second, sleeps[3])
	s.Require().Equal(1*time.Second, sleeps[4])
}

func someMethod(x int) (int, error) {
	if x <= 3 {
		return 0, errors.New("some error")
	}
	return 42, nil
}

func (s *RetryTestSuite) TestConstantRetry() {
	// ARRANGE
	x := 0
	fn := func() (int, error) {
		x++
		return someMethod(x)
	}
	startTime := time.Now()

	// ACT
	result, err := Try(context.Background(), fn, Const().WithInterval(50*time.Millisecond).WithTimeout(160*time.Millisecond))

	// ASSERT
	s.Require().NoError(err)
	s.Require().Equal(42, result)
	s.Require().GreaterOrEqual(time.Since(startTime), 150*time.Millisecond)
	s.Require().LessOrEqual(time.Since(startTime), 160*time.Millisecond)
}

func (s *RetryTestSuite) TestBackoffRetry() {
	// ARRANGE
	x := 0
	fn := func() (int, error) {
		x++
		return someMethod(x)
	}
	startTime := time.Now()

	// ACT
	result, err := Try(context.Background(), fn, Backoff().WithInterval(10*time.Millisecond).WithTimeout(100*time.Millisecond))

	// ASSERT
	s.Require().NoError(err)
	s.Require().Equal(42, result)
	s.Require().GreaterOrEqual(time.Since(startTime), (10+20+40)*time.Millisecond)
	s.Require().LessOrEqual(time.Since(startTime), (10+20+40+10)*time.Millisecond)
}

func (s *RetryTestSuite) TestConstantFailOnTimeout() {
	// ARRANGE
	x := 0
	fn := func() (int, error) {
		return someMethod(x)
	}
	startTime := time.Now()

	// ACT
	result, err := Try(context.Background(), fn, Const().WithInterval(10*time.Millisecond).WithTimeout(20*time.Millisecond))

	// ASSERT
	s.Require().Error(err)
	s.Require().Equal(0, result)
	s.Require().GreaterOrEqual(time.Since(startTime), 20*time.Millisecond)
	s.Require().LessOrEqual(time.Since(startTime), 30*time.Millisecond)
}

func (s *RetryTestSuite) TestBackoffFailOnTimeout() {
	// ARRANGE
	x := 0
	fn := func() (int, error) {
		return someMethod(x)
	}
	startTime := time.Now()

	// ACT
	result, err := Try(context.Background(), fn, Backoff().WithInterval(10*time.Millisecond).WithTimeout(60*time.Millisecond))

	// ASSERT
	s.Require().Error(err)
	s.Require().Equal(0, result)
	s.Require().GreaterOrEqual(time.Since(startTime), 60*time.Millisecond)
	s.Require().LessOrEqual(time.Since(startTime), 80*time.Millisecond)
}

func (s *RetryTestSuite) TestRetryContextCanceled() {
	// ARRANGE
	ctx, cancel := context.WithCancel(context.Background())
	fn := func() (int, error) {
		return someMethod(0)
	}
	startTime := time.Now()

	// ACT
	go func() {
		time.Sleep(11 * time.Millisecond) // ensure that the first attempt is made before canceling the context
		cancel()
	}()
	result, err := Try(ctx, fn, Const().WithInterval(10*time.Millisecond).WithTimeout(60*time.Millisecond))

	// ASSERT
	s.Require().Error(err)
	s.Require().Equal(0, result)
	s.Require().LessOrEqual(time.Since(startTime), 20*time.Millisecond)
}

func (s *RetryTestSuite) TestNoRetryPolicy() {
	// ARRANGE
	ctx, cancel := context.WithCancel(context.Background()) // To be able to end immediately after the first attempt
	fn := func() (int, error) {
		return someMethod(0)
	}
	startTime := time.Now()

	// ACT
	cancel()
	result, err := Try(ctx, fn)

	// ASSERT
	s.Require().Error(err)
	s.Require().Equal(0, result)
	s.Require().LessOrEqual(time.Since(startTime), 20*time.Millisecond)
}

func (s *RetryTestSuite) TestSleep() {
	// ARRANGE
	ctx, cancel := context.WithCancel(context.Background())
	startTime := time.Now()

	// ACT
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := Sleep(ctx, 20*time.Millisecond)

	// ASSERT
	s.Require().Error(err)
	s.Require().Equal(context.Canceled, err)
	s.Require().GreaterOrEqual(time.Since(startTime), 10*time.Millisecond)
	s.Require().LessOrEqual(time.Since(startTime), 30*time.Millisecond)
}
