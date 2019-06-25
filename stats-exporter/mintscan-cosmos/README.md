# cosmostation-stats-exporter

Mintscan Explorer 에서 보여주는 통계 데이터들을 데이터베이스에 저장하는 프로그램


## 1.준비사항

go mod로 모듈관리를 하기 때문에 go version v1.11 이상이 필요하다.

## 2.빌드 방법
```
1. 프로젝트 git clone
2. Build 빌드하여 Binary 파일 생성
3. 서비스 등록 후 daemon-reload 후 서비스 시작

$ sudo systemctl daemon-reload
$ sudo systemctl restart stats-exporter.service
```


## [참고] stats-exporter.service
위치 : /lib/systemd/system/stats-exporter.service

파일내용
```
[Unit]
Description=ES Crawler Service
Requires=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/home/ubuntu/cosmostation-stats-exporter/cosmostation-stats-exporter
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

## 자동 배포 스크립트

```
#!/bin/bash

# /lib/systemd/system/stats-exporter.service

# Stop the daemon service
echo ""
echo "Stop stats-exporter service"
sudo service stats-exporter stop

# Git pull the latest commit
echo ""
echo "Pull the latest commit"
cd cosmostation-stats-exporter && git pull origin master

echo ""
echo "Go build in progress..."
go build

echo ""
echo "Start service..."
sudo service stats-exporter start

echo ""
echo "Start tailing logs..."
tail -f /var/log/syslog
```