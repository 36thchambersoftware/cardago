# cardago

Golang solutions for Cardano stake pools

# config.yaml

Create a file called config.yaml in the root project directory and fill in at least these settings:

    nodeCertPath: '/path/to/cert.file'
    leaderLogDirectory: '/path/to/logs'
    leaderLogPrefix: 'yourPrefix_'
    leaderLogExtension: 'txt'
    discord:
      authenticationToken: '123456789'
      serverID: '123456789'
      channelID: '123456789'
      userID: '123456789'

Valid config locations for running the binaries are `../../.` and `.`
