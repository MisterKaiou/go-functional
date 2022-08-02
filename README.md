# go-functional

Languages: [PT](README_pt.md), **EN**

## What is it?

It is a small library in Go that contains the implementation of common monads in functional programming that makes use of [generics](https://go.dev/doc/tutorial/generics) support.

## How to install

```sh
go get github.com/MisterKaiou/go-functional
```

## Why?

The current way of returning a value or an error, with a ` ([value], error) ` tuple, is not very expandable. Since you would have to resort to the classic ` if err != nil { do() } ` to handle the error, which quickly gets tiresome.

Of course, in functional languages this would hardly be a problem as the entire language and its APIs are created with monads in mind. We don't have this on Go's API, but even if you do it only in your code, it's a drastic change that will require, at most, getting used to reading `result.Map(** parameters **)`, which in my opinion , it's not much to ask considering how clean the code can get without the countless lines of `if err != nil`.

## Behavior

### Result
- Every `Result` returned is a pointer (`*Result`);
- The value of *Ok* for an error `Result` is `nil`;
- The value of *Ok* for a successful `Result` is the value it holds;
- The *error* value for an error `Result` is the error that was used to create it;
- The *error* value for a successful `Result` is `nil`.

### Option
- Every `Option` returned is a pointer (`*Option`);
- The value of *Some* for an `Option` representing nothing is `nil`;
- The value of *Some* for an `Option` that represents something is the value it holds;
- The value of *None* for an `Option` representing nothing is a pointer to a newly created instance of `None`
- The value of *None* for an `Option` that represents something is `nil`

You don't have access to the internal values, and this is by design, you really shouldn't do this or the purpose of these types would be in vain. `Unwrap` is a *convenience* and I still discourage its use.

<br/>

## Types

In order not to prolong the description too much, examples of how to handle *error* or *None* cases will be shown only once. Since, the same behavior happens in all types and functions.

<br/>

### **Result[Of any]**
****
**Result[Of any]** is a struct that represents a result that contains *`Of`* or *`error`*. Since we have an interface to represent errors within the language itself, I chose to use it.

<br/>

### Examples:
<br/>

**String:** *`Result`* implements the `stringer` interface, and the value returned when calling the `String` method depends on its internal value.
- If **Ok**: `String` will return `fmt.String` applied to the internal value. Ex:
```go
str := result.Ok(42).String() // Same as fmt.String(42)
print(str) // Prints: "42"
```
- If it is **Error**: String will return `error.Error()` in the internal value. Ex:
```go
str := result.Error[int](errors.New("error")).String() // Same as err.Error()
print(str) // Prints: error
```
<br/>

**Map:** Given a *`Result`*, it accepts a function that takes its value and returns another value. Finally, the returned value is stored inside a new *`Result`*

**Ok:**
```go
res := result.Ok(42) // *Result[int]
mappedRes := result.Map(res, func(val int) bool { return val == 42 }) // *Result[bool]

printf("The value of mappedRes is: %s", mappedRes)
// Prints: The value of mappedRes is: true"
```
**Erro:**

```go
res := result.Error[int](errors.New("error")) // *Result[int]
mappedRes := result.Map(res, func(val int) bool { return val == 0 }) //*Result[bool]

printf("The value of mappedRes is: %s", mappedRes) 
// Prints: The value of mappedRes is: error"
```
<br/>

**MapError:** Given a *`Result`*, it accepts a function that receives its error and returns another error. Finally, the returned value is stored in a new *`Result`*.

```go
res := Error[int](errors.New("something failed"))

mapped := MapError(res, func(err error) error { return errors.New(fmt.Sprint("oh no ", err)) })
println(mapped.String()) // Prints: oh no something failed
```

<br/>

**Bind:** Given a *`Result`*, it accepts a function that takes its value and returns another *`Result`*.

```go
res := result.Ok(42) // *Result[int]
bRes := result.Bind(res, func(val int) bool { return result.Ok(val == 42) }) // *Result[bool]

printf("The value of bRes is: %s", bRes)
// Prints: The value of bRes is: true"
```
<br/>

**Match:** Given a *`Result`*, it accepts two functions that return a value of the same type, but the first one receives the result contained in it and the second one receives the error.

```go
res := result.Ok(42)

matchedRes := result.Match(res,
	func(ok int) string { return "All fine around here" },
	func(err error) string { return err.Error() })


printf("The value of matchedRes is: %s", matchedRes)
// Prints: The value of matchedRes is: All fine around here"
```
<br/>

**Unwrap:** It is a *`Result`* method. Returns the value contained in it or calls `panic`.

*Note: Prefer to use `Match` instead of `Unwrap`*
```go
res := result.Ok(69).Unwrap() // res = 69
// Ou
res := result.Error[bool](errors.New("error")).Unwrap() // Panic
```
<br/>

### **Option[Of any]**
****
**Option[Of any]** is a struct that represents either *something* or *nothing*.

<br/>

### Examples:
<br/>

**String:** *`Option`* implements the `stringer` interface, and the value when calling the `String` method depends on its internal value.
- If **Some**: `String` will return the value of `fmt.String` applied to the internal value. Ex:
```go
str := option.Some(42).String() // Same as fmt.String(42)
print(str) // Prints: "42"
```
- If **None**: String will return `"None"`. Ex:
```go
str := option.None[int]().String()
print(str) // Prints: None
```
<br/>

**Map:** Given an *`Option`*, it accepts a function that takes its value and returns another value. Finally, the returned value is stored inside a new *`Option`*

**Ok:**
```go
opt := option.Some(42) // *Option[int]
mappedOpt := option.Map(opt, func(val int) bool { return val == 42 }) // *Option[bool]

printf("The value of mappedOpt is: %s", mappedOpt)
// Prints: The value of mappedOpt is: true"
```
**Erro:**

```go
opt := option.None[int]() // *Option[int]
mappedOpt := option.Map(opt, func(val int) bool { return val == 0 }) //*Option[bool]

printf("The value of mappedOpt is: %s", mappedOpt) 
// Prints: The value of mappedOpt is: None"
```
<br/>

**Bind:** Given an *`Option`*, it accepts a function that takes its value and returns another *`Option`*.

```go
opt := option.Some(42) // *Option[int]
bRes := option.Bind(opt, func(val int) bool { return option.Some(val == 42) }) // *Option[bool]

printf("The value of bRes is: %s", bRes)
// Prints: The value of bRes is: true"
```
<br/>

**Match:** Given an *`Option`*, accepts two functions that return a value of the same type. The first receives the result contained in it and the second does not receive parameters.

```go
opt := option.Some(42)

matchedOpt := option.Match(opt,
	func(ok int) string { return "Something here" },
	func() string { return "Nothing here" })

printf("The value of matchedOpt is: %s", matchedOpt)
// Prints: The value of matchedOpt is: Something here"
```
<br/>

**Unwrap:** It is a *`Option`* method. Returns the value contained in it or calls `panic`.

*Note: Prefer to use `Match` instead of `Unwrap`*
```go
res := option.Some(69).Unwrap() // res = 69
// Ou
res := option.None[bool]().Unwrap() // Panic
```

<br/>

### [*More info on the docs!*](https://pkg.go.dev/github.com/MisterKaiou/go-functional)

<br/>

## Real Example

The code snippet below was taken from a real application of mine using this library together with [gin](https://github.com/gin-gonic/gin):
```go
func (h *Handler) Login() response.HttpResult {
	nonce := h.NonceGenerator()
	stateStr := r.Bind(nonce, func(it auth.Nonce) *r.Result[string] {
		saved := h.NonceStore.Save(it)
		return r.Map(saved, func(_ unit.Unit) string {
			return it.String()
		})
	})
	authUrl := r.Bind(stateStr, auth.GenerateAuthenticationURL)
	return response.Redirect(authUrl)
}
```
Apart from the abstractions I created, like `response.HttpResult` and `response.Redirect`, notice that every call to `Bind` or `Map` would be `if err != nil { return/panic }`. A method that just creates the OAuth login URL would be way longer than these 9 lines if we weren't using this code style.

Adding the `r` alias to the `import` helps make reading less tiring.

## Future plans

Add other monads like:
- IO
- Either
- State (The need for this in Go is debatable as we don't have immutability unless we consider *Pass by Value* something close to that)

Use cases for these are kind of rare, so it's not included in the library as it is now. But if it is a requirement for at least someone, it will be implemented.