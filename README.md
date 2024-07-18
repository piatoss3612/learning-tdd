# 8장 포트폴리오 평가하기

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [ ] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원

다음으로 구현할 피처는 혼합된 화폐를 더하는 것이다.

## 돈 섞기

다양한 화폐를 조합하려면 '한 통화에서 다른 통화로 변환'이라는 새로운 추상화가 필요하다. 이를 위한 몇 가지 기본 원칙을 세워보자.

### 1. 환전은 항상 한 쌍의 통화와 연관된다.

모든 환전이 독립적으로 이루어져야 한다. 즉, 한 통화를 다른 통화로 변환하는 환전은 항상 한 쌍의 통화와 연관된다.

### 2. 환전은 명확한 환율로써 한 통화를 다른 통화로 전환한다.

환율은 환전의 핵심 구성 요소다. 환율은 분수로 표현된다.

### 3. 한 쌍의 통화 간 두 환율은 서로 산술적 역수일 수도, 그렇지 않을 수도 있다.

유로에서 달러로의 환율은 달러에서 유로로의 환율에 수학적 역수일 수도, 그렇지 않을 수도 있다.

### 4. 한 통화에서 다른 통화로의 환율이 정의되지 않을 수도 있다.

환율이 정의되지 않은 경우, 두 통화 간의 환전은 불가능하다.

이 원칙들을 어떻게 코드로 구현할 수 있을까? 답은 테스트 주도 개발을 통해 한 번에 하나씩 구현하는 것이다.

`5달러 + 10유로 = 17달러` 피처를 구현하기 위한 테스트를 작성해보자.

```go
// money_test.go
func TestAdditionDollarsAndEuros(t *testing.T) {
	var portfolio s.Portfolio

	fiveDollars := s.NewMoney(5, "USD")
	tenEuros := s.NewMoney(10, "EUR")

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(tenEuros)

	expectedValue := s.NewMoney(17, "USD")
	actualValue := portfolio.Evaluate("USD")
	assertEqual(t, expectedValue, actualValue)
}
```

1유로 당 1.2달러를 가정하고 예상되는 결과를 17달러로 설정하였다. 이제 테스트를 실행해보자.

```bash
$ go test -v .
...
=== RUN   TestAdditionDollarsAndEuros
    money_test.go:53: Expected {17 USD}, got {15 USD}
--- FAIL: TestAdditionDollarsAndEuros (0.00s)
```

테스트가 실패했다. 예상되는 결과는 17달러지만, 실제 결과는 15달러다. 이는 Evaluate 메서드가 단순히 모든 Money 구조체의 amount를 더하는 방식으로 구현되어 있기 때문이다.

우선 각 Money의 금액을 목표 통화로 변환하는 방법을 생각해보자. 이를 위해 우선 Evaluate 메서드를 수정해보자.

```go
func (p Portfolio) Evaluate(currency string) Money {
	total := 0.0
	for _, m := range p {
		total += convert(m, currency)
	}

	return Money{amount: total, currency: currency}
}
```

convert 함수는 Money와 목표 통화를 인자로 받아 변환된 금액을 반환한다. convert 함수는 다음과 같이 구현할 수 있다.

```go
func convert(m Money, currency string) float64 {
	if m.currency == currency {
		return m.amount
	}

	return m.amount * 1.2
}
```

우선은 테스트를 통과하는 가장 간단한 방법(유로에 1.2를 곱하는 방식)으로 구현했다. 이제 테스트를 실행해보자.

```bash
$ go test -v .
...
PASS
ok      tdd     0.005s
```

테스트가 성공했다. 그러나 이 방법은 다음과 같은 문제가 있다.

1. 환율이 하드코딩되어 있다.
2. 환율이 하나의 통화에 의존한다.
3. 환율이 변하지 않는다.

우선은 첫 번째 문제를 해결하고, 나머지는 피처 목록에 추가해보자.

```go
func convert(m Money, currency string) float64 {
	eurToUsd := 1.2

	if m.currency == currency {
		return m.amount
	}

	return m.amount * eurToUsd
}
```

테스트는 여전히 그린이다.

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: conversion of Money from EUR to USD"
```

## 중간 점검

현재까지의 작업은 다음과 같다.

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [x] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원
- [ ] 연관된 통화에 기반한 환율 결정
- [ ] 환율 변경 허용

유로를 달러로 변환하는 피처를 구현했다. 그러나 변환은 한 가지 경우에만 적용되며, 환율을 추가하거나 변경할 방법이 없다. 다음으로는 연관된 통화에 기반한 환율 결정을 구현해보자.
