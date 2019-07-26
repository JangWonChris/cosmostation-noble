# cosmostation-chain-exporter-es

## 1.준비사항

1. go mod로 모듈관리를 하기 때문에 go version v1.11 이상이 필요. <br/>
   (현재 backend 서버에 go version 1.12.1으로 설치되어 있음)

2. 프로젝트 빌드를 위해서는 config.yaml 파일이 필요. <br/>
   프로젝트 루트위치에서 아래 명령어를 통해 config.yaml을 다운로드한다. <br/>
   ```
   curl https://gist.githubusercontent.com/wannabit-mina/761efa3482cb820ca9bfcb001bd480bb/raw/ac0ec1e59ca9d6e3389bfacad67682c5e14cdd9c/config.yaml > config.yaml
   ```

## 2.로컬에서 빌드 방법

```
1. git clone https://github.com/cosmostation/cosmostation-cosmos.git
2. cd cosmostation-cosmos/chain-exporter-es
3. go build
4. ./chain-exporter-es server --env=dev --network=cosmos
    ( * env, network 플래그는 필수 / env는 dev/prod, network는 cosmos/kava)
```

## 3. 서버에서 빌드 및 구동

```
1. cd $HOME/chain-exporter-es/chain-exporter-es
2. go build -o $HOME/go/bin/chain-exporter-es
3. sudo systemctl start es-crawler.service
```

<br/>

### Cosmos - es-crawler.service

`위치 : /etc/systemd/system/es-crawler.service`

```
[Unit]
Description=ES Crawler Service
Requires=network-online.target
After=network-online.target

[Service]
EnvironmentFile=/etc/es-crawler.conf
Type=simple
ExecStart=/home/ubuntu/go/bin/chain-exporter-es server --env=${ENV_PROD} --network=${NETWORK_COSMOS}
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

<br/>

### Kava - es-cralwer.service

`위치 : /etc/systemd/system/es-cralwer.service`

```
[Unit]
Description=ES Crawler Service
Requires=network-online.target
After=network-online.target

[Service]
EnvironmentFile=/etc/es-crawler.conf
Type=simple
ExecStart=/home/ubuntu/go/bin/chain-exporter-es server --env=${ENV_DEV} --network=${NETWORK_KAVA}
Restart=on-failure
RestartSec=10s


[Install]
WantedBy=multi-user.target
```

_EnvironmentFile_

```
ENV_DEV = dev
ENV_PROD = prod
NETWORK_COSMOS = cosmos
NETWORK_KAVA = kava
```
