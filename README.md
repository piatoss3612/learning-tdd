# 10장 오류 처리

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

다음으로 구현할 피처는 환율이 명시되지 않은 경우 오류 처리 개선이다.

## 오류 위시리스트

현 시점에서 환율이 누락된 경우에 대한 예외 처리가 필요하다. 이를 구현하기 위해 다음과 같은 위시리스트를 참고한다.

| 항목 | 설명                                                                                                                                |
| ---- | ----------------------------------------------------------------------------------------------------------------------------------- |
| 1    | 하나 이상 필요한 환율이 누락됐을 때 Evaluate 메서드는 명시적인 오류를 반환한다.                                                     |
| 2    | 오류 메시지는 `탐욕적`이어야 한다. 즉 첫 번째 누락된 환율 뿐만 아니라 Portfolio의 평가를 방해하는 모든 누락된 환율을 포함해야 한다. |
| 3    | 호출자가 오류를 무시하지 않도록, 오류가 발생할 때 유효한 Money가 반환돼서는 안 된다.                                                |

### 테스트 작성

환율이 누락된 경우를 명시적으로 처리하기 위해 `convert`와 `Evaluate` 메서드의 시그니처를 변경해야 한다. 현재는 단일한 반환값을 가지고 있지만, 두 번째 반환값을 추가하여 문제가 발생했음을 알리는 것이 좋다. 우선은 테스트 코드를 작성한다.

```go
func TestAdditionWithMultipleMissingExchangeRates(t *testing.T) {
	var portfolio s.Portfolio

	oneDollar := s.NewMoney(1, "USD")
	oneEuro := s.NewMoney(1, "EUR")
	oneWon := s.NewMoney(1, "KRW")

	portfolio = portfolio.Add(oneDollar)
	portfolio = portfolio.Add(oneEuro)
	portfolio = portfolio.Add(oneWon)

	expectedErrorMessage := "Missing exchange rate(s):[USD->Kalganid,EUR->Kalganid,KRW->Kalganid]"
	_, actualError := portfolio.Evaluate("Kalganid")

	if expectedErrorMessage != actualError.Error() {
		t.Errorf("Expected %v, got %v", expectedErrorMessage, actualError.Error())
	}
}
```

