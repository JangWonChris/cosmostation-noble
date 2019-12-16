## 프로젝트 설명 (Overview)

모바일 및 웹 지갑과 Cosmostaion의 API를 정의 합니다.


## 프로젝트 주요기능 (Feature)

* Version Check : 각 지갑 OS별(Android, IOS, Web) 최신 버전을 확인하는 기능.

* Push Token Update: 유저 address를 기반으로한 Push Token 등록/상태변경


## API Spec

* Version Check : simple get method for each client OS


* Push Token Update

| name  | type | description | example |
| :------------ | ------------: | :----------- | :----------- |
| address  | String  | user request address fo push notification | cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7 |
| chain  | int  | address with associcated chain | 0 - cosmos,  1 - iris, 2 - kava|
| os  | String  | user requset device type | Android, IOS |
| enable  | boolean  | enable or disable push | true, false |
| token  | String  | device token for notification | dIhWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuJtAajGjHuqwNGxtOH0LL1sRi3l4ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn |


## Notification Payload Spec