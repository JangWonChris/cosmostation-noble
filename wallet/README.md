## 프로젝트 설명 (Overview)

모바일 및 웹 지갑과 Cosmostaion의 API를 정의 합니다.


## 프로젝트 주요기능 (Feature)

* Version Check : 각 지갑 OS별(Android, IOS, Web) 최신 버전을 확인하는 기능.

* Push Token Update: 유저 address를 기반으로한 Push Token 등록/상태변경


## API Spec

* Version Check : simple get method for each client OS


* Push Token Update

| name         | type         | description                                                | example        |
| ------------ | ------------ | ---------------------------------------------------------- | -------------- |
| address      | String       | user request address fo push notification                  | cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7 |
| chain        | int          | address with associcated chain                             | 0 - cosmos,  1 - iris, 2 - kava |
| os           | String       | user requset device type                                   | Android, IOS  |
| enable       | boolean      | enable or disable push                                     | true, false |
| language     | String       | user setted language                                       | kr, en  |
| token        | String       | device token for notification                              | token string |


## GoRush set

###  고러쉬 레포 [GoRush](https://github.com/appleboy/gorush)

- 설치한 후 웹 서버 실행
- 옵션없이 돌린 후 /api/config를 통해 config 파일을 다운 받은후 나한테현재 상황에 맞게끔 고치면 편함.
- port번호와 android용 apikey만 바꿔주고 다시 -c 옵션으로 실행
- apikey = AIzaSyBGu7jJkZo5cIlINvWp8hDashE3_Oq9VXE  (대외비)

### Post 예제

 - post end point  http://localhost:8088/api/push 
 - payload body
```
{
  "notifications": [
    {
      "tokens": ["dMqkRx5wc2U:APA91bHnuigvPXh7kb8qzROKviED9K-iwucG1HYYmi7GImhs0fkuJw72lJYcAvq-uSYAlQ5WRT3-a8Wofw1Rv4MoSF-Shm2-T-SZJ0pdB_RC5vEyVeQZ2IHQEe5PkfueyH5fRnG3Gp17", "dIhWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuJtAajGjHuqwNGxtOH0LL1sRi3l4ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn"],
      "platform": 2,
      "message": "you received atom with txid 06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
      "title": "Received 11.434532Atom",
      "data": {"notifyto" : "cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7","txid" : "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827", ,"type" : "send"}
    }
  ]
}
```
