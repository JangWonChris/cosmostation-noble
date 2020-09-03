# Database Query Collection

## Account Table

```shell
# 특정 계정 조회 
SELECT * FROM account WHERE account_address = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee';

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
SELECT * FROM "validator" WHERE proposer = '099E2B09583331AFDE35E5FA96673D2CA7DEA316';
SELECT * FROM "validator" WHERE moniker = 'Cosmostation';
SELECT * FROM "validator" WHERE operator_address = 'cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn';

```

## Power Event History Table

```shell
# 특정 검증인 조회 (id_validator = 0)
SELECT * FROM "power_event_history" WHERE id_validator = 0;

# 현재까지 발생한 파워이벤트 트랜잭션 조회
SELECT COUNT(id) as total_count FROM "power_event_history" WHERE power_event_history.proposer = '099E2B09583331AFDE35E5FA96673D2CA7DEA316';
```

## Proposal Table

```shell
```


## Transaction Table 

```shell
-- 트랜잭션 해쉬로 트랜잭션 조회
SELECT * FROM "transaction" WHERE tx_hash = 'D0E335C3E8DCB0B06090DE1FB209AAC24C99FF43B1609E0BCE034300FE891AF1';

-- 메시지 타입별 트랜잭션 조회
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgSend' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgMultiSend' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgDelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgUndelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgBeginRedelegate' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgCreateValidator' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgEditValidator' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgSubmitProposal' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgVote' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgDeposit' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission' ORDER BY id DESC LIMIT 10;
SELECT * FROM "transaction" WHERE messages->0->>'type' = 'pricefeed/MsgPostPrice' LIMIT 10;

-- 메시지 타입 종류별 트랜잭션 조회
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%cosmos-sdk/%' ORDER BY id DESC LIMIT 50;
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%bep3/%' ORDER BY id DESC;
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%cdp/%' ORDER BY id DESC;
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%pricefeed/%' ORDER BY id DESC;
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%kava/%' ORDER BY id DESC;
SELECT * FROM "transaction" WHERE messages->0->>'type' LIKE '%incentive/%' ORDER BY id DESC;

-- 특정 메모가 포함된 트랜잭션 조회 
SELECT * FROM "transaction" WHERE TIMESTAMP > '2020-07-17T14:20:12Z' AND memo = 'Claim Atomic Swap via Cosmostation Wallet' ORDER BY id DESC;

-- 블록높이 사이에서 발생한 트랜잭션 조회
SELECT * FROM "transaction" WHERE height BETWEEN 605761 AND 605763 ORDER BY id DESC;

-- 특정 날짜 이후 발생한 언본딩 트랜잭션 조회
SELECT * FROM "transaction" WHERE (messages->0->>'type' = 'cosmos-sdk/MsgUndelegate') AND time > '2020.01.30' ORDER BY id DESC LIMIT 30;

-- 특정 계정에서 발생한 트랜잭션 조회 (검증인 주소 포함)
-- MsgSend: from_address, to_address
-- MultiSend: inputs, outputs
-- MsgDelegate, MsgWithdrawnDelegation, MsgUndelegate, MsgBeginRedelegate, MsgCreateValidator: delegator_address
-- MsgEditValidator: address
-- MsgSubmitProposal: proposer
-- MsgDeposit: depositor
-- MsgVote voter
SELECT * FROM "transaction" 
    WHERE messages->0->'value'->>'from_address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'to_address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->'inputs'->0->>'address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->'outputs'->0->>'address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'delegator_address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'address' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'proposer' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'depositor' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'voter' = 'cosmos1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqda85ee' OR
    messages->0->'value'->>'validator_address' = 'cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn' AND messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission'
    ORDER BY id DESC LIMIT 50;
```

## INDEXES

```shell
-- 트랜잭션 메시지 value에 GIN indexing 생성 방법
CREATE INDEX transaction_messages_value_idx ON transaction USING BTREE((messages->0->'value'));

-- 트랜잭션 메시지 symbol에 indexing 생성 방법 
CREATE INDEX transaction_messages_symbol_idx ON transaction USING BTREE((messages->0->'value'->>'symbol'));

-- 계정별 트랜잭션을 위한 인덱스
CREATE INDEX CONCURRENTLY transaction_messages_sender_idx ON transaction USING BTREE((messages->0->'value'->>'sender'));
CREATE INDEX CONCURRENTLY transaction_messages_intputs_address_idx ON transaction USING BTREE((messages->0->'value'->'inputs'->0->>'address'));
CREATE INDEX CONCURRENTLY transaction_messages_outputs_address_idx ON transaction USING BTREE((messages->0->'value'->'outputs'->0->>'address'));
CREATE INDEX CONCURRENTLY transaction_messages_voter_idx ON transaction USING BTREE((messages->0->'value'->>'voter'));
CREATE INDEX CONCURRENTLY transaction_messages_proposer_idx ON transaction USING BTREE((messages->0->'value'->>'proposer'));
CREATE INDEX CONCURRENTLY transaction_messages_depositor_idx ON transaction USING BTREE((messages->0->'value'->>'depositor'));
CREATE INDEX CONCURRENTLY transaction_messages_from_idx ON transaction USING BTREE((messages->0->'value'->>'from'));
CREATE INDEX CONCURRENTLY transaction_messages_to_idx ON transaction USING BTREE((messages->0->'value'->>'to'));
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
DELETE FROM "transaction" WHERE height >= 31079;
DELETE FROM "power_event_history" WHERE height >= 31079;
DELETE FROM "deposit" WHERE height >= 31079;
DELETE FROM "vote" WHERE height >= 31079;
```