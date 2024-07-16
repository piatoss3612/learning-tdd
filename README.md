# 1장 돈 문제

## 레드-그린-리팩터: 테스트 주도 개발 구성 요소

1. 레드: 실패하는 테스트를 작성한다(컴파일 실패 포함). 테스트 스위트(suite)를 실행해서 테스트가 실패하는 것을 확인한다.
2. 그린: 테스트를 통과할 만큼의 최소한의 코드를 작성한다. 테스트 스위트를 실행해서 테스트가 성공하는 것을 확인한다.
3. 리팩터: 중복 코드, 하드 코딩된 값, 프로그래밍 이디엄(idiom)의 부적절한 사용 등을 제거한다. 이 과정에서 테스트가 깨진다면, 깨진 모든 테스트를 그린으로 만드는 것을 우선시한다.

## 문제 인식

- 여러 통화로 돈을 관리하거나 주식 포트폴리오를 관리하는 스프레드시트를 만들어야 한다고 가정하자.

### 요구사항

1. 단일 통화로 된 숫자상에서 간단한 산술 연산이 가능해야 한다.
   - [ ] 5달러 \* 2 = 10달러
   - [ ] 10유로 \* 2 = 20유로
   - [ ] 4002원 / 4 = 1000.5원
2. 통화 간 환전을 지원해야 한다.
   - [ ] 5달러 + 10유로 = 17달러
   - [ ] 1달러 + 1100원 = 2200원

- 각 항목은 테스트 주도 개발로 구현할 하나의 피처(feature)가 된다.

### 첫 번째 실패하는 테스트

- 단일 통화의 곱셈을 테스트한다.

```go
package main

import "testing"

func TestMultiplication(t *testing.T) {
	fiver := Dollar{
		amount: 5,
	}
	tenner := fiver.times(2)
	if tenner.amount != 10 {
		t.Errorf("Expected 10, got %d", tenner.amount)
	}
}
```

- fiver 변수에 amount 필드가 5인 Dollar 구조체를 초기화하여 '5달러'를 나타내는 엔티티를 선언한다.
- fiver.times(2)를 호출하여 5달러를 2배로 만들어 '10달러'를 나타내는 엔티티를 얻는다.
- tenner의 amount 필드가 10인지 확인한다.

> Dollar 구조체와 times 메서드는 아직 존재하지 않는다.

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:6:11: undefined: Dollar
FAIL    tdd [build failed]
FAIL
```

- Dollar가 정의되어 있지 않다는 오류와 함께 테스트가 실패한다.

### 그린으로 전환

1. Dollar의 추상(abstraction)을 만든다.

```go
type Dollar struct {
}
```

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:9:3: unknown field amount in struct literal of type Dollar
./money_test.go:11:18: fiver.times undefined (type Dollar has no field or method times)
FAIL    tdd [build failed]
FAIL
```

- amount 필드와 times 메서드가 없다는 오류와 함께 테스트가 실패한다.

2. Dollar에 amount 필드를 추가한다.

```go
type Dollar struct {
	amount int
}
```

```bash
$ go test -v .
# tdd [tdd.test]
./money_test.go:13:18: fiver.times undefined (type Dollar has no field or method times)
FAIL    tdd [build failed]
FAIL
```

- times 메서드가 없다는 오류와 함께 테스트가 실패한다.

3. Dollar에 times 메서드를 추가한다.

```go
func (d Dollar) times(multiplier int) Dollar {
	return Dollar{
		amount: 10,
	}
}
```

- multiplier와 amount를 곱한 값을 반환하는 것이 올바른 산술 연산이지만,
- 우선은 테스트 예상 결과를 반환하는 가장 간단한 코드를 작성한다.
- 하드 코딩된 값인 10을 반환한다.

```bash
$ go test -v .
=== RUN   TestMultiplication
--- PASS: TestMultiplication (0.00s)
PASS
ok      tdd     0.002s
```

- 테스트가 성공한다.

### 마무리하기

- 이제 리팩터링 단계로 넘어가야 한다.
- 리팩터링은 중복 코드, 하드 코딩된 값, 프로그래밍 이디엄의 부적절한 사용 등을 제거하는 것이다.

### 이상한 점 찾기

1. 결합: '5달러를 2배하면 10달러'를 검증하는 코드를 작성했지만, 이를 '10달러를 2배하면 20달러'로 변경하면 테스트가 실패한다. 테스트를 변경하면 코드도 변경해야 하는 의존성(논리적 결합)이 존재한다는 것을 알 수 있다.
2. 중복: 테스트와 코드에 10이라는 값이 중복되어 있다. 10이라는 값은 실제로 5와 2를 곱한 결과이므로, 이 값을 계산하는 코드를 작성하면 중복을 제거할 수 있다.

```go
func (d Dollar) times(multiplier int) Dollar {
	return Dollar{
		amount: d.amount * multiplier,
	}
}
```

- 테스트를 실행하면 성공한다.

### 변경 사항 반영하기

- 첫 번째 피처를 구현했으므로, 코드를 버전 관리 시스템에 반영해 준다.

```bash
$ git add .
$ git commit -m "feat: first green test"
```

- 커밋 메시지는 시맨틱 커밋 메시지 규칙을 따른다.

- 구현한 피처는 체크리스트에 체크한다.

  - [x] 5달러 \* 2 = 10달러
  - [ ] 10유로 \* 2 = 20유로
  - [ ] 4002원 / 4 = 1000.5원
  - [ ] 5달러 + 10유로 = 17달러
  - [ ] 1달러 + 1100원 = 2200원
