# 11장 은행 업무로 재설계

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [x] 5달러 + 10유로 = 17달러
- [x] 1달러 + 1100원 = 2200원
- [x] 연관된 통화에 기반한 환율 결정 (환전 전 -> 환전 후)
- [x] 환율이 명시되지 않은 경우 오류 처리 개선
- [ ] 환율 구현 개선
- [ ] 환율 수정 허용

현 시점에서 Portfolio는 Money 엔티티의 저장소 & 환율표를 관리하고 환전을 수행하는 여러 책임을 가지고 있다. 환전과 관련된 기능은 Portfolio 뿐만 아니라 다른 여러 곳에서도 사용될 수 있으므로, 이 기능을 별도의 서비스로 분리해야 한다. 현실에서 환전은 은행을 통해 이루어지므로, 이를 은행 업무로 재설계하고자 한다.

Bank 엔티티는 환율을 관리하고 환전을 수행하는 책임을 가진다. 또한 비대칭 환율을 지원하며, 환율이 누락되어 환전이 불가능한 경우에 대한 오류 처리도 담당한다.

## 의존성 주입

세 개의 엔티티는 상호 의존성을 가진다. 이는 다음과 같이 표현할 수 있다.

TODO: 관계도 이미지 추가

```plantuml
Portfolio -(aggregation)-> Money
Portfolio -(uses)-> Bank
Bank -(uses)-> Money

aggregation (집합 연관)
uses (인터페이스 의존성)
```

Portfolio는 Bank의 Convert 메서드에 의존하고 있으며, 이를 통해 Bank에 대한 의존성을 최소화할 수 있다. Bank는 Evaluate 메서드를 통해 Bank를 매개변수로 받아들이는데, 이처럼 의존성이 필요한 시점에 메서드를 통해 주입받는 방식을 메서드 주입이라고 한다.

## 모두 합치기

Bank를 이용해 하나의 Money 구조체를 다른 통화로 환전하는 테스트를 작성해보자. Bank는 아직 구현되지 않은 상태이다.

```go
func TestConversion(t *testing.T) {
	bank := s.NewBank()
	bank.AddExchangeRate("EUR", "USD", 1.2)
	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "USD")
	assertNil(t, err)
	assertEqual(t, s.NewMoney(12, "USD"), actualConvertedMoney)
}

func assertNil(t *testing.T, err error) {
	if actual != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}
```

코드의 동작 과정은 다음과 같다.

1. Bank(미구현) 인스턴스를 생성한다.
2. Bank에 환율을 추가(미구현)한다.
3. 10유로를 생성한다.
4. Bank에게 10유로를 USD로 환전하도록 요청(미구현)한다.
5. 환전된 결과에 오류가 없는지 확인한다.
6. 환전된 결과가 예상한 값과 일치하는지 확인한다.

테스트를 실행하면 당연하게도 실패한다.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:84:12: undefined: s.NewBank
FAIL    tdd [build failed]
FAIL
```

이 테스트를 '그린'으로 만들려면 '미구현' 항목들을 구현해야 한다. stocks 패키지에 bank.go 파일을 생성하고, Bank 구조체를 정의한다.

```go
package stocks

import "errors"

type Bank struct {
	exchangeRates map[string]float64
}

func (b Bank) AddExchangeRate(currencyFrom string, currencyTo string, rate float64) {
	key := currencyFrom + "->" + currencyTo
	b.exchangeRates[key] = rate
}

func (b Bank) Convert(money Money, currencyTo string) (Money, error) {
	if money.currency == currencyTo {
		return money, nil
	}

	key := money.currency + "->" + currencyTo
	rate, ok := b.exchangeRates[key]
	if !ok {
		return Money{}, errors.New("failed")
	}

	return Money{amount: money.amount * rate, currency: currencyTo}, nil
}