여러 화폐를 칼가니드로 환전하려고 시도할 때, 누락된 환율이 존재하므로 오류가 발생해야 한다. 이 때, 첫 번째 반환값은 무시하고 두 번째 반환값을 통해 오류 메시지를 확인한다. 예상되는 오류 메시지에는 누락된 환율이 모두 나열되어야 하며, ','로 구분되어야 한다.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:78:20: assignment mismatch: 2 variables but portfolio.Evaluate returns 1 value
FAIL    tdd [build failed]
FAIL
```

예상된 오류와 달리, 실제로 발생한 오류는 `Evaluate` 메서드가 현 시점에서는 단일 반환값을 가지고 있기 때문에 컴파일 단계에서 발생한 것이다. 이를 해결하기 위해 `Evaluate` 메서드의 시그니처를 변경한여 두 번째 반환값을 통해 `error`를 반환하도록 수정해야 한다.

그런데 에러가 발생하는지 어떻게 확인할 수 있을까? 이를 위해 `convert` 메서드에서 환율이 누락되었는지 아닌지를 확인하는 로직을 추가해야 한다. 또한 환율이 누락된 경우, `Evaluate` 메서드에서 확인할 수 있도록 `bool` 타입의 반환값을 추가해야 한다.

```go
func convert(m Money, currency string) (float64, bool) {
	if m.currency == currency {
		return m.amount, true
	}

	key := m.currency + "->" + currency
	exchangeRates := map[string]float64{
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	rate, ok := exchangeRates[key]

	return m.amount * rate, ok
}
```

환율이 누락된 경우, `ok` 변수가 `false`를 반환하도록 수정했다. 이제 `Evaluate` 메서드에서 `convert` 메서드의 반환값을 통해 환율이 누락된 경우를 확인할 수 있다.

```go
func (p Portfolio) Evaluate(currency string) (Money, error) {
	total := 0.0
	failedConversions := make([]string, 0, len(p))

	for _, m := range p {
		if converted, ok := convert(m, currency); ok {
			total += converted
		} else {
			failedConversions = append(failedConversions, m.currency+"->"+currency)
		}
	}

	if len(failedConversions) == 0 {
		return Money{amount: total, currency: currency}, nil
	}

	failures := fmt.Sprintf("Missing exchange rate(s):[%s]", strings.Join(failedConversions, ","))

	return Money{}, fmt.Errorf(failures)
}
```

환율이 누락된 경우, `failedConversions` 슬라이스에 환율이 누락된 화폐 쌍을 추가하고, 이를 이용해 오류 메시지를 생성한다. 이제 테스트 코드를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:32:23: assignment mismatch: 1 variable but portfolio.Evaluate returns 2 values
./money_test.go:47:17: assignment mismatch: 1 variable but portfolio.Evaluate returns 2 values
./money_test.go:61:17: assignment mismatch: 1 variable but portfolio.Evaluate returns 2 values
FAIL    tdd [build failed]
FAIL
```

이번에는 이전에 작성했던 테스트 코드에서 `Evaluate`의 두 번째 반환값을 받지 못했기 때문에 컴파일 에러가 발생했다. 이를 해결하기 위해 테스트 코드를 다음과 같이 수정한다.

```go
portfolioInDollars, _ = portfolio.Evaluate("USD")
```

이제 테스트 코드를 실행해보자.

```bash
$ go test -v .$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
=== RUN   TestDivision
--- PASS: TestDivision (0.00s)
=== RUN   TestAddition
--- PASS: TestAddition (0.00s)
=== RUN   TestAdditionDollarsAndEuros
--- PASS: TestAdditionDollarsAndEuros (0.00s)
=== RUN   TestAdditionDollarsAndWon
--- PASS: TestAdditionDollarsAndWon (0.00s)
=== RUN   TestAdditionWithMultipleMissingExchangeRates
--- PASS: TestAdditionWithMultipleMissingExchangeRates (0.00s)
PASS
ok      tdd     0.002s
```

모든 테스트가 성공적으로 통과되었다.

### 리팩터링

기존에 예상값과 실제값을 비교하기 위해 사용하던 `assertEqual` 함수는 `Money` 타입을 비교하는 데에만 사용되었으므로,

```go
func assertEqual(t *testing.T, expected, actual s.Money) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
```

새로운 테스트 코드에서는 `string` 타입을 비교하는 데에는 사용할 수 없었다.

```go
if expectedErrorMessage != actualError.Error() {
	t.Errorf("Expected %v, got %v", expectedErrorMessage, actualError.Error())
}
```

불필요한 코드 중복을 줄이기 위해 `assertEqual` 함수를 다음과 같이 수정한다.

```go
func assertEqual(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
```

빈 인터페이스 타입은 모든 타입을 받을 수 있으므로, 이제 `string` 타입을 비교하는 데에도 사용할 수 있다. 이제 에러 메시지를 비교하는 코드를 다음과 같이 수정한다.

```go
assertEqual(t, expectedErrorMessage, actualError.Error())
```

이제 테스트 코드를 실행해보자.

```bash
$ go test -v .
...
PASS
ok      tdd     0.003s
```

모든 테스트가 성공적으로 통과되었다.

그러나 아직 `convert`와 `Evaluate` 메서드에서 키를 구성하는 코드가 중복되어 있다. 이를 제거할 필요가 있으므로, 환율 구현 개선과 더불어 다음 단계에서 진행하도록 한다.

- [ ] 환율 구현 개선

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: improved error handling for missing exchange rates"
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
- [ ] 환율 구현 개선
- [ ] 환율 수정 허용
