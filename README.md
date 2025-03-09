# alog

An opinionated slog handler.

## Features

- JSON logging with customizable attributes
- Configurable output, log level, and source inclusion
- Contextual logging support

## Installation

To install the package, run:

```sh
go get github.com/arthurdotwork/alog
```

## Usage

### Basic Usage

```go
package main

import (
    "github.com/arthurdotwork/alog"
    "log/slog"
)

func main() {
    logger := alog.Logger()
    logger.Info("This is an info message")
}
```

### Set as default

You can set the logger as the default logger:

```go
package main

import (
    "github.com/arthurdotwork/alog"
    "log/slog"
)

func main() {
    slog.SetDefault(alog.Logger())
    slog.Info("This is an info message")
}
```


### Custom Options

You can customize the logger by providing options:

```go
package main

import (
    "bytes"
    "github.com/arthurdotwork/alog"
    "log/slog"
)

func main() {
    var buf bytes.Buffer
    logger := alog.Logger(
        alog.WithOutput(&buf),
        alog.WithLevel(slog.LevelDebug),
        alog.WithSource(false),
        alog.WithAttrs(slog.Attr{Key: "app", Value: slog.StringValue("myapp")}),
    )
    logger.Debug("This is a debug message")
}
```

### Contextual Logging

You can add contextual information to your logs:

```go
package main

import (
    "context"
    "github.com/arthurdotwork/alog"
    "log/slog"
)

func main() {
    logger := alog.Logger()
    ctx := alog.Append(context.Background(), "request_id", "12345")
    logger.InfoContext(ctx, "This is an info message with context")
}
```

## Testing

To run the tests, use:

```sh
go test ./...
```

## Changelog

All notable changes to this project will be documented in the `CHANGELOG.md` file.

## License

This project is licensed under the MIT License.
