# Changelog

## Unreleased - 2020.08.31

### 추가된 기능 (Features)

### 개선된 기능 (Improvements)

- model.repond() 함수에서 interface{}를 전달 받았을 때, []byte 형 타입이면, json encode를 하지 않도록 변경
- model.validator.go 에 있던 bonded, unbonded, unbonding status 상수 삭제 : sdktypes.Bonded, sdktypes.Unbonding, sdktypes.Unbonded 변수를 int형으로 type cast하여 사용하도록 변경
- 모든 rest-server 질의를 걷어내고(폐기 예정인 API제외) 모두 GRPC 질의로 변경

### 수정된 버그 (Bug fixes)

- (validator/PowerEventHistory API) chain-exporter 에서 정의된 MsgType을 사용하도록 변경 (EventType -> MsgType)
- (Governonce/Proposal id 조회 API) option의 const string이 "yes" -> "Yes"로 변경
- (Post Txs API) list를 전달 할 때, 변경 된 Tx 구조로 리턴할 수 있게끔 변경

# Changelog 작성 템플릿 (Format)

## 버전 (Version) - 날짜 (Date)

### 추가된 기능 (Features)

- [1] status를 메모리에 담아두고 5초마다 업데이트하도록 변경. 이제, 기존의 GetStatus는 값을 복사해 리턴하도록 변경, SetStatus 추가
- [2]

### 개선된 기능 (Improvements)

- [1]
- [2]

### 수정된 버그 (Bug fixes)

- [1]
- [2]
