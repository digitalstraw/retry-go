# `re.Try` Function
`re.Try` provides a simple and idiomatic way to implement retry logic in Go. It executes a function repeatedly until 
it succeeds or a timeout occurs.

# Import
```bash
go get github.com/digitalstraw/retry-go/re
```

# Usage
The function accepts two parameters:

- `fn`: A function that returns a value of type `T` and an error.  
- `rp`: An optional implementation of the `re.TryPolicy` interface defining the retry behavior.

## Default Usage
If no retry policy is specified, `re.Const()` is used by default, with a 30-second timeout and a 1-second interval between attempts.

```go
db, err := re.Try(ctx, func() (*sql.DB, error) {
    return sql.Open("mysql", dsn)
})
```

## Using the Default Backoff Retry Policy
The `Backoff` policy in its default configuration has a 2-minute timeout with a doubling interval between attempts, starting at 100ms.

```go
fn := func() (*sql.DB, error) {
    return sql.Open("mysql", dsn)
}
db, err := re.Try(ctx, fn, re.Backoff())
```

## Fully Customized Backoff Retry Policy
```go
fn := func() (*sql.DB, error) {
    return sql.Open("mysql", dsn)
}
// Genrates the following sleep durations: 100ms, 300ms, 900ms, 2.7s, 8.1s, 24.3s, 60s, 60s, ...
db, err := re.Try(ctx, fn, re.Backoff().WithFactor(3).WithInterval(100*time.Millisecond).WithMaxInterval(60*time.Second).WithTimeout(5*time.Minute))
```

# Custom Retry Policies
You can define your own retry policies by implementing the `re.TryPolicy` interface.

## Example: Custom Policy
```go
const (
	defaultTimeout = 10 * time.Second
)

type CustomPolicy struct {
	*re.Policy

	failures int
}

var _ re.TryPolicy = (*CustomPolicy)(nil)

func Custom() *CustomPolicy {
	c := &CustomPolicy{}
	c.Policy = &re.Policy{
		Interval: 1 * time.Second,
		StopAt:   time.Now().Add(defaultTimeout),
		Self:     c,
	}
	return c
}

// SleepDuration produces following durations in default: 3s, 5s, 7s, timeout.
func (c *CustomPolicy) SleepDuration(attempt int, _ time.Duration) time.Duration {
	c.failures++
	return c.Interval + time.Duration(attempt)*2*time.Second // Custom logic
}

func (c *CustomPolicy) Failures() int {
	return c.failures
}
```

## Using Custom Policy
```go
fn := func() (*sql.DB, error) {
    return sql.Open("mysql", dsn)
}
db, err := re.Try(ctx, fn, Custom().WithTimeout(1*time.Minute))
```

# How to Run Linter
To run the linter, use the following command in the terminal:
```bash
golangci-lint run --fix --build-tags=tests
```


# Contributing
Contributions are welcome! Please open an issue or submit a pull request on GitHub. 



# License
This code is licensed under the MIT License. See the LICENSE file for details.

# Changelog
- **v1.0.0** 
  - Added context parameter to `Try()` function. 
  - Immediate reaction to context cancellation implemented. 
  - `re.Sleep(ctx, duration)` function added to support context-aware sleeping.
- **v0.0.1** 
  - Initial release with basic retry functionality and default policies.
