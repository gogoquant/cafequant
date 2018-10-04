package config

// API KEY
const (
	// todo: replace with your own AccessKey and Secret Key
	ACCESS_KEY string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	SECRET_KEY string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	// default to be disabled, please DON'T enable it unless it's officially announced.
	ENABLE_PRIVATE_SIGNATURE bool = false

	// generated the key by: openssl ecparam -name prime256v1 -genkey -noout -out privatekey.pem
	// only required when Private Signature is enabled
	// todo: replace with your own PrivateKey from privatekey.pem
	PRIVATE_KEY_PRIME_256 string = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDRalZO/ASg5cISM1jVKqBWzlg3xhEAULskUCD5yp/zVoAoGCCqGSM49
AwEHoUQDQgAEV8HcHMPhSSQWfo1LIx8W/RkimDVBFiQ6fDB/lHucIf9huu+uatTX
YLPI1ILcq6MqujUCzW1PhFSUFtVQqsUiHA==
-----END EC PRIVATE KEY-----`
)

// API请求地址, 不要带最后的/
const (
	//todo: replace with real URLs and HostName
	MARKET_URL string = "https://api.hadax.com/"
	TRADE_URL  string = ""
	HOST_NAME  string = ""
)
