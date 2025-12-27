# gotots

Inspired by [Eden](https://elysiajs.com/eden/overview.html) from the creators of **ElysiaJS**, this project aims to provide **end-to-end type safety** between a **Golang backend** and a **TypeScript frontend**.

Define your models once in Go, and this tool automatically generates corresponding TypeScript definitions. These generated types can be consumed by the frontend to ensure compile-time API safety and prevent contract mismatches between backend and frontend.

## Status

- ✅ Convert Go `struct` to TypeScript `interface`  
- 🚧 Convert Go `iota` constants to TypeScript `enum`


## Installation

### CLI

```bash
go install github.com/sairash/gotots/cmd/gotots@latest
```

### Library

```bash
go get github.com/sairash/gotots
```

## Usage

### CLI

```bash
gotots -dir models -output api/types.ts
```

### Library

```go
package main

import (
	"log"

	"github.com/sairash/gotots"
)

func main() {
	err := gotots.New().
		FromDir("models").
		ToFile("api/types.ts").
		Generate()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example

Given this Go code in `models/user.go`:

```go
package models

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Avatar    *string
	CreatedAt time.Time
}

type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  User   `json:"author,omitempty"`
	Tags    struct {
		Tag        string `json:"tag"`
		TaggedUser User   `json:"tagged_user"`
	} `json:"tags"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

```

Running `gotots -dir models -output types.ts` generates:

```typescript
/* Do not change, this code is generated from Golang structs */

export interface User {
	ID: number;
	Name: string;
	Email: string;
	Avatar?: string | null;
	CreatedAt: string;
};

export interface Post {
	id: number;
	title: string;
	content: string;
	author?: User;
	tags: {
		tag: string;
		tagged_user: User;
	};
	created_at: string;
	updated_at: string;
};
```

***For more examples see [Examples.md](./EXAMPLES.md)***

## Type Mappings

| Go Type | TypeScript Type |
|---------|-----------------|
| `string` | `string` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `*T` | `T \| null` (optional) |
| `map[K]V` | `Record<K, V>` |
| `time.Time` | `string` |
| `interface{}`, `any` | `any` |
| `uuid.UUID` | `string` |
| `json.RawMessage` | `any` |
| `sql.NullString` | `string \| null` |
| `sql.NullInt64` | `number \| null` |
| `decimal.Decimal` | `string` |

## Refrenced Projects
- [github.com/tkrajina/typescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs)
- [github.com/StirlingMarketingGroup/go2ts](https://github.com/StirlingMarketingGroup/go2ts)
- [github.com/OneOfOne/struct2ts](https://github.com/OneOfOne/struct2ts)
