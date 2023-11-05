# cardago

Golang solutions for Cardano stake pools

# config.yaml

Create a file called config.yaml in the root project directory and fill in at least these settings:

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
      authenticationToken: '123456789'
      serverID: '123456789'
      channelID: '123456789'
      voiceChannelID: '123456789'
      userID: '123456789'

Valid config locations for running the binaries are `../../.` and `.`
