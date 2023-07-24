# cosmos-monitor

This program is used to monitor specific validators on a blockchain deployed using the Cosmos SDK. Currently supports 21 blockchains including Cosmos Hub
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
  acrechainUser: "acrechain database user"
  acrechainPassword: "acrechain database password"
  acrechainName: "acrechain database name"
  acrechainHost: "acrechain database host ip"
  acrechainPort: "acrechain database host port"

  akashUser: "akash database user"
  akashPassword: "akash database password"
  akashName: "akash database name"
  akashHost: "akash database host ip"
  akashPort: "akash database host port"

  apolloUser: "apollo database user"
  apolloPassword: "apollo database password"
  apolloName: "apollo database name"
  apolloHost: "apollo database host ip"
  apolloPort: "apollo database host port"

  axelarUser: "axelar database user"
  axelarPassword: "axelar database password"
  axelarName: "axelar database name"
  axelarHost: "axelar database host ip"
  axelarPort: "axelar database host port"

  bandUser: "band database user"
  bandPassword: "band database password"
  bandName: "band database name"
  bandHost: "band database host ip"
  bandPort: "band database host port"

  cosmosUser: "cosmos database user"
  cosmosPassword: "cosmos database password"
  cosmosName: "cosmos database name"
  cosmosHost: "cosmos database host ip"
  cosmosPort: "cosmos database host port"

  evmosUser: "evmos database user"
  evmosPassword: "evmos database password"
  evmosName: "evmos database name"
  evmosHost: "evmos database host ip"
  evmosPort: "evmos database host port"

  gopherUser: "gopher database user"
  gopherPassword: "gopher database password"
  gopherName: "gopher database name"
  gopherHost: "gopher database host ip"
  gopherPort: "gopher database host port"

  heroUser: "hero database user"
  heroPassword: "hero database password"
  heroName: "hero database name"
  heroHost: "hero database host ip"
  heroPort: "hero database host port"

  injectiveUser: "injective database user"
  injectivePassword: "injective database password"
  injectiveName: "injective database name"
  injectiveHost: "injective database host ip"
  injectivePort: "injective database host port"

  junoUser: "juno database user"
  junoPassword: "juno database password"
  junoName: "juno database name"
  junoHost: "juno database host ip"
  junoPort: "juno database host port"

  neutronconsumerUser: "neutron consumer database user"
  neutronconsumerPassword: "neutron consumer database password"
  neutronconsumerName: "neutron consumer database name"
  neutronconsumerHost: "neutron consumer database host ip"
  neutronconsumerPort: "neutron consumer database host port"

  neutronUser: "neutron database user"
  neutronPassword: "neutron database password"
  neutronName: "neutron database name"
  neutronHost: "neutron database host ip"
  neutronPort: "neutron database host port"

  nyxUser: "nyx database user"
  nyxPassword: "nyx database password"
  nyxName: "nyx database name"
  nyxHost: "nyx database host ip"
  nyxPort: "nyx database host port"

  okp4User: "okp4 database user"
  okp4Password: "okp4 database password"
  okp4Name: "okp4 database name"
  okp4Host: "okp4 database host ip"
  okp4Port: "okp4 database host port"

  persistenceUser: "persistence database user"
  persistencePassword: "persistence database password"
  persistenceName: "persistence database name"
  persistenceHost: "persistence database host ip"
  persistencePort: "persistence database host port"

  providerUser: "provider database user"
  providerPassword: "provider database password"
  providerName: "provider database name"
  providerHost: "provider database host ip"
  providerPort: "provider database host port"

  rizonUser: "rizon database user"
  rizonPassword: "rizon database password"
  rizonName: "rizon database name"
  rizonHost: "rizon database host ip"
  rizonPort: "rizon database host port"

  secretUser: "secret database user"
  secretPassword: "secret database password"
  secretName: "secret database name"
  secretHost: "secret database host ip"
  secretPort: "secret database host port"

  sommelierUser: "sommelier database user"
  sommelierPassword: "sommelier database password"
  sommelierName: "sommelier database name"
  sommelierHost: "sommelier database host ip"
  sommelierPort: "sommelier database host port"

  sputnikUser: "sputnik database user"
  sputnikPassword: "sputnik database password"
  sputnikName: "sputnik database name"
  sputnikHost: "sputnik database host ip"
  sputnikPort: "sputnik database host port"

  teritoriUser: "teritori database user"
  teritoriPassword: "teritori database password"
  teritoriName: "teritori database name"
  teritoriHost: "teritori database host ip"
  teritoriPort: "teritori database host port"

  xplaUser: "xpla database user"
  xplaPassword: "xpla database password"
  xplaName: "xpla database name"
  xplaHost: "xpla database host ip"
  xplaPort: "xpla database host port"

  zetaUser: "zeta database user"
  zetaPassword: "zeta database password"
  zetaName: "zeta database name"
  zetaHost: "zeta database host ip"
  zetaPort: "zeta database host port"

