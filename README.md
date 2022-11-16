# cosmos-monitor

This program is used to monitor specific validators on a blockchain deployed using the Cosmos SDK. Currently supports 16 blockchains including Cosmos Hub
For blockchains that provide consensus services, the program will send alarm emails in the following cases:
- If the validator is jailed, `receiver1` and `receiver2` will be notified by email;
- If the validator is inactive, `receiver1` and `receiver2` will be notified by email;
- If the validator has not signed five blocks in a row or the unsigned rate of the last `n` blocks reaches `m`, `receiver1` and `receiver2` will be notified by email;
- If the validator rank exceeds the set rank threshold, `receiver1` and `receiver2` will be notified by email;
- If there is a new proposal, the `receiver1` will be notified by email

For the consumer chain, the program will send alarm emails in the following cases:
- If the validator has not signed five blocks in a row or the unsigned rate of the last `n` blocks reaches `m`, `receiver1` and `receiver2` will be notified by email;
- If the validator rank exceeds the set rank threshold, `receiver1` and `receiver2` will be notified by email;

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
postgres:
  apolloUser: "apollo_database_user"
  apolloPassword: "apollo_database_password"
  apolloName: "apollo_database_name"
  apolloHost: "apollo_database_host_ip"
  apolloPort: "apollo_database_host_port"

  bandUser: "band_database_user"
  bandPassword: "band_database_password"
  bandName: "band_database_name"
  bandHost: "band_database_host_ip"
  bandPort: "band_database_host_port"

  cosmosUser: "cosmos_database_user"
  cosmosPassword: "cosmos_database_password"
  cosmosName: "cosmos_database_name"
  cosmosHost: "cosmos_database_host_ip"
  cosmosPort: "cosmos_database_host_port"

  evmosUser: "evmos_database_user"
  evmosPassword: "evmos_database_password"
  evmosName: "evmos_database_name"
  evmosHost: "evmos_database_host_ip"
  evmosPort: "evmos_database_host_port"

  injectiveUser: "injective_database_user"
  injectivePassword: "injective_database_password"
  injectiveName: "injective_database_name"
  injectiveHost: "injective_database_host_ip"
  injectivePort: "injective_database_host_port"

  junoUser: "juno_database_user"
  junoPassword: "juno_database_password"
  junoName: "juno_database_name"
  junoHost: "juno_database_host_ip"
  junoPort: "juno_database_host_port"

  neutronUser: "neutron_database_user"
  neutronPassword: "neutron_database_password"
  neutronName: "neutron_database_name"
  neutronHost: "neutron_database_host_ip"
  neutronPort: "neutron_database_host_port"

  nyxUser: "nyx_database_user"
  nyxPassword: "nyx_database_password"
  nyxName: "nyx_database_name"
  nyxHost: "nyx_database_host_ip"
  nyxPort: "nyx_database_host_port"

  persistenceUser: "persistence_database_user"
  persistencePassword: "persistence_database_password"
  persistenceName: "persistence_database_name"
  persistenceHost: "persistence_database_host_ip"
  persistencePort: "persistence_database_host_port"

  providerUser: "provider_database_user"
  providerPassword: "provider_database_password"
  providerName: "provider_database_name"
  providerHost: "provider_database_host_ip"
  providerPort: "provider_database_host_port"

  rizonUser: "rizon_database_user"
  rizonPassword: "rizon_database_password"
  rizonName: "rizon_database_name"
  rizonHost: "rizon_database_host_ip"
  rizonPort: "rizon_database_host_port"

  secretUser: "secret_database_user"
  secretPassword: "secret_database_password"
  secretName: "secret_database_name"
  secretHost: "secret_database_host_ip"
  secretPort: "secret_database_host_port"

  sommelierUser: "sommelier_database_user"
  sommelierPassword: "sommelier_database_password"
  sommelierName: "sommelier_database_name"
  sommelierHost: "sommelier_database_host_ip"
  sommelierPort: "sommelier_database_host_port"

  sputnikUser: "sputnik_database_user"
  sputnikPassword: "sputnik_database_password"
  sputnikName: "sputnik_database_name"
  sputnikHost: "sputnik_database_host_ip"
  sputnikPort: "sputnik_database_host_port"

  teritoriUser: "teritori_database_user"
  teritoriPassword: "teritori_database_password"
  teritoriName: "teritori_database_name"
  teritoriHost: "teritori_database_host_ip"
  teritoriPort: "teritori_database_host_port"

  xplaUser: "xpla_database_user"
  xplaPassword: "xpla_database_password"
  xplaName: "xpla_database_name"
  xplaHost: "xpla_database_host_ip"
  xplaPort: "xpla_database_host_port"

