package constants

const (
	TEST_SENT_URI = `https://api-rinkeby.etherscan.io/api?apikey=Y5BJ5VA3XZ59F63XQCQDDUWU2C29144MMM&module=logs&action=getLogs&fromBlock=0&toBlock=latest&address=0x29317B796510afC25794E511e7B10659Ca18048B&topic0=0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef&topic0_1_opr=and&topic1=`
	TEST_SENT_URI2 = `&topic1_2_opr=or&topic2=`
	ZFill = "000000000000000000000000"
	DANTE_MN = "https://ton.sentinelgroup.io"
	Success = "Congratulations!! please click the button below to connect to the sentinel dVPN node and next time use /mynode to access this node"
	CheckWalletOptionsError = "error while fetching user wallet address. in case you have not attached your wallet address, please share your wallet address again."
	WarningMsg = `Please do not share this url with anyone else. 
who ever has this url, can access your node and hence use the bandwidth assigned to you`
)