# gRPC
gRpc:
  acrechainIp: "acrechain grpc host ip"
  acrechainPort: "acrechain grpc host grpc port"

  akashIp: "akash grpc host ip"
  akashPort: "akash grpc host grpc port"

  apolloIp: "apollo grpc host ip"
  apolloPort: "apollo grpc host grpc port"

  axelarIp: "axelar grpc host ip"
  axelarPort: "axelar grpc host grpc port"

  bandIp: "band grpc host ip"
  bandPort: "band grpc host grpc port"

  cosmosIp: "cosmos grpc host ip"
  cosmosPort: "cosmos grpc host grpc port"

  evmosIp: "evmos grpc host ip"
  evmosPort: "evmos grpc host grpc port"

  gopherIp: "gopher grpc host ip"
  gopherPort: "gopher grpc host grpc port"

  heroIp: "hero grpc host ip"
  heroPort: "hero grpc host grpc port"

  injectiveIp: "injective grpc host ip"
  injectivePort: "injective grpc host grpc port"

  junoIp: "juno grpc host ip"
  junoPort: "juno grpc host grpc port"

  neutronconsumerIp: "neutron consumer grpc host ip"
  neutronconsumerPort: "neutron consumer grpc host grpc port"

  neutronIp: "neutron grpc host ip"
  neutronPort: "neutron grpc host grpc port"

  nyxIp: "nyx grpc host ip"
  nyxPort: "nyx grpc host grpc port"

  okp4Ip: "okp4 grpc host ip"
  okp4Port: "okp4 grpc host grpc port"

  persistenceIp: "persistence grpc host ip"
  persistencePort: "persistence grpc host grpc port"

  providerIp: "provider grpc host ip"
  providerPort: "provider grpc host grpc port"

  rizonIp: "rizon grpc host ip"
  rizonPort: "rizon grpc host grpc port"

  secretIp: "secret grpc host ip"
  secretPort: "secret grpc host grpc port"

  sommelierIp: "sommelier grpc host ip"
  sommelierPort: "sommelier grpc host grpc port"

  sputnikIp: "sputnik grpc host ip"
  sputnikPort: "sputnik grpc host grpc port"

  teritoriIp: "teritori grpc host ip"
  teritoriPort: "teritori grpc host grpc port"

  xplaIp: "xpla grpc host ip"
  xplaPort: "xpla grpc host grpc port"

  zetaIp: "zeta grpc host ip"
  zetaPort: "zeta grpc host grpc port"

