# 3장 포트폴리오

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.5원
- [ ] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원

이번에는 통화를 혼합한 덧셈을 구현해보자.

## 다음 테스트 설계하기

- [ ] 5달러 + 10유로 = 17달러

다음 피처를 개발하기에 앞서, 먼저 밑그림을 그려보는 것이 좋다. 여기서 1유로가 1.2달러로 교환된다고 가정한다.그러면 다음과 같은 것도 고려해야 한다.

- 1유로 + 1유로 = 2.4달러
- 1유로 + 1유로 = 2유로

두 개의 Money 엔티티를 더한 결과는, 연관된 모든 통화간 환율을 알 수 있다면 어떤 통화로도 표현될 수 있다. '어떤 통화로도 표현'이라는 말은 도메인 모델이 확장되어야 한다는 것을 의미한다. 이처럼 테스트를 설계할 때는 도메인 모델이 어떻게 확장될지 고려해야 한다.

확장된 도메인 모델은 Portfolio라는 새로운 개념을 도입할 수 있다. Portfolio는 여러 통화를 포함할 수 있으며, 각 통화는 환율을 가지고 있다. 이제 테스트를 작성해보자.

```go
func TestAddition(t *testing.T) {
	var portfolio Portfolio
	var portfolioInDollars Money

	fiveDollars := Money{amount: 5, currency: "USD"}
	tenDollars := Money{amount: 10, currency: "USD"}
	fifteenDollars := Money{amount: 15, currency: "USD"}

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(tenDollars)
	portfolioInDollars = portfolio.Evaluate("USD")

	assertEqual(t, fifteenDollars, portfolioInDollars)
}
```

포트폴리오에 5달러와 10달러를 추가하고, 포트폴리오에 들어있는 돈을 달러로 평가하면 15달러가 나와야 한다. 이제 테스트를 실행해보자.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:64:16: undefined: Portfolio
FAIL    tdd [build failed]
FAIL
```

당연하게도 Portfolio가 정의되지 않았다는 에러가 발생한다. Portfolio를 정의해보자.

```go
type Portfolio []Money

func (p Portfolio) Add(money Money) Portfolio {
	return p
}

func (p Portfolio) Evaluate(currency string) Money {
	return Money{amount: 15, currency: "USD"}
}
```

Portfolio를 Money 타입의 슬라이스로 정의하고, Add와 Evaluate 메서드를 정의했다. 각 메서드는 테스트를 통과하기 위한 최소한의 구현을 가지고 있다. 하드코딩된 값을 반환하도록 구현했기 때문에 테스트를 통과할 것이다. 이제 테스트를 실행해보자.

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
ok      tdd     (cached)
```

테스트가 그린 상태가 되었다. 이제 테스트를 통과하기 위해 구현한 코드를 리팩토링해보자. 테스트와 프로덕션 코드에서 중복된 코드를 찾아내고, 중복을 제거하면서 테스트를 통과하는 코드를 작성해보자.

우선은 테스트와 Portfolio의 Evaluate 메서드에서 15라는 값이 중복되어 나타나는 것을 발견했다. Evaluate 메서드에서 실제로 계산된 값을 반환하도록 수정해보자.

```go
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
...
=== RUN   TestAddition
    money_test.go:59: Expected {15 USD}, got {0 USD}
--- FAIL: TestAddition (0.00s)
```

테스트가 실패했다. 이번에는 15라는 값이 아닌 0이라는 값이 반환되었다. 이는 빈 슬라이스를 순회하면서 total 변수에 더해지는 값이 없기 때문이다. 이번에는 Add 메서드가 실제로 Money를 추가하도록 수정해보자.

```go
func (p Portfolio) Add(money Money) Portfolio {
	return append(p, money)
}
```

이제 테스트를 실행해보자.

```bash
$ go test -v .
...
=== RUN   TestAddition
--- PASS: TestAddition (0.00s)
```

테스트가 그린 상태가 되었다. 그러나 테스트로는 발견되지 않은 문제가 있다. Evaluate 메서드는 통화를 고려하지 않고 단순히 금액을 더하고 있다.

코드에서 이런 '미련한' 동작을 지우는 방법을 테스트해야 할까 아니면 '리팩터링'을 해야 할까? 여기에 만능인 답은 없다. 테스트 주도 개발은 스스로 속도를 조절할 수 있다. 여기서는 아직 환율의 개념이 정의되지 않았으므로, 우선은 속도를 맞추기 위해 '미련한' 동작의 수정을 나중으로 미루도록 하자.

## 변경 사항 반영하기

```bash
$ git add .
$ git commit -m "feat: addition feature for Moneys in the same currency done"
```

## 중간 점검

1. 다른 통화를 더하는 피처를 구현하기 위해 새로운 도메인 모델인 Portfolio를 도입했다.
2. 한 번에 모두 해결하기에는 양이 많으므로, 먼저 동일한 통화를 가진 두 Money 인스턴스를 더하는 피처를 구현했다.
3. 개선해야 할 부분이 있지만, 일단은 미뤄두고 새로운 개념을 도입하는 과정에서 함께 살펴볼 수 있도록 했다.
4. 테스트와 프로덕션 코드가 늘어남에 따라 하나의 파일에 모든 코드를 작성하는 것이 불편해졌다. 이제 코드를 분리해보자.

- [x] 5달러 \* 2 = 10달러
- [x] 10유로 \* 2 = 20유로
- [x] 4002원 / 4 = 1000.
- [ ] 5달러 + 10달러 = 15달러
- [ ] 5달러 + 10유로 = 17달러
- [ ] 1달러 + 1100원 = 2200원
