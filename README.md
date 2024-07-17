# 5장 Go의 패키지 및 모듈

## 코드를 패키지로 분리하기

현재 `money_test.go` 파일에 테스트 코드와 프로덕션 코드가 섞여 있다. 이를 분리해보자. 이는 다음 두 가지의 분리 작업을 수반한다.

1. 프로덕션 코드와 테스트 코드를 분리한다.
2. 테스트 코드에서 프로덕션 코드로만 의존성을 가지도록 한다.

먼저 `money.go` 파일과 `portfolio.go` 파일을 생성한다.

```go
// money.go
package main

type Money struct {
	amount   float64
	currency string
}

func (m Money) times(multiplier float64) Money {
	return Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}

func (m Money) Divide(divisor float64) Money {
	return Money{
		amount:   m.amount / divisor,
		currency: m.currency,
	}
}
```

```go
// portfolio.go
package main

type Portfolio []Money

func (p Portfolio) Add(money Money) Portfolio {
	return append(p, money)
}

func (p Portfolio) Evaluate(currency string) Money {
	total := 0.0
	for _, m := range p {
		total += m.amount
	}

	return Money{amount: total, currency: currency}
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestMultiplicationInEuros
--- PASS: TestMultiplicationInEuros (0.00s)
=== RUN   TestDivision
--- PASS: TestDivision (0.00s)
=== RUN   TestAddition
--- PASS: TestAddition (0.00s)
PASS
ok      tdd     0.002s
```

소스 코드를 별도의 파일로 분리했지만, 상위 레벨에서 코드의 구성은 어떤가? Portfolio와 Money가 모두 main 패키지에 속해 있다. 이를 '주식' 마켓과 관련된 네임스페이스로 묶어보고 싶다. 그 전에 Go의 모듈 시스템에 대해 알아보자.

## Go 모듈

Go 프로그램은 일반적으로 여러 개의 소스 파일로 구성되며, 각 소스 파일은 `package` 키워드를 사용하여 자신이 속한 패키지를 선언한다. 그리고 이러한 패키지들은 하나의 모듈로 묶일 수 있다.

애플리케이션으로 실행되어야 하는 Go 프로그램은 `main` 패키지에 위치해야 한다. 이 패키지는 `main` 함수를 가져야 하며, 이 함수는 프로그램의 시작점이 된다.

현재까지 작성된 프로그램은 다음과 같은 구조를 가지고 있다.

```bash
.
├── go.mod
├── money.go
├── money_test.go
└── portfolio.go
```

`go mod init tdd`를 실행해 Go 모듈을 초기화하고 이 과정에서 `go.mod` 파일이 생성된다.

tdd 모듈 내부에는 모든 코드가 main 패키지에 속해 있다. 이를 수정해보자.

## 패키지 생성하기

`money.go`와 `portfolio.go` 파일을 `stocks` 패키지로 분리해보자.

폴더 구조는 다음과 같다.

```bash
.
├── go.mod
├── money_test.go
└── stocks
    ├── money.go
    └── portfolio.go
```

`money.go` 파일과 `portfolio.go` 파일의 패키지 선언을 `stocks`로 변경한다.

```go
// money.go
package stocks
```

```go
// portfolio.go
package stocks
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:6:11: undefined: Money
./money_test.go:11:20: undefined: Money
./money_test.go:19:14: undefined: Money
./money_test.go:24:25: undefined: Money
./money_test.go:32:19: undefined: Money
./money_test.go:34:32: undefined: Money
./money_test.go:38:49: undefined: Money
./money_test.go:45:16: undefined: Portfolio
./money_test.go:46:25: undefined: Money
./money_test.go:48:17: undefined: Money
./money_test.go:48:17: too many errors
FAIL    tdd [build failed]
FAIL
```

테스트가 실패했다. 이는 테스트 코드에서 `Money`와 `Portfolio` 타입을 찾을 수 없기 때문이다. 이를 해결하기 위해 `money_test.go` 파일에서 `stocks` 패키지를 임포트해야 한다.