func NewBank() Bank {
	return Bank{exchangeRates: make(map[string]float64)}
}
```

Bank 구조체는 환율을 저장하는 exchangeRates 필드를 가지고 있다. AddExchangeRate 메서드는 환율을 추가하는 메서드이며, Convert 메서드는 환전을 수행하는 메서드이다. Convert 메서드는 환전할 금액과 환전할 통화를 매개변수로 받아들인다. 환전할 금액의 통화가 환전할 통화와 같은 경우에는 환전할 필요가 없으므로, 환전할 금액을 그대로 반환한다. 환전할 금액의 통화가 환전할 통화와 다른 경우에는 exchangeRates 필드에서 환율을 찾아 환전을 수행한다. 환율이 존재하지 않는 경우에는 오류를 반환한다.

이제 테스트를 실행하면 성공한다.

```bash
$ go test -v .
...
=== RUN   TestConversion
--- PASS: TestConversion (0.00s)
PASS
ok      tdd     0.004s
```

그런데 Portfolio의 기존의 Evaluate 메서드는 변환을 수행하는 과정에서 누락된 모든 환율을 담은 오류를 반환하도록 구현되어 있다. ex) "Missing exchange rate(s):[USD->EUR]". 누락된 환율은 어느 부분에서 처리해야 할까? 기존의 Evaluate 메서드에서 convert 메서드를 호출하는 부분에서 그러했듯이, Bank의 Convert 메서드에서 처리하는 것이 가장 적절하다. 지금은 Bank의 Convert 메서드에서 누락된 환율을 처리하지 않고 있다. 또한 오류가 발생했을 때, 빈 Money 구조체가 반환되는데, 이보다는 nil을 반환하는 것이 더 적절하다.

이를 위한 테스트를 우선 작성해보자.

```go
func TestConversionWithMissingExchangeRate(t *testing.T) {
	bank := s.NewBank()
	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "Kalganid")
	if actualConvertedMoney != nil {
		t.Errorf("Expected money to be nil, got %v", actualConvertedMoney)
	}
	assertEqual(t, "EUR->Kalganid", err.Error())
}
```

이 테스트는 다음과 같은 동작을 수행한다.

1. Bank 인스턴스를 생성한다.
2. 10유로를 생성한다.
3. Bank에게 10유로를 Kalganid로 환전하도록 요청한다.
4. 오류가 발생했을 때, 환전된 결과가 nil인지 확인한다.
5. 오류가 발생했을 때, 오류 메시지가 예상한 값과 일치하는지 확인한다.

이 테스트를 실행하면 실패한다.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:96:29: invalid operation: actualConvertedMoney != nil (mismatched types stocks.Money and untyped nil)
FAIL    tdd [build failed]
FAIL
```

이 테스트를 '그린'으로 만들기 위해 Bank의 Convert 메서드의 시그니처를 변경하여 첫 번째 반환값을 포인터로 변경하고, 오류가 발생했을 때 nil을 반환하도록 수정한다.

```go
func (b Bank) Convert(money Money, currencyTo string) (*Money, error) {
	if money.currency == currencyTo {
		return &money, nil
	}

	key := money.currency + "->" + currencyTo
	rate, ok := b.exchangeRates[key]
	if !ok {
		return nil, errors.New("failed")
	}

	return &Money{amount: money.amount * rate, currency: currencyTo}, nil
}
```

다시 테스트를 실행하면 이번에는 다음과 같이 실패한다.

```bash
=== RUN   TestConversion
    money_test.go:110: Expected {12 USD}, got &{12 USD}
--- FAIL: TestConversion (0.00s)
=== RUN   TestConversionWithMissingExchangeRate
    money_test.go:110: Expected EUR->Kalganid, got failed
--- FAIL: TestConversionWithMissingExchangeRate (0.00s)
```

이전에 작성한 테스트가 실패했는데, 이는 Convert 메서드가 nil을 반환하도록 수정되었기 때문이다. 이를 해결하기 위해 테스트에서 반환값의 포인터를 역참조하여 비교하도록 수정한다.

```go
assertEqual(t, s.NewMoney(12, "USD"), *actualConvertedMoney)
```

