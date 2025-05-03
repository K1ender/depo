# Depo

**Depo** is a tiny dependency injection helper for Go. It allows you to register struct pointers and inject them into functions automatically based on parameter types.

## Features

* ⚙️ Lightweight, no external dependencies
* 🔍 Uses Go’s reflection to resolve dependencies by type name
* 🧪 Suitable for small projects, tools, or test setups

## Installation

```bash
go get github.com/k1ender/depo
```

## Usage

### Registering and using a dependency

```go
package main

import (
	"fmt"

	"github.com/k1ender/depo"
)

type Logger struct {
	Prefix string
}

func NewLogger() *Logger {
	return &Logger{Prefix: "[Depo] "}
}

func main() {
	logger := NewLogger()

	container := depo.New(logger)

	container.Use(func(l *Logger) {
		fmt.Println(l.Prefix + "Dependency injection works!")
	})
}
```

### Sharing mutable dependencies

```go
package main

import (
	"fmt"

	"github.com/k1ender/depo"
)

type Store struct {
	Data map[string]string
}

func NewStore() *Store {
	return &Store{Data: map[string]string{}}
}

func main() {
	store := NewStore()
	container := depo.New(store)

	container.Use(func(s *Store) {
		s.Data["hello"] = "world"
	})

	container.Use(func(s *Store) {
		fmt.Println("Value:", s.Data["hello"]) // Output: Value: world
	})
}
```

## API

### `depo.New(deps ...any) *DependencyPool`

Registers any number of pointer dependencies. Only the concrete type name is used, so multiple types with the same name will conflict.

### `(*DependencyPool).Use(fun any) error`

Injects registered dependencies into the given function by matching parameter types by name.

## Limitations

* 💡 Dependencies must be **pointers to structs**
* 🧱 No interface injection — only concrete types
* ⚠️ Types with the same name will overwrite each other

## License

MIT