# cosmos-monitor

The program used to cosmos-monitor the status of specific validators on the Cosmos Hub. 
The program will send an alarm by email in the following cases:
- If the validator is jailed, `receiver1` and `receiver2` will be notified by email;
- If the validator is inactive, `receiver1` and `receiver2` will be notified by email;
- If the validator has not signed five blocks in a row or the unsigned rate of the last `n` blocks reaches `m`, `receiver1` and `receiver2` will be notified by email;
- If there is a new proposal, the `receiver1` will be notified by email

> receiver1: receive all alert emails;
> 
> receiver2: only receive emergency alert emails, like jailed/inactive, etc.
## 1、Roadmap
- Alarm channels will be added, like Telegram, Discord;
- Monitoring of stake income will be added;
- Program will support other chains developed using the Cosmos SDK;

## 2. Hardware requirements
We recommennd the following hardware resources:
- CPU: 4 cores
- Memory: 8GM RAM
- Database: PostgreSQL

## 3. configuration file
- `vim /data/config/conf.yaml`
- Configuration example
```yaml
# database config
# database config
postgres:
  cosmosUser: "cosmos_database_user"
  cosmosPassword: "cosmos_database_password"
  cosmosName: "cosmos_database_name"
  cosmosHost: "cosmos_database_host_ip"
  cosmosPort: "cosmos_database_host_port"

  injectiveUser: "injective_database_user"
  injectivePassword: "injective_database_password"
  injectiveName: "injective_database_name"
  injectiveHost: "injective_database_host_ip"
  injectivePort: "injective_database_host_port"

# gRPC
gRpc:
  cosmosIp: "cosmos_grpc_host_ip"
  cosmosPort: "cosmos_grpc_host_grpc_port"
  injectiveIp: "injective_grpc_host_ip"
  injectivePort: "injective_grpc_host_grpc_port"

alert:
  cosmosOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  injectiveOperatorAddr: "injvaloperxxxxx,injvaloperyyyy"    # The address starts with "injvaloper"
  timeInterval: 600  # Monitoring time interval. 600 means the cosmos-monitor runs every 600 seconds
  blockInterval: 100 # Used to calculate the recent signature rate. 100 means to count the signatures rate of the last 100 blocks
  proportion: 0.05 # Signature rate
  cosmosStartingBlockHeight: 123456 # The cosmos starting block height of the monitor program
  injectiveStartingBlockHeight: 1234567 # The injective starting block height of the monitor program
mail:
  host: "mail_host"
  port: "mail_port"
  username: "mail_username"
  password: "mail_password"
  sender: "sender_mail_address"
  cosmosReceiver1: "cosmos_receiver1_mail_address" # Receive all alert emails of cosmos
  cosmosReceiver2: "cosmos_receiver2_mail_address" # Only receive emergency alert emails of cosmos, like jailed/inactive, etc.
  injectiveReceiver1: "injective_receiver1_mail_address" # Receive all alert emails of injective
  injectiveReceiver2: "injective_receiver2_mail_address" # Only receive emergency alert emails of injective, like jailed/inactive, etc.

log:
  path: "log_save_location"
  level: "log_level" # info, error
  eventlogpath: "event_log_save_location"
```
## 4. Deployment monitoring
- 4.1 Install go
> TIP
> 
> go 1.18++ is required for building and installing the monitor software.

[Golang Official website](https://go.dev/doc/install)

- 4.2 Build monitor binary file
```shell
git clone https://github.com/HashQuark-Research1/cosmos-monitor.git
go build -o cosmos-monitor
```

- 4.3 Edit configuration file
```shell
vim config/conf.yaml
```

- 4.4 Automating your monitor with systemd
```bash
vim /etc/systemd/system/cosmos-monitor.servic
```
```shell
[Unit]
Description=Cosmos Monitor Daemon
After=network.target
[Service]
User=ubuntu
ExecStart=/home/ubuntu/cosmosmonitor/cosmos-monitor -c /data/config/conf.yaml
KillSignal=SIGINT
Restart=on-failure
RestartSec=5s
LimitNOFILE=1000000
[Install]
WantedBy=multi-user.target
[Manager]
DefaultLimitNOFILE=1000000
```

```shell script
systemctl start cosmos-monitor.servic
```