alert:
  acrechainOperatorAddr: "acrevaloperxxxxx,acrevaloperyyyy"    # The address starts with "acrevaloper"
  akashOperatorAddr: "akashvaloperxxxxx,akashvaloperyyyy"    # The address starts with "akashvaloper"
  apolloOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  axelarOperatorAddr: "axelarvaloperxxxxx,axelarvaloperyyyy"    # The address starts with "axelarvaloper"
  bandOperatorAddr: "bandvaloperxxxx,bandvaloperyyyy" # The address starts with "bandvaloper"
  cosmosOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  evmosOperatorAddr: "evmosvaloperxxxx,evmosvaloperyyyy" # The address starts with "evmosvaloper"
  gopherOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  heroOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  injectiveOperatorAddr: "injvaloperxxxxx,injvaloperyyyy"    # The address starts with "injvaloper"
  junoOperatorAddr: "junovaloperxxxxx,junovaloperyyyy"    # The address starts with "cosmosvaloper"
  neutronconsumerOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "neutronconsumer"
  neutronOperatorAddr: "neutronvaloperxxxx,neutronvaloperyyyy" # The address starts with "neutronvaloper"
  nyxOperatorAddr: "nvaloperxxxxx,nvaloperyyyy"    # The address starts with "nvaloper"
  okp4OperatorAddr: "okp4valoperxxxxx,okp4valoperyyyy"    # The address starts with "okp4valoper"
  persistenceOperatorAddr: "persistencevaloperxxxx,persistencevaloperyyyy" # The address starts with "persistencevaloper"
  providerOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  rizonOperatorAddr: "rizonvaloperxxxxx,rizonvaloperyyyy"    # The address starts with "rizonvaloper"
  secretOperatorAddr: "secretvaloperxxxxx,secretvaloperyyyy"    # The address starts with "secretvaloper"
  sommelierOperatorAddr: "sommeliervaloperxxxx,sommeliervaloperyyyy" # The address starts with "sommeliervaloper"
  sputnikOperatorAddr: "cosmosvaloperxxxxx,cosmosvaloperyyyy"    # The address starts with "cosmosvaloper"
  teritoriOperatorAddr: "teritorivaloperxxxxx,teritorivaloperyyyy"    # The address starts with "teritorivaloper"
  xplaOperatorAddr: "xplavaloperxxxxx,xplavaloperyyyy"    # The address starts with "xplavaloper"
  zetaOperatorAddr: "zetavaloperxxxxx,zetavaloperyyyy"    # The address starts with "zetavaloper"

  timeInterval: 600  # Monitoring time interval. 600 means the cosmos-monitor runs every 600 seconds
  blockInterval: 100 # Used to calculate the recent signature rate. 100 means to count the signatures rate of the last 100 blocks
  proportion: 0.05 # Signature rate

  acrechainRankingThreshold: 100    # The  acrechain validator ranking threshold
  akashRankingThreshold: 100    # The akash validator ranking threshold
  apolloRankingThreshold: 100    # The apollo validator ranking threshold
  axelarRankingThreshold: 100    # The axelar validator ranking threshold
  bandRankingThreshold: 100  # The band validator ranking threshold
  cosmosRankingThreshold: 100    # The cosmos validator ranking threshold
  evmosRankingThreshold: 100 # The evmos validator ranking threshold
  gopherRankingThreshold: 100 # The gopher validator ranking threshold
  heroRankingThreshold: 100 # The hero validator ranking threshold
  injectiveRankingThreshold: 100 # The injective validator ranking threshold
  junoRankingThreshold: 100  # The juno validator ranking threshold
  nyxRankingThreshold: 100   # The nyx validator ranking threshold
  okp4RankingThreshold: 100   # The okp4 validator ranking threshold
  neutronconsumerRankingThreshold: 100   # The neutron consumer validator ranking threshold
  neutronRankingThreshold: 100   # The neutron validator ranking threshold
  persistenceRankingThreshold: 100    # The persistence validator ranking threshold
  providerRankingThreshold: 100  # The provider validator ranking threshold
  rizonRankingThreshold: 100 # The rizon validator ranking threshold
  secretRankingThreshold: 100    # The secret validator ranking threshold
  sommelierRankingThreshold: 100 # The sommelier validator ranking threshold
  sputnikRankingThreshold: 100   # The sputnik validator ranking threshold
  teritoriRankingThreshold: 100  # The teritori validator ranking threshold
  xplaRankingThreshold: 100  # The xpla validator ranking threshold
  zetaRankingThreshold: 100  # The zeta validator ranking threshold

  #  Whether the project is monitored, if set true it will be monitored
  acrechainIsMonitored: true
  akashIsMonitored: true
  apolloIsMonitored: true
  axelarIsMonitored: true
  bandIsMonitored: true
  cosmosIsMonitored: true
  evmosIsMonitored: true
  gopherIsMonitored: true
  heroIsMonitored: true
  injectiveIsMonitored: true
  junoIsMonitored: true
  nyxIsMonitored: true
  okp4IsMonitored: true
  neutronconsumerIsMonitored: true
  neutronIsMonitored: true
  persistenceIsMonitored: true
  providerIsMonitored: true
  rizonIsMonitored: true
  secretIsMonitored: true
  sommelierIsMonitored: true
  sputnikIsMonitored: true
  teritoriIsMonitored: true
  xplaIsMonitored: true
  zetaIsMonitored: true

