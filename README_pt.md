# go-functional

Languages: **PT**, [EN](README.md)

## O que é?

É uma pequena biblioteca em Go que contém a implementação de mônadas (monads) comuns em programação funcional que faz uso do suporte a [generics](https://go.dev/doc/tutorial/generics).

## Como instalar
```sh
go get github.com/MisterKaiou/go-functional
```

## Por quê?

A maneira atual de se retornar um valor ou um erro, com uma tupla ` ([valor], error) `, não é muito expansível. Visto que você teria de recorrer ao clássico ` if err != nil { do() } ` para tratar o erro, o que fica rapidamente cansativo.

É claro que em linguagens funcionais isso dificilmente seria um problema já que toda a linguagem e suas APIs são criadas com mônadas em mente. Não temos isso nas APIs em Go, mas mesmo que faça somente no seu código, é uma drástica mudança que vai requerer, no máximo, se acostumar a ler `result.Map(** parâmetros **)`, o que na minha opinião, não é muito a se pedir levando em consideração quão enxuto o código fica sem as inúmeras linhas de `if err != nil`.

## Comportamento

### Result
- Todo `Result` retornado é um ponteiro (`*Result`);
- O valor de *Ok* para um `Result` de erro é `nil`;
- O valor de *Ok* para um `Result` de sucesso é o valor que ele guarda;
- O valor de *erro* para um `Result` de erro é o erro que foi usado para criá-lo;
- O valor de *erro* para um `Result` de sucesso é `nil`.

### Option
- Todo `Option` retornado é um ponteiro (`*Option`);
- O valor de *Some* para um `Option` que representa nada é `nil`;
- O valor de *Some* para um `Option` que representa algo é o valor que ele guarda;
- O valor de *None* para um `Option` que representa nada é um ponteiro para um instância recêm criada de `None`
- O valor de *None* para um `Option` que representa algo é `nil`

Voce não tem acesso aos valores internos, e isso é por design, você realmente não deveria fazer isso ou o propósito destes tipos seria em vão. `Unwrap` é uma *conveniência* e ainda assim desencorajo o uso dela.

<br/>

## Tipos
<br/>

Visando não prolongar demais a descrição, os exemplos de como manusear casos de *erro* ou *None* serão mostrados apenas uma vez. Já que, o mesmo comportamento acontece em todos os tipos e funções.

### **Result[Of any]**
****
**Result[Of any]** é uma struct que representa um resultado que contém *`Of`* ou *`error`*. Já que temos uma interface para representar erros dentro da própria linguagem, optei por fazer uso dela.

<br/>

### Exemplos:
<br/>

**String:** *`Result`* implementa a interface `stringer`, e o valor retornado ao chamar o método `String` depende do seu valor interno.
- Caso seja **Ok**: `String` vai retornar `fmt.String` aplicado ao valor interno. Ex:
```go
str := result.Ok(42).String() // É o mesmo que fmt.String(42)
print(str) // Imprime: "42"
```
- Caso seja **Error**: String vai retornar `error.Error()` no valor interno. Ex:
```go
str := result.Error[int](errors.New("erro")).String() // É o mesmo que err.Error()
print(str) // Imprime: erro
```
<br/>

**Map:** Dado um *`Result`*, ela aceita uma função que recebe o valor deste e retorna outro valor. Por fim, o valor retornado é guardado dentro de um novo *`Result`*

**Ok:**
```go
res := result.Ok(42) // *Result[int]
mappedRes := result.Map(res, func(val int) bool { return val == 42 }) // *Result[bool]

printf("O valor de mappedRes é: %s", mappedRes)
// Imprime na tela: O valor de mappedRes é: true"
```
**Erro:**

```go
res := result.Error[int](errors.New("erro")) // *Result[int]
mappedRes := result.Map(res, func(val int) bool { return val == 0 }) //*Result[bool]

printf("O valor de mappedRes é: %s", mappedRes) 
// Imprime na tela: O valor de mappedRes é: erro"
```
<br/>

**MapError:** Dado um *`Result`*, ela aceita uma função que recebe o erro deste e retorna outro erro. Por fim, o valor retornado é guardado em um novo *`Result`*.

```go
res := Error[int](errors.New("something failed"))

mapped := MapError(res, func(err error) error { return errors.New(fmt.Sprint("oh no ", err)) })
println(mapped.String()) // Imprime: oh no something failed
```

<br/>

**Bind:** Dado um *`Result`*, ela aceita uma função que recebe o valor deste e retorna outro *`Result`*.

```go
res := result.Ok(42) // *Result[int]
bRes := result.Bind(res, func(val int) bool { return result.Ok(val == 42) }) // *Result[bool]

printf("O valor de bRes é: %s", bRes)
// Imprime na tela: O valor de bRes é: true"
```
<br/>

**Match:** Dado um *`Result`*, ela aceita duas funções que retornam um valor do mesmo tipo, mas a primeira recebe o resultado contido nele e a segunda recebe o erro.

```go
res := result.Ok(42)

matchedRes := result.Match(res,
	func(ok int) string { return "Tudo okay por aqui" },
	func(err error) string { return err.Error() })


printf("O valor de matchedRes é: %s", matchedRes)
// Imprime na tela: O valor de matchedRes é: Tudo okay por aqui"
```
<br/>

**Unwrap:** É um método de *`Result`*. Retorna o valor contido nele ou chama `panic`.

*Obs: Prefira usar `Match` ao invés de `Unwrap`*
```go
res := result.Ok(69).Unwrap() // res = 69
// Ou
res := result.Error[bool](errors.New("erro")).Unwrap() // Panic
```
<br/>

### **Option[Of any]**
****
**Option[Of any]** é uma struct que representa ou *algo* ou *nada*.

<br/>

### Exemplos:
<br/>

**String:** *`Option`* implementa a interface `stringer`, e o valor ao chamar o método `String` depende do seu valor interno.
- Caso seja **Some**: `String` vai retornar o valor de `fmt.String` aplicado ao valor interno. Ex:
```go
str := option.Some(42).String() // É o mesmo que fmt.String(42)
print(str) // Imprime: "42"
```
- Caso seja **None**: String vai retornar `"None"`. Ex:
```go
str := option.None[int]().String()
print(str) // Imprime: None
```
<br/>

**Map:** Dado um *`Option`*, ela aceita uma função que recebe o valor deste e retorna outro valor. Por fim, o valor retornado é guardado dentro de um novo *`Option`*

**Ok:**
```go
opt := option.Some(42) // *Option[int]
mappedOpt := option.Map(opt, func(val int) bool { return val == 42 }) // *Option[bool]

printf("O valor de mappedOpt é: %s", mappedOpt)
// Imprime na tela: O valor de mappedOpt é: true"
```
**Erro:**

```go
opt := option.None[int]() // *Option[int]
mappedOpt := option.Map(opt, func(val int) bool { return val == 0 }) //*Option[bool]

printf("O valor de mappedOpt é: %s", mappedOpt) 
// Imprime na tela: O valor de mappedOpt é: None"
```
<br/>

**Bind:** Dado um *`Option`*, ela aceita uma função que recebe o valor deste e retorna outro *`Option`*.

```go
opt := option.Some(42) // *Option[int]
bRes := option.Bind(opt, func(val int) bool { return option.Some(val == 42) }) // *Option[bool]

printf("O valor de bRes é: %s", bRes)
// Imprime na tela: O valor de bRes é: true"
```
<br/>

**Match:** Dado um *`Option`*, aceita duas funções que retornam um valor do mesmo tipo. A primeira recebe o resultado contido nele e a segunda não recebe parâmetros.

```go
opt := option.Some(42)

matchedOpt := option.Match(opt,
	func(ok int) string { return "Algo aqui" },
	func() string { return "Nada aqui" })

printf("O valor de matchedOpt é: %s", matchedOpt)
// Imprime na tela: O valor de matchedOpt é: Algo aqui"
```
<br/>

**Unwrap:** É um método de *`Option`*. Retorna o valor contido nele ou chama `panic`.

*Obs: Prefira usar `Match` ao invés de `Unwrap`*
```go
res := option.Some(69).Unwrap() // res = 69
// Ou
res := option.None[bool]().Unwrap() // Panic
```

<br/>

### [*Mais informações na documentação!*](https://pkg.go.dev/github.com/MisterKaiou/go-functional)

<br/>

## Exemplo Real

O trecho de código abaixo foi retirado de uma aplicação real minha utilizando esta biblioteca em conjunto com  [gin](https://github.com/gin-gonic/gin):
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
Fora as abstrações que eu criei, como `response.HttpResult` e `response.Redirect`, perceba que cada chamada a `Bind` ou `Map` seria `if err != nil { return/panic }`. Um método que somente cria a URL de login OAuth teria bem mais do que estas 9 linhas caso não estivéssemos usando este estilo de código.

Adicionar o alias `r` para o `import` ajuda a tornar a leitura menos cansativa.

## Planos Futuros

Adicionar outras mônadas como:
- IO
- Either
- State (A necessidade desta em Go é discutível já que não temos imutabilidade, a menos que possamos considerar *Passar por Valor* algo próximo a isso)

Os casos de uso para estas são meio raros, por isso não está incluso na biblioteca como está agora. Mas, caso seja um requisito para pelo menos alguns, será implementada.