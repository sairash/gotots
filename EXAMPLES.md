# Examples
This file shows what are the things that are supported

## Simple Struct

**Go**

```go
type User struct {
    ID   int
    Name string
}
```

**TypeScript**

```ts
export interface User {
    ID: number;
    Name: string;
}
```

---

## Pointer Fields (Optional + Nullable)

**Go**

```go
type Profile struct {
    Bio     *string
    Website *string
}
```

**TypeScript**

```ts
export interface Profile {
    Bio?: string | null;
    Website?: string | null;
}
```

---

## JSON Tags

**Go**

```go
type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"lastName"`
    Email     string
}
```

**TypeScript**

```ts
export interface User {
    id: number;
    first_name: string;
    lastName: string;
    Email: string;
}
```

---

## `omitempty` JSON Tag

**Go**

```go
type Profile struct {
    Bio      string `json:"bio,omitempty"`
    Website  string `json:"website,omitempty"`
    Location string `json:"location"`
}
```

**TypeScript**

```ts
export interface Profile {
    bio?: string;
    website?: string;
    location: string;
}
```

---

## Slices → Arrays

**Go**

```go
type Post struct {
    Tags     []string
    Comments []int
}
```

**TypeScript**

```ts
export interface Post {
    Tags: string[];
    Comments: number[];
}
```

---

## Maps → Record

**Go**

```go
type Config struct {
    Settings map[string]string
    Counts   map[string]int
}
```

**TypeScript**

```ts
export interface Config {
    Settings: Record<string, string>;
    Counts: Record<string, number>;
}
```

---

## Nested Structs

**Go**

```go
type Address struct {
    Street string
    City   string
}

type Person struct {
    Name    string
    Address Address
}
```

**TypeScript**

```ts
export interface Address {
    Street: string;
    City: string;
}

export interface Person {
    Name: string;
    Address: Address;
}
```

---

## Anonymous / Inline Structs

**Go**

```go
type Response struct {
    Data struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    } `json:"data"`
}
```

**TypeScript**

```ts
export interface Response {
    data: {
        id: number;
        name: string;
    };
}
```

---

## Multiple Files & Subdirectories

**Go**

``` go
// user.go
type User struct {
    Name string
}

// sub/product.go
type Product struct {
    Price float64
}
```

**TypeScript**

``` ts
export interface User {
    Name: string;
}

export interface Product {
    Price: number;
}
```

---

## Ignoring Test Files

``` bash
user.go        ✅ included
user_test.go   ❌ ignored
```