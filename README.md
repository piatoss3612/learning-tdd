# 9장 여기도 통화, 저기도 통화

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [x] 5달러 + 10달러 = 15달러
- [x] 프로덕션 코드와 테스트 코드 분리
- [x] 중복된 테스트 제거
- [x] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원

`Portfolio`의 `Money` 엔티티에 대한 평가(Evaluate) 피처의 현 상태는 다음과 같다.

1. `Money`의 통화와 같은 통화인 경우, `Money`의 `amount`를 반환한다. 이는 올바른 동작이며, 어떤 통화든 그 자신에 대한 환율은 1이다.
2. 같은 통화가 아닌 다른 모든 경우, 고정된 숫자(1.2)를 곱한 값을 반환한다. 이는 'USD'를 'EUR'로 환전하는 경우에 대해서만 올바른 동작이다. 또한, 이 환율을 수정할 수 없다는 문제가 있다.

이번 장에서는 통화별 환율을 사용해 한 통화에서 다른 통화로 환전하는 기능을 구현해보자.

## 해시맵 만들기

'환전 전' 통화와 '환전 후' 통화로 주어진 환율을 조회할 수 있는 해시맵이 필요하다. 이는 환율표를 나타낸다.

| 환전 전 | 환전 후 | 환율    |
| ------- | ------- | ------- |
| EUR     | USD     | 1.2     |
| USD     | EUR     | 0.82    |
| USD     | KRW     | 1100    |
| KRW     | USD     | 0.00090 |

환율표를 통해 알 수 있듯이, 통화 쌍에 대한 상호간 환율은 서로 산술적 역수가 아닐 수 있다. 예를 들어, '100 EUR'를 'USD'로 환전하면 '120 USD'가 되지만, '120 USD'를 'EUR'로 환전하면 '98.4 EUR'가 된다. 이는 은행이 돈을 버는 방법 중 하나이다.

환율표를 구현하기 위해 피처 목록에 다음 두 가지 항목을 추가한다.

- [ ] 연관된 통화에 기반한 환율 결정 (환전 전 -> 환전 후)
- [ ] 환율 수정 허용

또한 새로운 통화(한국 원화)를 추가함에 따라 우선순위 전제 행위 변환(TPP, Transformation Priority Premise)이 실행되는 것을 확인할 수 있다. 즉, `if-else` 브랜치를 여러 가닥으로 나누면서 구조를 변경하지 않고, 새로운 자료 구조를 도입하여 행위를 변경하는 것이다.

- [ ] 1달러 + 1100원 = 2200원

### 테스트 작성

```go
func TestAdditionDollarsAndWon(t *testing.T) {
	var portfolio s.Portfolio

	fiveDollars := s.NewMoney(1, "USD")
	elevenHundredWon := s.NewMoney(1100, "KRW")

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(elevenHundredWon)

	expectedValue := s.NewMoney(2200, "KRW")
	actualValue := portfolio.Evaluate("KRW")

	assertEqual(t, expectedValue, actualValue)
}
```

달러와 원화를 더하는 테스트를 작성한다. 기댓값은 2200원으로, 1 달러당 1100원을 가정한다.

이제 테스트를 실행하면 실패할 것이다. 이는 아직 올바른 환율을 적용하는 메커니즘이 구현되지 않았기 때문이다.

```bash
$ go test .
...
=== RUN   TestAdditionDollarsAndWon
    money_test.go:68: Expected {2200 KRW}, got {1101.2 KRW}
--- FAIL: TestAdditionDollarsAndWon (0.00s)
...
```

환율을 나타내는 `map[string]float64` 타입의 `exchangeRates`를 `Portfolio` 타입의 `convert` 메서드에 추가한다.

```go
exchangeRates := map[string]float64{
	"EUR->USD": 1.2,
	"USD->KRW": 1100,
}
```

이제 `eurToUsd`를 제거하고 `m.currency`와 `currency`를 이용해 `key`를 생성한다. 그리고 `exchangeRates`에서 `key`를 이용해 환율을 조회한다. 이를 이용해 환전을 수행한다.

```go
func convert(m Money, currency string) float64 {
	if m.currency == currency {
		return m.amount
	}

	key := m.currency + "->" + currency
	exchangeRates := map[string]float64{
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	return m.amount * exchangeRates[key]
}
```

이제 테스트를 실행하면 성공할 것이다.

```bash
$ go test .
...
=== RUN   TestAdditionDollarsAndWon
--- PASS: TestAdditionDollarsAndWon (0.00s)
...
```

그런데 이렇게 하면 환율을 수정할 수 없다. 만약 `exchangeRates`에 명시되지 않은 환율을 사용하려고 한다면, 무슨 일이 벌어질까? `exchangeRates`의 모든 키-값 쌍을 주석 처리하고 테스트를 실행해보자.

```go
exchangeRates := map[string]float64{
	// "EUR->USD": 1.2,
	// "USD->KRW": 1100,
}
```

```bash
$ go test .
...
=== RUN   TestAdditionDollarsAndEuros
    money_test.go:68: Expected {17 USD}, got {5 USD}
--- FAIL: TestAdditionDollarsAndEuros (0.00s)
=== RUN   TestAdditionDollarsAndWon
    money_test.go:68: Expected {2200 KRW}, got {1100 KRW}
--- FAIL: TestAdditionDollarsAndWon (0.00s)
...
```

`TestAdditionDollarsAndEuros`와 `TestAdditionDollarsAndWon` 테스트가 실패한다. 이는 환율을 찾을 수 없기 때문에 `float64`의 기본값인 0과 곱셈을 하기 때문이다. 더 나은 예외 처리 방법이 필요하므로, 이를 피처 목록에 추가한다.

- [ ] 환율이 명시되지 않은 경우 오류 처리 개선

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: conversion between currencies with defined exchange rates"
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
- [ ] 환율이 명시되지 않은 경우 오류 처리 개선
- [ ] 환율 수정 허용
