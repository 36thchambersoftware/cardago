# Cardago

Golang solutions for Cardano stake pools.

## Commands

Currently there are 3 features.

1. scheduledblocks is a command that will use your pool info to check for scheduled blocks in the next epoch. This is designed to be used as part of a cron that will check on a regular basis. Common blockers to this command running would be your pool being out of sync, you're using the wrong shelley-genesis file, the epoch isn't yet 75% complete, or the command has already run successfully for the next epoch (meaning a log file has already been created for the next epoch)
2. kesperiodinfo is for alerting you when your kes key rotation needs to be completed by.
3. syncprogress is ideal for alerting a discord server that your block producer is online and fully synced. This is a way to give delegators/discord members peace of mind that the pool is up and running.

## Run Locally

Clone the project

```bash
  gh repo clone 36thchambersoftware/cardago
```

Go to desired command directory

```bash
  cd cardago/cmd/scheduledblocks
```

Create the binary for your system (Linux shown, though this is a good resource https://freshman.tech/snippets/go/cross-compile-go-programs/)

```bash
  GOOS=linux GOARCH=amd64 go build
```

Create a config.yaml file in the same directory as the command.

```bash
logs:
  leader:
    directory: '/home/cardano/cardano/logs'
    prefix: 'leaderSchedule_'
    extension: ''
cardano:
  nodeCertPath: '/home/cardano/cardano/node.cert'
  shelleyGenesisFilePath: '/path/to/shelley-genesis.json'
  stakepoolid: '<stake pool id>'
  vrfskeyfilepath: '/path/to/vrf.skey'
discord:
  webhookURL: 'https://your.discord/webhook/url?wait=true'
  userID: '123456789'
```

Run the program

```bash
  ./scheduledblocks
```

## Acknowledgements

- [PREEB Pool](https://preeb.cloud)

## Authors

- [@36thchambersoftware](https://github.com/36thchambersoftware)

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)
[![AGPL License](https://img.shields.io/badge/license-AGPL-blue.svg)](http://www.gnu.org/licenses/agpl-3.0)