mail:
  host: "mail host"
  port: "mail port"
  username: "mail username"
  password: "mail password"
  sender: "sender mail address"
  acrechainReceiver1: "acrechain receiver1 mail address" # Receive all alert emails of acrechain
  acrechainReceiver2: "acrechain receiver2 mail address" # Receive all alert emails of acrechain
  akashReceiver1: "akash receiver1 mail address" # Receive all alert emails of akash
  akashReceiver2: "akash receiver2 mail address" # Receive all alert emails of akash
  apolloReceiver1: "apollo receiver1 mail address" # Receive all alert emails of apollo
  apolloReceiver2: "apollo receiver2 mail address" # Receive all alert emails of apollo
  axelarReceiver1: "axelar receiver1 mail address" # Receive all alert emails of axelar
  axelarReceiver2: "axelar receiver2 mail address" # Receive all alert emails of axelar
  bandReceiver1: "band receiver1 mail address" # Receive all alert emails of band
  bandReceiver2: "band receiver2 mail address" # Only receive emergency alert emails of band, like jailed/inactive, etc.
  cosmosReceiver1: "cosmos receiver1 mail address" # Receive all alert emails of cosmos
  cosmosReceiver2: "cosmos receiver2 mail address" # Only receive emergency alert emails of cosmos, like jailed/inactive, etc.
  evmosReceiver1: "evmos receiver1 mail address" # Receive all alert emails of evmos
  evmosReceiver2: "evmos receiver2 mail address" # Only receive emergency alert emails of evmos, like jailed/inactive, etc.
  gopherReceiver1: "gopher receiver1 mail address" # Receive all alert emails of gopher
  gopherReceiver2: "gopher receiver2 mail address" # Only receive emergency alert emails of gopher, like jailed/inactive, etc.
  heroReceiver1: "gopher receiver1 mail address" # Receive all alert emails of hero
  heroReceiver2: "gopher receiver2 mail address" # Only receive emergency alert emails of hero, like jailed/inactive, etc.
  injectiveReceiver1: "injective receiver1 mail address" # Receive all alert emails of injective
  injectiveReceiver2: "injective receiver2 mail address" # Only receive emergency alert emails of injective, like jailed/inactive, etc.
  junoReceiver1: "juno receiver1 mail address" # Receive all alert emails of juno
  junoReceiver2: "juno receiver2 mail address" # Only receive emergency alert emails of juno, like jailed/inactive, etc.
  neutronconsumerReceiver1: "neutron consumer receiver1 mail address" # Receive all alert emails of neutron consumer
  neutronconsumerReceiver2: "neutron consumer receiver2 mail address" # Only receive emergency alert emails of neutron consumer, like jailed/inactive, etc.
  neutronReceiver1: "neutron receiver1 mail address" # Receive all alert emails of neutron
  neutronReceiver2: "neutron receiver2 mail address" # Only receive emergency alert emails of neutron, like jailed/inactive, etc.
  nyxReceiver1: "nyx receiver1 mail address" # Receive all alert emails of nyx
  nyxReceiver2: "nyx receiver2 mail address" # Only receive emergency alert emails of nyx, like jailed/inactive, etc.
  okp4Receiver1: "okp4 receiver1 mail address" # Receive all alert emails of okp4
  okp4Receiver2: "okp4 receiver2 mail address" # Only receive emergency alert emails of okp4, like jailed/inactive, etc.
  persistenceReceiver1: "persistence receiver1 mail address" # Receive all alert emails of persistence
  persistenceReceiver2: "persistence receiver2 mail address" # Only receive emergency alert emails of persistence, like jailed/inactive, etc.
  providerReceiver1: "provider receiver1 mail address" # Receive all alert emails of provider
  providerReceiver2: "provider receiver2 mail address" # Only receive emergency alert emails of provider, like jailed/inactive, etc.
  rizonReceiver1: "rizon receiver1 mail address" # Receive all alert emails of rizon
  rizonReceiver2: "rizon receiver2 mail address" # Only receive emergency alert emails of rizon, like jailed/inactive, etc.
  secretReceiver1: "secret receiver1 mail address" # Receive all alert emails of secret
  secretReceiver2: "secret receiver2 mail address" # Only receive emergency alert emails of secret, like jailed/inactive, etc.
  sommelierReceiver1: "sommelier receiver1 mail address" # Receive all alert emails of sommelier
  sommelierReceiver2: "sommelier receiver2 mail address" # Only receive emergency alert emails of sommelier, like jailed/inactive, etc.
  sputnikReceiver1: "sputnik receiver1 mail address" # Receive all alert emails of sputnik
  sputnikReceiver2: "sputnik receiver2 mail address" # Only receive emergency alert emails of sputnik, like jailed/inactive, etc.
  teritoriReceiver1: "teritori receiver1 mail address" # Receive all alert emails of teritori
  teritoriReceiver2: "teritori receiver2 mail address" # Only receive emergency alert emails of teritori, like jailed/inactive, etc.
  xplaReceiver1: "xpla receiver1 mail address" # Receive all alert emails of xpla
  xplaReceiver2: "xpla receiver2 mail address" # Only receive emergency alert emails of xpla, like jailed/inactive, etc.
  zetaReceiver1: "zeta receiver1 mail address" # Receive all alert emails of zeta
  zetaReceiver2: "zeta receiver2 mail address" # Only receive emergency alert emails of zeta, like jailed/inactive, etc.

log:
  path: "log save location"
  level: "log level" # info, error
  eventlogpath: "event log save location"
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