# gRPC
gRpc:
  apolloIp: "apollo_grpc_host_ip"
  apolloPort: "apollo_grpc_host_grpc_port"

  bandIp: "band_grpc_host_ip"
  bandPort: "band_grpc_host_grpc_port"

  cosmosIp: "cosmos_grpc_host_ip"
  cosmosPort: "cosmos_grpc_host_grpc_port"

  evmosIp: "evmos_grpc_host_ip"
  evmosPort: "evmos_grpc_host_grpc_port"

  injectiveIp: "injective_grpc_host_ip"
  injectivePort: "injective_grpc_host_grpc_port"

  junoIp: "juno_grpc_host_ip"
  junoPort: "juno_grpc_host_grpc_port"

  neutronIp: "neutron_grpc_host_ip"
  neutronPort: "neutron_grpc_host_grpc_port"

  nyxIp: "nyx_grpc_host_ip"
  nyxPort: "nyx_grpc_host_grpc_port"

  persistenceIp: "persistence_grpc_host_ip"
  persistencePort: "persistence_grpc_host_grpc_port"

  providerIp: "provider_grpc_host_ip"
  providerPort: "provider_grpc_host_grpc_port"

  rizonIp: "rizon_grpc_host_ip"
  rizonPort: "rizon_grpc_host_grpc_port"

  secretIp: "secret_grpc_host_ip"
  secretPort: "secret_grpc_host_grpc_port"

  sommelierIp: "sommelier_grpc_host_ip"
  sommelierPort: "sommelier_grpc_host_grpc_port"

  sputnikIp: "sputnik_grpc_host_ip"
  sputnikPort: "sputnik_grpc_host_grpc_port"

  teritoriIp: "teritori_grpc_host_ip"
  teritoriPort: "teritori_grpc_host_grpc_port"

  xplaIp: "xpla_grpc_host_ip"
  xplaPort: "xpla_grpc_host_grpc_port"

alert:
  apolloOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  bandOperatorAddr: "bandvaloperxxxx,bandvaloperyyyy" # The address starts with "bandvaloper"
  cosmosOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  evmosOperatorAddr: "evmosvaloperxxxx,evmosvaloperyyyy" # The address starts with "evmosvaloper"
  injectiveOperatorAddr: "injvaloperxxxxx,injvaloperyyyy"    # The address starts with "injvaloper"
  junoOperatorAddr: "junovaloperxxxxx,junovaloperyyyy"    # The address starts with "cosmosvaloper"
  neutronOperatorAddr: "neutronvaloperxxxx,neutronvaloperyyyy" # The address starts with "neutronvaloper"
  nyxOperatorAddr: "nvaloperxxxxx,nvaloperyyyy"    # The address starts with "nvaloper"
  persistenceOperatorAddr: "persistencevaloperxxxx,persistencevaloperyyyy" # The address starts with "persistencevaloper"
  providerOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  rizonOperatorAddr: "rizonvaloperxxxxx,rizonvaloperyyyy"    # The address starts with "rizonvaloper"
  secretOperatorAddr: "secretvaloperxxxxx,secretvaloperyyyy"    # The address starts with "secretvaloper"
  sommelierOperatorAddr: "sommeliervaloperxxxx,sommeliervaloperyyyy" # The address starts with "sommeliervaloper"
  sputnikOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  teritoriOperatorAddr: "teritorivaloperxxxxx,teritorivaloperyyyy"    # The address starts with "teritorivaloper"
  xplaOperatorAddr: "xplavaloperxxxxx,xplavaloperyyyy"    # The address starts with "xplavaloper"

  timeInterval: 600  # Monitoring time interval. 600 means the cosmos-monitor runs every 600 seconds
  blockInterval: 100 # Used to calculate the recent signature rate. 100 means to count the signatures rate of the last 100 blocks
  proportion: 0.05 # Signature rate

  apolloRankingThreshold: 100    # The apollo validator ranking threshold
  bandRankingThreshold: 100  # The band validator ranking threshold
  cosmosRankingThreshold: 100    # The cosmos validator ranking threshold
  evmosRankingThreshold: 100 # The evmos validator ranking threshold
  injectiveRankingThreshold: 100 # The injective validator ranking threshold
  junoRankingThreshold: 100  # The juno validator ranking threshold
  nyxRankingThreshold: 100   # The nyx validator ranking threshold
  neutronRankingThreshold: 100   # The neutron validator ranking threshold
  persistenceRankingThreshold: 100    # The persistence validator ranking threshold
  providerRankingThreshold: 100  # The provider validator ranking threshold
  rizonRankingThreshold: 100 # The rizon validator ranking threshold
  secretRankingThreshold: 100    # The secret validator ranking threshold
  sommelierRankingThreshold: 100 # The sommelier validator ranking threshold
  sputnikRankingThreshold: 100   # The sputnik validator ranking threshold
  teritoriRankingThreshold: 100  # The teritori validator ranking threshold
  xplaRankingThreshold: 100  # The xpla validator ranking threshold

