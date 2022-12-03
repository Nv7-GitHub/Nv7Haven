# Poll Data Types
## Combo Poll
```go
struct {
  els []any (float64)
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
  elems []any (float64)
}
```

## Un-Categorize Poll
```go
struct {
  cat string
  elems []any (float64)
}
```

## Query Poll
```go
struct {
  query string
  edit bool
  kind string // see below
  data any // jsonb
}
```

### Query types
**Element**
```go
struct {
  elem float64
}
```
**Products**
```go
struct {
  query string
}
```
**Parents**
```go
struct {
  query string
}
```
**Inventory**
```go
struct {
  user string
}
```
**Elements**
```go
struct {}
```
**Regex**
```go
struct {
  query string
}
```
**Comparison**
```go
struct {
  query string
  field string // Valid fields: name, image, color, comment, creator, commenter, colorer, imager, treesize
  typ string `enum:"notequal,equal,greater,less"`
  value any // float64 or string
}
```