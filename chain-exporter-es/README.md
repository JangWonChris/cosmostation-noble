# cosmostation-chain-exporter-es

## 1.준비사항

go mod로 모듈관리를 하기 때문에 go version v1.11 이상이 필요하다.
(현재 backend 서버에 go version 1.12.1으로 설치되어 있음)

## 2.빌드 방법

```
1. https://github.com/cosmostation-cosmos/chain-exporter-es.git
2. cd cosmostation-chain-exporter-es
3. sudo nano main.go
    1. func run() 에서
      esconfig = config.ElasticConfig{}.GetProdConfig()
      //esconfig = config.ElasticConfig{}.GetAppConfig()
      //esconfig = config.ElasticConfig{}.GetDevConfig()

     배포 환경에 맞는 config 파일 셋팅해주기 (일단은 주석처리로...프로세스 보완필요)

3. go build
4. sudo systemctl daemon-reload
5. sudo systemctl restart escrawler.service
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
