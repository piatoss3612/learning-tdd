# 2장 다양한 통화로 돈 계산

## 유로에 발 들이기

- [x] 5달러 \* 2 = 10달러
- [] 10유로 \* 2 = 20유로
- [] 4002원 / 4 = 1000.5원
- [] 5달러 + 10유로 = 17달러
- [] 1달러 + 1100원 = 2200원

피처 목록 중 두 번째 항목을 구현해보자.

달러에 더해서 유로를 지원하려면 1장에서 만든 Dollar보다 더 일반적인 엔티티, 이미 정의된 amount에 더해 currency 속성을 가진 새로운 엔티티가 필요하다.

```go
func TestMultiplicationInEuros(t *testing.T) {
	tenEuros := Money{
		amount:   10,
		currency: "EUR",
	}
	twentyEuros := tenEuros.times(2)
	if twentyEuros.amount != 20 {
		t.Errorf("Expected 20, got %d", twentyEuros.amount)
	}
	if twentyEuros.currency != "EUR" {
		t.Errorf("Expected EUR, got %s", twentyEuros.currency)
	}
}
```

테스트는 금액(amount) 뿐만 아니라 통화(currency)를 포함하는 구조체 인스턴스인 '10 EUR'과 '20 EUR'의 개념을 표현한다. 이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:26:14: undefined: Money
FAIL    tdd [build failed]
FAIL
```

Money가 정의되지 않았다는 에러가 발생한다. Money 구조체를 정의해보자.

```go
type Money struct {
	amount   int
	currency string
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:35:26: tenEuros.times undefined (type Money has no field or method times)
FAIL    tdd [build failed]
FAIL
```

times 메서드가 정의되지 않았다는 에러가 발생한다. times 메서드를 정의해보자.

```go
func (m Money) times(multiplier int) Money {
	return Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestMultiplicationInEuros
--- PASS: TestMultiplicationInEuros (0.00s)
PASS
ok      tdd     0.002s
```

테스트가 그린 상태가 되었다.

## DRY한 코드를 유지하라

- DRY(Don't Repeat Yourself): 반복하지 말라, 중복을 피하라

테스트를 통과하기 위해 새로운 구조체 Money를 만들었지만, 이 구조체는 Dollar와 유사하다. 두 구조체를 합칠 수 있을 것 같다. Money 구조체가 Dollar 구조체를 포함하는 관계이므로 Dollar 구조체를 Money 구조체로 대체할 수 있다. Dollar 구조체를 제거하고 Money 구조체로 대체해보자.

```go
func TestMultiplication(t *testing.T) {
	fiver := Money{
		amount:   5,
		currency: "USD",
	}
	tenner := fiver.times(2)
	if tenner.amount != 10 {
		t.Errorf("Expected 10, got %d", tenner.amount)
	}
	if tenner.currency != "USD" {
		t.Errorf("Expected USD, got %s", tenner.currency)
	}
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestMultiplicationInEuros
--- PASS: TestMultiplicationInEuros (0.00s)
PASS
ok      tdd     0.002s
```

테스트가 그린 상태가 되었다.

## 반복하지 말라고 하지 않았나?

- currency 속성을 추가하고 Money 구조체를 만드는 과정에서 두 개의 테스트가 비슷한 코드를 포함하고 있다.
- 테스트 중 하나를 삭제할 수도 있고 그냥 내버려 둘 수도 있지만, 코드에 발생될 리그레션을 방지하기 위한 대비책이 필요하다.
- 우선은 두 테스트를 모두 유지하도록 하자.

> 리그레션: 소프트웨어의 변경으로 인해 기존의 기능이 올바르게 작동하지 않게 되는 현상

## 분할 정복

다음 요구사항인 나눗셈을 구현해보자. 이번에는 두 가지의 새로운 하위 요구사항이 있다.

1. 새로운 통화 단위: 대한민국 원(KRW)
2. 소수부를 포함하는 금액

```go
func TestDivision(t *testing.T) {
	originalMoney := Money{amount: 4002, currency: "KRW"}
	actualMoneyAfterDivision := originalMoney.Divide(4)
	expectedMoneyAfterDivision := Money{amount: 1000.5, currency: "KRW"}
	if actualMoneyAfterDivision != expectedMoneyAfterDivision {
		t.Errorf("Expected %v, got %v", expectedMoneyAfterDivision, actualMoneyAfterDivision)
	}
}
```

이전과 다르게 각 필드를 비교하지 않고 예상되는 Money 인스턴스와 실제 Money 인스턴스를 비교하고 있다. 이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:47:44: originalMoney.divide undefined (type Money has no field or method divide)
./money_test.go:48:46: cannot use 1000.5 (untyped float constant) as int value in struct literal (truncated)
FAIL    tdd [build failed]
FAIL
```

divide 메서드가 정의되지 않았다는 에러와 1000.5가 int 타입이 아니라는 에러가 발생한다.

Divide 메서드를 정의해보자.

```go
func (m Money) Divide(divisor int) Money {
	return Money{
		amount:   m.amount / divisor,
		currency: m.currency,
	}
}
```

다음으로 amount 필드가 소수부를 포함할 수 있도록 타입을 변경해보자.

```go
type Money struct {
	amount   float64
	currency string
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:26:13: invalid operation: m.amount * multiplier (mismatched types float64 and int)
./money_test.go:33:13: invalid operation: m.amount / divisor (mismatched types float64 and int)
FAIL    tdd [build failed]
FAIL
```

이번에는 float64 타입의 amount와 int 타입의 multiplier, divisor 간의 연산이 불가능하다는 에러가 발생한다. 모든 피연산자가 같은 타입(float64)을 가지도록 수정해보자.

```go
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

이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd
# [tdd]
./money_test.go:12:3: (*testing.common).Errorf format %d has arg tenner.amount of wrong type float64
./money_test.go:45:3: (*testing.common).Errorf format %d has arg twentyEuros.amount of wrong type float64
FAIL    tdd [build failed]
FAIL
```

테스트가 실패했다. 이번에는 기존의 amount가 int 타입이었던 것과 달리 float64 타입이 되었기 때문에 테스트 코드에서 값을 출력하기 위해 사용한 포맷 문자열(%d)이 float64 타입에 맞지 않아서 에러가 발생한 것이다. 포맷 문자열을 모든 타입의 값을 출력할 수 있는 %v로 수정해보자.

```go
func TestMultiplication(t *testing.T) {
	fiver := Money{
		amount:   5,
		currency: "USD",
	}
	tenner := fiver.times(2)
	if tenner.amount != 10 {
		t.Errorf("Expected 10, got %v", tenner.amount)
	}
	if tenner.currency != "USD" {
		t.Errorf("Expected USD, got %s", tenner.currency)
	}
}
```

다음으로 테스트를 실행해보자.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestMultiplicationInEuros
--- PASS: TestMultiplicationInEuros (0.00s)
=== RUN   TestDivision
--- PASS: TestDivision (0.00s)
PASS
ok      tdd     0.002s
```

테스트가 그린 상태가 되었다. 이제 테스트를 통과하기 위해 구현한 코드를 리팩토링해보자.

## 마무리하기

각 테스트에서 공통된 코드를 찾아내고, 중복을 제거하면서 테스트를 통과하는 코드를 작성해보자.

예상되는 Money 인스턴스와 실제 Money 인스턴스를 비교하는 코드를 함수로 추출해보자.

```go
func assertEqual(t *testing.T, expected, actual Money) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
```

이제 테스트 코드를 리팩토링해보자.

```go
func TestMultiplication(t *testing.T) {
	fiver := Money{
		amount:   5,
		currency: "USD",
	}
	tenner := fiver.times(2)
	expectedTenner := Money{
		amount:   10,
		currency: "USD",
	}
	assertEqual(t, expectedTenner, tenner)
}

func TestMultiplicationInEuros(t *testing.T) {
	tenEuros := Money{
		amount:   10,
		currency: "EUR",
	}
	twentyEuros := tenEuros.times(2)
	expectedTwentyEuros := Money{
		amount:   20,
		currency: "EUR",
	}
	assertEqual(t, expectedTwentyEuros, twentyEuros)
}

func TestDivision(t *testing.T) {
	originalMoney := Money{amount: 4002, currency: "KRW"}
	actualMoneyAfterDivision := originalMoney.Divide(4)
	expectedMoneyAfterDivision := Money{amount: 1000.5, currency: "KRW"}
	assertEqual(t, expectedMoneyAfterDivision, actualMoneyAfterDivision)
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
PASS
ok      tdd     0.002s
```

테스트가 그린 상태가 되었다. 이제 변경 사항을 반영하고 마무리하자.

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: division and multiplication features done"
```

## 중간 점검

1. 달러와 유로를 모두 지원하는 Money 구조체를 만들었다.
2. 나눗셈을 구현했고 실수를 사용할 수 있도록 설계를 변경했다.
3. 테스트 코드를 리팩토링하여 중복을 제거했다.

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.5원
- [] 5달러 + 10유로 = 17달러
- [] 1달러 + 1100원 = 2200원
