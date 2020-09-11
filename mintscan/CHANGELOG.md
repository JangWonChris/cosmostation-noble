# Changelog

## Unreleased - 2020.08.31

### 추가된 기능 (Features)

### 개선된 기능 (Improvements)

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