다음으로 오류 메시지를 비교하는 부분이 실패하는데, 이는 Bank의 Convert 메서드에서 반환하는 오류 메시지에 key를 반환하도록 수정하면 해결된다.

```go
func (b Bank) Convert(money Money, currencyTo string) (*Money, error) {
	if money.currency == currencyTo {
		return &money, nil
	}

	key := money.currency + "->" + currencyTo
	rate, ok := b.exchangeRates[key]
	if !ok {
		return nil, errors.New(key)
	}

	return &Money{amount: money.amount * rate, currency: currencyTo}, nil
}
```

이제 테스트를 실행하면 성공한다.

```bash
$ go test -v .
...
=== RUN   TestConversionWithMissingExchangeRate
--- PASS: TestConversionWithMissingExchangeRate (0.00s)
PASS
ok      tdd     0.002s
```

현 시점에서 리팩토링을 수행할 수 있는 부분은 assertNil 함수이다. assertNil 함수는 error 타입이 nil인지 확인하는 경우에만 사용되는데, 이를 더 일반적인 용도로 변경하면 더 유용하게 사용할 수 있다.

```go
func assertNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Errorf("Expected nil, got %v", actual)
	}
}
```

이렇게 변경된 assertNil 함수를 사용해 TestConversionWithMissingExchangeRate 테스트를 수정해보자.

```go
func TestConversionWithMissingExchangeRate(t *testing.T) {
	bank := s.NewBank()
	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "Kalganid")
	assertNil(t, actualConvertedMoney)
	assertEqual(t, "EUR->Kalganid", err.Error())
}
```

그런데 테스트를 실행하면 어딘가 이상한 실패가 발생한다.

```bash
$ go test -v .
...
=== RUN   TestConversionWithMissingExchangeRate
    money_test.go:105: Expected nil, got <nil>
--- FAIL: TestConversionWithMissingExchangeRate (0.00s)
```

이 실패는 오류 메시지의 두번째 nil을 감싸고 있는 홑화살괄호 <>에서 단서를 찾을 수 있다.

Go의 인터페이스는 타입 T와 값 V를 가지는데, V가 nil이더라도 T는 인터페이스가 나타내는 타입의 포인터를 가질 수 있다.

```go
assertNil(t, actualConvertedMoney)
```

결과적으로 actualConvertedMoney는 nil이 아니라 \*Money 타입이므로, nil과 비교하는 것이 잘못되었다. 이를 해결하기 위해 assertNil 함수를 다음과 같이 수정한다.

```go
func assertNil(t *testing.T, actual interface{}) {
	if actual != nil && !reflect.ValueOf(actual).IsNil() {
		t.Errorf("Expected nil, got %v", actual)
	}
}
```

이 함수는 actual이 nil이 아니고, reflect 패키지의 ValueOf 함수를 통해 actual의 값이 nil이 아닌지 확인한다. 이렇게하면 인터페이스의 값이 nil이 아닌 경우에만 오류를 반환한다.

이제 테스트를 실행하면 성공한다.

```bash
$ go test -v .
...
=== RUN   TestConversionWithMissingExchangeRate
--- PASS: TestConversionWithMissingExchangeRate (0.00s)
PASS
ok      tdd     0.003s
```

이제 Portfolio의 Evaluate 메서드에 Bank를 주입하도록 수정해보자.

```go
func (p Portfolio) Evaluate(bank Bank, currency string) (*Money, error) {
	total := 0.0
	failedConversions := make([]string, 0, len(p))

	for _, m := range p {
		if converted, err := bank.Convert(m, currency); err == nil {
			total += converted.amount
		} else {
			failedConversions = append(failedConversions, err.Error())
		}
	}

	if len(failedConversions) == 0 {
		return &Money{amount: total, currency: currency}, nil
	}

	failures := fmt.Sprintf("Missing exchange rate(s):[%s]", strings.Join(failedConversions, ","))

	return nil, fmt.Errorf(failures)
}
```

