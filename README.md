# 프로젝트 설명 (Overview)

Cosmos Hub 관련 프로젝트를 다루는 Repository 입니다.

# 프로젝트 구조 (Structure)

## 1. api 

> Mintscan 익스플로러와 지갑에 필요한 API

## 2. chain-exporter

> PostgreSQL 데이터베이스에 블록 데이터 저장하는 프로젝트

## 3. chain-exporter-es 

> AWS Elasticsearch에 트랜잭션만 저장하는 프로젝트

## 4. stats-exporter

> PostgreSQL 데이터베이스에 통계 데이터를 저장하는 프로젝트

# 이슈 또는 기능 제안 (Bug Report or Feature Request)

버그 및 에러 또는 기능 제안은 Issues에 남겨주세요

# 테스트넷 

## gaia-13006 테스트넷에서 사용중인 주소 (multiple tokens)

```
[Address]

cosmos1x5wgh6vwye60wv3dtshs9dmqggwfx2ldnqvev0

[Mnemonic phrases]

```

# 브랜치 / 배포 / 태그

각 작업자는 `이름(name)/브랜치이름(branch) 또는 본인이 원하는 브랜치 이름`을 생성하여 개발을 진행하면 된다.

배포(release)는 `master` 브랜치에서만 병합을 할 것이며 병합 할 때 `release`를 할 예정이다.

태그(tag)는 본인이 어느정도 작업을 마친 후 기록하고 싶은 부분에 하면 된다. 아래와 같은 규칙대로 `tag` 하기로 했다.

```
v0.1.0-rc1
v0.1.0-rc2
v0.1.0-rc3
```
