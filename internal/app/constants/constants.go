package constants

const (
	DevEnvironment = "dev"

	ServerHost           = "HTTPServer.listen"
	ServerPort           = "HTTPServer.port"
	ReadTimeout          = "HTTPServer.readTimeout"
	WriteTimeout         = "HTTPServer.writeTimeout"
	MaxConnsPerIP        = "HTTPServer.maxConnsPerIP"
	MaxRequestsPerConn   = "HTTPServer.maxRequestsPerConn"
	MaxKeepaliveDuration = "HTTPServer.maxKeepaliveDuration"

	JWTExpiration                 = "JWT.expiration"
	JWTRefreshExpiration          = "JWT.refresh_expiration"
	JWTIssuer                     = "JWT.issuer"
	JWTAccessTokenPublicKeyPath   = "JWT.access_token_public_key_path"
	JWTAccessTokenPrivateKeyPath  = "JWT.access_token_private_key_path"
	JWTRefreshTokenPublicKeyPath  = "JWT.refresh_token_public_key_path"
	JWTRefreshTokenPrivateKeyPath = "JWT.refresh_token_private_key_path"

	CorrelationId = "X-Correlation-ID"
	CtxClaimKey   = "context_claim"

	GetJWTTokenPath = "/token"

	OTPDURATION  = "OtpDuration"
	OTPSECRETKEY = "OtpSecretKey"
)

var CODEVALIDITY = 900
var DBLOGMODE bool

const (
	TEMPLATE_ID__ACCOUNT_VERIFICATION = "d-be66506f47f74fba979a9af69eb2dfc1"
	TEMPLATE_ID_PASSWORD_RESET        = "d-d71cbf79d2964c9d92f1361e00c87c31 "
)
