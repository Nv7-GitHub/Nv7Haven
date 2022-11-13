# Poll Data Types
## Combo Poll
```go
struct {
  els []float64
  result float64 | string // String if new, float64 if exists
}
```

## Comment Poll
```go
struct {
  elem float64
  new string
  old string
}
```

## Cat Comment Poll
```go
struct {
  cat string
  new string
  old string
}
```

## Image Poll
```go
struct {
  elem float64
  new string
  old string
}
```

## Cat Image Poll
```go
struct {
  cat string
  new string
  old string
}
```

## Color Poll
```go
struct {
  elem float64
  new float64
  old float64
}
```

## Cat Color Poll
```go
struct {
  cat string
  new float64
  old float64
}
```

## Categorize Poll
```go
struct {
  cat string
  elems []float64
}
```

## Un-Categorize Poll
```go
struct {
  cat string
  elems []float64
}
```