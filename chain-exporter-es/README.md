# cosmostation-chain-exporter-es

## 1.준비사항

1. go mod로 모듈관리를 하기 때문에 go version v1.11 이상이 필요. <br/>
(현재 backend 서버에 go version 1.12.1으로 설치되어 있음)

2. 프로젝트 빌드를 위해서는 config.yaml 파일이 필요. <br/>
config.yaml은 따로 공유하는 것으로 한다.

## 2.로컬에서 빌드 방법
```
1. git clone https://github.com/cosmostation/cosmostation-cosmos.git
2. cd cosmostation-cosmos/chain-exporter-es
3. go build
4. ./chain-exporter-es server --env=dev --network=cosmos
    ( * env, network 플래그는 필수 / env는 dev/prod, network는 cosmos/kava)
```

## [참고]escrawler.service

위치 : /lib/systemd/system/escrawler.service

파일내용

```
[Unit]
Description=ES Crawler Service
Requires=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/home/ubuntu/elasticsearch/cosmostation-chain-exporter-es/cosmostation-chain-exporter-es
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```
