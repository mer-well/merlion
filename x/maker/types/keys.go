package types

const (
	// ModuleName defines the module name
	ModuleName = "maker"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_maker"
)

const (
	prefixCollateralRatio = iota + 1
	prefixCollateralRatioLastBlock
	prefixBackingParams
	prefixCollateralParams
	prefixBackingTotal
	prefixCollateralTotal
	prefixBackingPool
	prefixCollateralPool
	prefixBackingAccount
	prefixCollateralAccount
)

var (
	KeyPrefixCollateralRatio          = []byte{prefixCollateralRatio}
	KeyPrefixCollateralRatioLastBlock = []byte{prefixCollateralRatioLastBlock}
	KeyPrefixBackingParams            = []byte{prefixBackingParams}
	KeyPrefixCollateralParams         = []byte{prefixCollateralParams}
	KeyPrefixBackingTotal             = []byte{prefixBackingTotal}
	KeyPrefixCollateralTotal          = []byte{prefixCollateralTotal}
	KeyPrefixBackingPool              = []byte{prefixBackingPool}
	KeyPrefixCollateralPool           = []byte{prefixCollateralPool}
	KeyPrefixBackingAccount           = []byte{prefixBackingAccount}
	KeyPrefixCollateralAccount        = []byte{prefixCollateralAccount}
)