mail:
  host: "mail_host"
  port: "mail_port"
  username: "mail_username"
  password: "mail_password"
  sender: "sender_mail_address"
  apolloReceiver1: "band_receiver1_mail_address" # Receive all alert emails of apollo
  apolloReceiver2: "band_receiver2_mail_address" # Receive all alert emails of apollo
  bandReceiver1: "band_receiver1_mail_address" # Receive all alert emails of band
  bandReceiver2: "band_receiver2_mail_address" # Only receive emergency alert emails of band, like jailed/inactive, etc.
  cosmosReceiver1: "cosmos_receiver1_mail_address" # Receive all alert emails of cosmos
  cosmosReceiver2: "cosmos_receiver2_mail_address" # Only receive emergency alert emails of cosmos, like jailed/inactive, etc.
  evmosReceiver1: "evmos_receiver1_mail_address" # Receive all alert emails of evmos
  evmosReceiver2: "evmos_receiver2_mail_address" # Only receive emergency alert emails of evmos, like jailed/inactive, etc.
  injectiveReceiver1: "injective_receiver1_mail_address" # Receive all alert emails of injective
  injectiveReceiver2: "injective_receiver2_mail_address" # Only receive emergency alert emails of injective, like jailed/inactive, etc.
  junoReceiver1: "juno_receiver1_mail_address" # Receive all alert emails of juno
  junoReceiver2: "juno_receiver2_mail_address" # Only receive emergency alert emails of juno, like jailed/inactive, etc.
  neutronReceiver1: "neutron_receiver1_mail_address" # Receive all alert emails of neutron
  neutronReceiver2: "neutron_receiver2_mail_address" # Only receive emergency alert emails of neutron, like jailed/inactive, etc.
  nyxReceiver1: "nyx_receiver1_mail_address" # Receive all alert emails of nyx
  nyxReceiver2: "nyx_receiver2_mail_address" # Only receive emergency alert emails of nyx, like jailed/inactive, etc.
  persistenceReceiver1: "persistence_receiver1_mail_address" # Receive all alert emails of persistence
  persistenceReceiver2: "persistence_receiver2_mail_address" # Only receive emergency alert emails of persistence, like jailed/inactive, etc.
  providerReceiver1: "provider_receiver1_mail_address" # Receive all alert emails of provider
  providerReceiver2: "provider_receiver2_mail_address" # Only receive emergency alert emails of provider, like jailed/inactive, etc.
  rizonReceiver1: "rizon_receiver1_mail_address" # Receive all alert emails of rizon
  rizonReceiver2: "rizon_receiver2_mail_address" # Only receive emergency alert emails of rizon, like jailed/inactive, etc.
  secretReceiver1: "secret_receiver1_mail_address" # Receive all alert emails of secret
  secretReceiver2: "secret_receiver2_mail_address" # Only receive emergency alert emails of secret, like jailed/inactive, etc.
  sommelierReceiver1: "sommelier_receiver1_mail_address" # Receive all alert emails of sommelier
  sommelierReceiver2: "sommelier_receiver2_mail_address" # Only receive emergency alert emails of sommelier, like jailed/inactive, etc.
  sputnikReceiver1: "sputnik_receiver1_mail_address" # Receive all alert emails of sputnik
  sputnikReceiver2: "sputnik_receiver2_mail_address" # Only receive emergency alert emails of sputnik, like jailed/inactive, etc.
  teritoriReceiver1: "teritori_receiver1_mail_address" # Receive all alert emails of teritori
  teritoriReceiver2: "teritori_receiver2_mail_address" # Only receive emergency alert emails of teritori, like jailed/inactive, etc.
  xplaReceiver1: "rizon_receiver1_mail_address" # Receive all alert emails of xpla
  xplaReceiver2: "rizon_receiver2_mail_address" # Only receive emergency alert emails of xpla, like jailed/inactive, etc.

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
ExecStart=/home/ubuntu/go/bin/cosmos-monitor -c /data/config/.conf.yaml
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