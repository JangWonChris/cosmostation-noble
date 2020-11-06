# Database Query Collection

## Account Table

```shell
# 특정 계정 조회 
SELECT * FROM account WHERE account_address = 'kava1m36xddywe0yneykv34az8smzhtxy3nyc6v9jdj';

# 네트워크에 있는 주소 갯수 조회
SELECT id FROM account ORDER BY id DESC LIMIT 1;
SELECT count(*) FROM account;

# 가장 많은 토큰을 보유한 계정 조회

# 언본딩 자산이 있는 계정 조회
```

## Block Table

```shell
# 특정 블록 높이 블록정보 조회
SELECT * FROM block WHERE height = '375804';

# 블록 높이 구간 조회
SELECT * FROM block WHERE height BETWEEN '366313' AND '366317';
```

## Validator Table

```shell
# 특정 검증인 조회
SELECT * FROM "validator" WHERE proposer = '7495072779ED8CD1FC19DB02232BB5FCB862FEA6';
SELECT * FROM "validator" WHERE moniker = 'Cosmostation';
SELECT * FROM "validator" WHERE operator_address = 'kavavaloper12g40q2parn5z9ewh5xpltmayv6y0q3zs6ddmdg';

```

## Power Event History Table

```shell
# 특정 검증인 조회 (id_validator = 0)
SELECT * FROM "power_event_history" WHERE id_validator = 0;

# 현재까지 발생한 파워이벤트 트랜잭션 조회
SELECT COUNT(id) as total_count FROM "power_event_history" WHERE power_event_history.proposer = '5D451E3630313E8C04D3D6292A318DC65160A644';
```

## Proposal Table

```shell
```


## Transaction Table 

```shell
# 트랜잭션 해쉬로 트랜잭션 조회
SELECT * FROM "transaction_legacy" WHERE tx_hash = 'C5D57A1A35CFE4EFB58BE35B220B8FD677BFF4FB839DC4A845FFF8AFE8B660B2';

# 메시지 타입별 트랜잭션 조회
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgSend' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgMultiSend' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgDelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgUndelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgBeginRedelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgCreateValidator' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgEditValidator' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgSubmitProposal' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgVote' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgDeposit' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' = 'pricefeed/MsgPostPrice' LIMIT 10;

# 메시지 타입 종류별 트랜잭션 조회
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%cosmos-sdk/%' ORDER BY id DESC LIMIT 50;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%bep3/%' ORDER BY id DESC;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%cdp/%' ORDER BY id DESC;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%pricefeed/%' ORDER BY id DESC;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%kava/%' ORDER BY id DESC;
SELECT * FROM "transaction_legacy" WHERE messages->0->>'type' LIKE '%incentive/%' ORDER BY id DESC;

# 특정 메모가 포함된 트랜잭션 조회 
SELECT * FROM "transaction_legacy" WHERE TIMESTAMP > '2020-07-17T14:20:12Z' AND memo = 'Claim Atomic Swap via Cosmostation Wallet' ORDER BY id DESC;

# 블록높이 사이에서 발생한 트랜잭션 조회
SELECT * FROM "transaction_legacy" WHERE height BETWEEN 605761 AND 605763 ORDER BY id DESC;

# 특정 날짜 이후 발생한 언본딩 트랜잭션 조회
SELECT * FROM "transaction_legacy" WHERE (messages->0->>'type' = 'cosmos-sdk/MsgUndelegate') AND time > '2020.01.30' ORDER BY id DESC LIMIT 30;

# 특정 계정에서 발생한 트랜잭션 조회 (검증인 주소 포함)
# MsgSend: from_address, to_address
# MultiSend: inputs, outputs
# MsgDelegate, MsgWithdrawnDelegation, MsgUndelegate, MsgBeginRedelegate, MsgCreateValidator: delegator_address
# MsgEditValidator: address
# MsgSubmitProposal: proposer
# MsgDeposit: depositor
# MsgVote voter
SELECT * FROM "transaction_legacy" 
    WHERE messages->0->'value'->>'from_address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'to_address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->'inputs'->0->>'address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->'outputs'->0->>'address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'delegator_address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'address' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'proposer' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'depositor' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'voter' = 'kava140g8fnnl46mlvfhygj3zvjqlku6x0fwuhfj3uf' OR
    messages->0->'value'->>'validator_address' = 'kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7' AND messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission'
    ORDER BY id DESC LIMIT 50;
```

## RESET SEQUENCE ID

```shell
-- Get Max ID from table
SELECT MAX(id) FROM stats_market5m;
SELECT MAX(id) FROM stats_market1h;
SELECT MAX(id) FROM stats_market1d;
SELECT MAX(id) FROM stats_network1h;
SELECT MAX(id) FROM stats_network1d;
SELECT MAX(id) FROM stats_validators1h;
SELECT MAX(id) FROM stats_validators1d;

-- Get Next ID from table
SELECT nextval('stats_market5m_id_seq');
SELECT nextval('stats_market1h_id_seq');
SELECT nextval('stats_market1d_id_seq');
SELECT nextval('stats_network1h_id_seq');
SELECT nextval('stats_network1d_id_seq');
SELECT nextval('stats_validators1h_id_seq');
SELECT nextval('stats_validators1d_id_seq');

-- Set Next ID Value to MAX ID
SELECT setval('stats_market5m_id_seq', (SELECT MAX(id) FROM stats_market5m));
SELECT setval('stats_market1h_id_seq', (SELECT MAX(id) FROM stats_market1h));
SELECT setval('stats_market1d_id_seq', (SELECT MAX(id) FROM stats_market1d));
SELECT setval('stats_network1h_id_seq', (SELECT MAX(id) FROM stats_network1d));
SELECT setval('stats_network1d_id_seq', (SELECT MAX(id) FROM stats_network1d));
SELECT setval('stats_validators1h_id_seq', (SELECT MAX(id) FROM stats_validators1h));
SELECT setval('stats_validators1d_id_seq', (SELECT MAX(id) FROM stats_validators1d));
```

## DELETE

```shell
# 특정 블록높이 이후 모든 데이터 삭제
DELETE FROM "account" WHERE height >= 31079;
DELETE FROM "block" WHERE height >= 31079;
DELETE FROM "evidence" WHERE height >= 31079;
DELETE FROM "miss_detail" WHERE height >= 31079;
DELETE FROM "miss" WHERE start_height >= 31079;
DELETE FROM "transaction_legacy" WHERE height >= 31079;
DELETE FROM "power_event_history" WHERE height >= 31079;
DELETE FROM "deposit" WHERE height >= 31079;
DELETE FROM "vote" WHERE height >= 31079;
```