이 때 stocks 패키지의 경로는 `tdd/stocks`이다. 따라서 `money_test.go` 파일의 상단에 다음과 같이 임포트한다. `s`는 stocks 패키지의 별칭이다.

```go
// money_test.go
package main

import
(
	"testing"
	s "tdd/stocks"
)
```

이제 Money와 Portfolio의 모든 참조를 `s.Money`와 `s.Portfolio`로 변경한다.

```go
func assertEqual(t *testing.T, expected, actual s.Money)
```

다음으로 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:10:3: unknown field amount in struct literal of type stocks.Money
./money_test.go:11:3: unknown field currency in struct literal of type stocks.Money
./money_test.go:13:18: fiver.times undefined (type stocks.Money has no field or method times)
./money_test.go:15:3: unknown field amount in struct literal of type stocks.Money
./money_test.go:16:3: unknown field currency in struct literal of type stocks.Money
./money_test.go:23:3: unknown field amount in struct literal of type stocks.Money
./money_test.go:24:3: unknown field currency in struct literal of type stocks.Money
./money_test.go:26:26: tenEuros.times undefined (type stocks.Money has no field or method times)
./money_test.go:28:3: unknown field amount in struct literal of type stocks.Money
./money_test.go:29:3: unknown field currency in struct literal of type stocks.Money
./money_test.go:29:3: too many errors
FAIL    tdd [build failed]
FAIL
```

테스트가 실패했다. 이번에는 Money와 Portfolio의 필드와 메서드를 찾을 수 없다는 에러가 발생했다. 뭐가 문제일까?

## 캡슐화

이전에는 모든 코드가 main 패키지에 속해 있었기 때문에, 테스트에서 Money와 Portfolio의 필드와 메서드에 접근할 수 있었다. 하지만 이제는 stocks 패키지에 속해 있기 때문에, 테스트에서 Money와 Portfolio의 필드와 메서드에 접근할 수 없다. Go는 대소문자로 구분하여 접근 제어를 구현한다. 여기서는 Money와 Portfolio의 필드와 메서드가 소문자로 시작하므로, 이를 테스트에서 접근할 수 없다.

Go는 이처럼 패키지 내부에서만 접근할 수 있는 멤버를 캡슐화할 수 있다.

이를 해결하기 위해 Money와 Portfolio의 필드와 메서드를 대문자로 시작하도록 변경할 수도 있지만, 부가적인 행동을 제공하여 캡슐화를 유지하고 코드의 불변성을 보장할 수 있다.

그 방법은 다음과 같이 생성자 함수를 사용하는 것이다.

```go
func NewMoney(amount float64, currency string) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}
```

> Money 구조체의 times 메서드를 Times로 바꿔야 하는데, 이 부분은 처음에 의도되지 않은 오타인 것 같다.

이제 생성자 함수를 사용하여 테스트 코드를 수정해보자.

```go
fiveDollars := s.NewMoney(5, "USD")
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestMultiplicationInEuros
--- PASS: TestMultiplicationInEuros (0.00s)
=== RUN   TestDivision
--- PASS: TestDivision (0.00s)
=== RUN   TestAddition
--- PASS: TestAddition (0.00s)
PASS
ok      tdd
```

테스트가 성공했다.

## 테스트에서 중복 제거하기

`TestMultiplicationInEuros`는 `TestMultiplication`과 동일한 테스트이다. 따라서 이를 제거한다. 그러고 나면 테스트 파일에는 세 개의 테스트만 남게 된다.

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "refactor: moved Money and Portfolio to stocks package"
```

## 중간 점검

현재까지의 작업은 다음과 같다.

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [ ] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원

이번 장에서는 프로덕션 코드와 테스트 코드를 분리하고, 패키지를 생성하여 코드를 캡슐화했다. 또한 중복된 테스트를 제거하고, 테스트에서 프로덕션 코드로 단방향 의존성을 유지하도록 변경했다.
