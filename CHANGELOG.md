# Changelog

## 요구 사항

- genesis.json parsing
	- accounts coins_total을 포함한 coin컬럼에 multi denom에 대응가능한 구조로 변경해야 한다.

	- 현재 우리는, 자산의 변동을 기록하지 않고, 최신 결과를 보여주기 때문에, denom종류별 합계를 구할 수 있는 형태로 변경이 필요하다.

	- genesis.json 파일을 이용해 genesis account를 구할 때, account type 별 테스트가 필요하다.

	- genesis.json 파일을 이용하기 때문에, node home 및 동적 경로 입력이 가능해야 한다.

	- validator set, module account 등 genesis에서 추출하여 저장할 로직을 구현하여 저장한다.
	
	- 최종적으로 기존 로직과 분리하여 flag 옵션을 통해 실행하고 종료하도록 한다.

- transaction 테이블에는 기존 테이블과 같은 형태로 데이터를 삽입한다.

- transaction에 저장된 데이터와 별개로 transaction json chunk를 별도의 테이블에 저장하고, 추후 가공 할 때 활용한다. (새로 싱크 방지)



## Unreleased - 2020.08.31

### 추가된 기능 (Features)

### 개선된 기능 (Improvements)

### 수정된 버그 (Bug fixes)

- (Validator) inactive validator rank가 0인 문제 수정


# Changelog 작성 템플릿 (Format)

## 버전 (Version) - 날짜 (Date)

### 추가된 기능 (Features)

- [1]
- [2]

### 개선된 기능 (Improvements)

- [1]
- [2]

### 수정된 버그 (Bug fixes)

- [1]
- [2]

