## branch 정보

cosmoshub-4 이후로, 특정 블록 높이에서 업그레이드가 진행되어 아래와 같은 브런치 네이밍 룰을 따름

branch 이름 규칙 : `<chain-id>-<app version>-<block height>`

| chain-id | branch-name | mintscan-backend-library | mintscan-database | note |
|---|---|---|---|---|
| cosmoshub-1 | cosmoshub-1 | sdk-0.33.x-1 | master | |
| cosmoshub-2 | cosmoshub-2 | sdk-0.34.x-1 | master | |
| cosmoshub-3 | cosmoshub-3 | sdk-0.37.x-1 | master | |
| cosmoshub-4 | cosmoshub-4 | sdk-0.42.x | master | 필요시 새로 작업 필요(v5로 덮어써짐)|
| cosmoshub-4 | cosmoshub-4-v5-6910000 | sdk-0.42.x | master | |
| cosmoshub-4 | cosmoshub-4-v6-8695000 | sdk-0.44.x | master | |


## 작업 시 

현재 동작 중인 버전은 
- 해당 브런치에 업데이트 후 master에 머지

과거 버전은 
- 해당 브런치에 업데이트

이름/기능 으로 branch를 만들어 작업 후 해당 브런치에 머지한다. (예, jeonghwan/featureX)

## 참고 사항


### chain 정보

| chain-id | initial-height | num. of blocks | num. of txs | app version |
|---|---|---|---|---|
| cosmoshub-1 | 1 | 500,000 | 20,228 | |
| cosmoshub-2 | 1 | 2,902,000 | 662,621 | |
| cosmoshub-3 | 1 | 5,200,790 | 2,420,082 | gaia v2 |
| cosmoshub-4 |5,200,791 | | | gaia v4 |
| cosmoshub-4 |6,910,000 | | | gaia v5 |
| cosmoshub-4 |8,695,000 | | | gaia v6(vega) |
