# 12장 테스트 순서

11장에서 Bank 엔티티를 추가하고 환율을 관리 및 환전하는 기능을 구현했다. 이 과정에서 환율을 수정하는 피처가 이미 포함되어 있다고 믿을 수 있다. 그러나 실제로 의도한 대로 동작하는지 확신을 얻으려면 테스트를 작성하여 이를 확인해야 한다.

왜 이미 구현한 피처에 대한 테스트를 작성해야 할까? 이유는 다음과 같다.

1. 반복적인 테스트를 통해 코드를 검증하면 코드의 신뢰도를 높일 수 있다.
2. 새로운 테스트는 해당 피처의 실행 가능한 문서 역할을 한다.
3. 새로운 테스트는 기존 테스트가 생각하지 못한 상호작용을 발견할 수 있도록 도와줌으로써 이를 해결하도록 유도한다.

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

## 환율 변경

환전에 대한 기존 테스트를 변경하는 것에서 시작해보자. 이전 환율을 통해 환전을 수행하고, 환율을 변경한 후에도 예상대로 동작하는지 확인해야 한다.

TestConversion 테스트 함수의 마지막에 환율을 수정하고 다시 환전을 수행하는 코드를 추가한다.

```go
func TestConversion(t *testing.T) {
	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "USD")
	assertNil(t, err)
	assertEqual(t, s.NewMoney(12, "USD"), *actualConvertedMoney)

	// Add a new exchange rate
	bank.AddExchangeRate("EUR", "USD", 1.3)
	actualConvertedMoney, err = bank.Convert(tenEuros, "USD")
	assertNil(t, err)
	assertEqual(t, s.NewMoney(13, "USD"), *actualConvertedMoney)
}
```

이제 테스트를 실행하면 성공한다!

```bash
$ go test -v .
...
=== RUN   TestConversion
--- PASS: TestConversion (0.00s)
```

테스트를 수정했으니, 리팩토링을 통해 테스트 이름을 변경하여 테스트의 의도를 명확하게 드러내자.

```go
func TestConversionWithDifferentRatesBetweenTwoCurrencies(t *testing.T) {
	...
}
```

그런데 앞서 11장에서 Bank 인스턴스를 전역 변수로 선언했기 때문에 이 테스트가 다른 테스트에 영향을 줄 수 있다. 이를 확인하기 위해 바로 아래에 다른 테스트를 추가하고 실행해보자.

```go
func TestWhatIsTheConversionRateFromEURToUSD(t *testing.T) {
	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "USD")
	assertNil(t, err)
	assertEqual(t, s.NewMoney(12, "USD"), *actualConvertedMoney)
}
```

```bash
$ go test -v .
...
=== RUN   TestWhatIsTheConversionRateFromEURToUSD
    money_test.go:132: Expected {12 USD}, got {13 USD}
```

예상대로 변경된 환율에 따라 테스트가 실패했다. 이는 init 함수가 각 테스트 함수가 실행되기 전에 매번 호출되는 것이 아닌, 프로그램이 시작될 때 한 번만 호출되기 때문이다.

문제를 정리하면 다음과 같다.

1. bank 인스턴스가 여러 테스트에서 공유된다.
2. 각 테스트 실행 전 bank의 상태를 초기화하지 않는다.

이를 해결하기 위해 각 테스트 함수가 실행되기 전에 bank 인스턴스를 초기화할 수 있는 함수를 작성해보자.

```go
func initExchangeRates() {
	bank = s.NewBank()
	bank.AddExchangeRate("EUR", "USD", 1.2)
	bank.AddExchangeRate("USD", "KRW", 1100)
}
```

그리고 이 함수를 bank 인스턴스를 사용하는 테스트 함수의 첫 번째 줄에 추가한다.

```go
func TestWhatIsTheConversionRateFromEURToUSD(t *testing.T) {
	initExchangeRates()

	tenEuros := s.NewMoney(10, "EUR")
	actualConvertedMoney, err := bank.Convert(tenEuros, "USD")
	assertNil(t, err)
	assertEqual(t, s.NewMoney(12, "USD"), *actualConvertedMoney)
}
```

이제 테스트를 실행하면 성공한다.

```bash
$ go test -v .
...
PASS
ok      tdd     0.002s
```

그러면 이제 TestWhatIsTheConversionRateFromEURToUSD를 유지하는 것이 맞을까? 이 테스트는 어떠한 새로운 피처도 테스트하지 않고 그저 기존 테스트의 취약점을 드러내는 역할만을 한다. 더 중요한 것은 테스트의 순서를 임의로 변경하면서 테스트 간의 의도치 않은 의존성을 발견하는 것이다.

이를 위해 Go는 `-shuffle` 플래그를 제공한다. 이 플래그를 사용하면 테스트를 실행할 때마다 테스트 함수의 순서가 무작위로 변경된다. 이를 통해 테스트 간의 의존성을 확인할 수 있다.

```bash
$ go test -v -shuffle on .
-test.shuffle 1722826536991840026
...
PASS
ok      tdd     0.002s
```

테스트의 순서를 변경하여 우연한 결합을 발견하는 더 나은 방법을 발견했으므로 이제 TestWhatIsTheConversionRateFromEURToUSD 테스트를 제거하자.

매 테스트 함수가 실행되기 전에 bank 인스턴스를 초기화하는 initExchangeRates 함수를 호출하는 방식 외에도, 테이블 기반 테스트를 사용하여 여러 환율을 테스트할 수도 있다. 이 내용은 다루지 않는다.

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "test: verified behavior when modifying an existing exchange rate"
```

## 중간 점검

- 기존 피처를 문서화하기 위해 테스트를 추가했고 테스트 의존성을 배웠다.
- 의도되지 않은 사이드 이펙트를 발견하기 위해 테스트 순서를 변경하는 방법을 배웠다. (Go의 `-shuffle` 플래그)
- 목표한 피처 목록의 모든 피처를 구현했다.

<br>

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
- [x] 환율 수정 허용