Evaluate 메서드는 Bank를 매개변수로 받아들이도록 변경되었다. 주입된 Bank를 통해 Convert 메서드를 호출하여 환전을 수행한다. 환전에 실패한 경우에는 failedConversions 슬라이스에 오류 메시지를 추가하고, 모든 환전이 성공한 경우에는 총 금액을 반환한다. 이때, 총 금액은 Money 구조체로 반환되며, 오류가 발생한 경우에는 오류 메시지를 반환한다.

또한 반환값도 변경되어 포인터를 반환하도록 수정되었다.

이제 테스트를 실행해보자.

```go
$ go test -v .
# tdd [tdd.test]
./money_test.go:33:26: cannot use portfolio.Evaluate("USD") (value of type *stocks.Money) as stocks.Money value in assignment
./money_test.go:33:45: not enough arguments in call to portfolio.Evaluate
        have (string)
        want (stocks.Bank, string)
./money_test.go:48:39: not enough arguments in call to portfolio.Evaluate
        have (string)
        want (stocks.Bank, string)
./money_test.go:62:39: not enough arguments in call to portfolio.Evaluate
        have (string)
        want (stocks.Bank, string)
./money_test.go:79:39: not enough arguments in call to portfolio.Evaluate
        have (string)
        want (stocks.Bank, string)
FAIL    tdd [build failed]
FAIL
```

테스트가 실패하는데, 이는 Portfolio의 Evaluate 메서드의 시그니처가 변경되었기 때문이다. 이를 해결하기 위해 테스트를 수정해야 한다. 그런데 매번 Bank를 생성하는 것은 번거롭기 때문에, 테스트 전체에서 사용할 수 있는 Bank 인스턴스를 생성하도록 수정한다.

```go
var bank s.Bank

func init() {
	bank = s.NewBank()
	bank.AddExchangeRate("EUR", "USD", 1.2)
	bank.AddExchangeRate("USD", "KRW", 1100)
}
```

전역 변수 bank를 선언하고, init 함수를 통해 bank를 초기화한다. 이제 테스트에서 bank를 사용할 수 있다. 테스트를 수정해보자.

```go
func TestAddition(t *testing.T) {
	var portfolio s.Portfolio

	fiveDollars := s.NewMoney(5, "USD")
	tenDollars := s.NewMoney(10, "USD")
	fifteenDollars := s.NewMoney(15, "USD")

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(tenDollars)
	portfolioInDollars, err := portfolio.Evaluate(bank, "USD")
	assertNil(t, err)
	assertEqual(t, fifteenDollars, *portfolioInDollars)
}
```

신경 써야 할 부분은 bank를 Evaluate 메서드에 주입하는 것 뿐만 아니라, assertEqual 함수를 사용할 때도 포인터를 역참조하여 비교해야 한다. 또한 Evaluate 메서드에서 반환된 오류가 nil인지 확인하는 부분도 추가해야 한다. 다른 테스트도 동일한 방식으로 수정한다.

그리고 더 이상 사용되지 않는 Portfolio의 convert 메서드를 제거한다.

이제 Money, Bank, Portfolio의 책임이 명확하게 분리되었다. Money는 금액과 통화를 관리하고, Bank는 환율을 관리하고 환전을 수행하며, Portfolio는 Money와 Bank를 이용해 포트폴리오의 총 금액을 계산한다.

마지막으로 테스트를 실행해보면 성공하는 것을 확인할 수 있다.

```bash
$ go test -v .
...
PASS
ok      tdd     0.003s
```

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: added Bank; refactored Portfolio to use Bank"
```

## 중간 점검

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [x] 5달러 + 10유로 = 17달러
- [x] 1달러 + 1100원 = 2200원
- [x] 연관된 통화에 기반한 환율 결정 (환전 전 -> 환전 후)
- [x] 환율이 명시되지 않은 경우 오류 처리 개선
- [x] 환율 구현 개선
- [ ] 환율 수정 허용
