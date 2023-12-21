###
# Sync Progress
###

buildSyncProgress:
	cd cmd/syncprogress; GOOS=linux GOARCH=amd64 go build -o syncprogress.bin

uploadSyncProgress:
	scp -P 7811 -i /Users/priebe/.ssh/cardano_mainnet cmd/syncprogress/syncprogress.bin cardano@$$PREEBCORE:/home/cardano/scripts/cardago/syncprogress

refreshProdSyncProgress: buildSyncProgress uploadSyncProgress

###
# KES Period Info
###

buildKESPeriodInfo:
	cd cmd/kesperiodinfo; GOOS=linux GOARCH=amd64 go build -o kesperiodinfo.bin

uploadKESPeriodInfo:
	scp -P 7811 -i /Users/priebe/.ssh/cardano_mainnet cmd/kesperiodinfo/kesperiodinfo.bin cardano@$$PREEBCORE:/home/cardano/scripts/cardago/kesperiodinfo

refreshKESPeriodInfo: buildKESPeriodInfo uploadKESPeriodInfo

###
# Scheduled Blocks
###

buildScheduledBlocks:
	cd cmd/scheduledblocks; GOOS=linux GOARCH=amd64 go build -o scheduledblocks.bin

uploadScheduledBlocks:
	scp -P 7811 -i /Users/priebe/.ssh/cardano_mainnet cmd/scheduledblocks/scheduledblocks.bin cardano@$$PREEBCORE:/home/cardano/scripts/cardago/blocksthisepoch

refreshScheduledBlocks: buildScheduledBlocks uploadScheduledBlocks