//go:build legacy
// +build legacy

package models

import (
	"time"
)

// BlockchainComponent represents a component registered on blockchain
type BlockchainComponent struct {
	ID                 int64     `json:"id"`
	ComponentID        int64     `json:"component_id"`
	CreatorID          int64     `json:"creator_id"`
	BlockchainAddress  string    `json:"blockchain_address"` // Smart contract address
	TransactionHash    string    `json:"transaction_hash"`
	ContractVersion    string    `json:"contract_version"`   // "1.0", "2.0"
	ComponentHash      string    `json:"component_hash"`     // SHA-256 of component
	RegistryTimestamp  time.Time `json:"registry_timestamp"`
	VerificationStatus string    `json:"verification_status"` // "verified", "pending", "failed"
	RightsType         string    `json:"rights_type"`         // "full", "limited", "commercial"
	ExclusiveRights    bool      `json:"exclusive_rights"`
	LicenseType        string    `json:"license_type"`        // "MIT", "Apache", "GPL", "Proprietary"
	RoyaltyPercentage  float64   `json:"royalty_percentage"`  // Creator royalty (0-100%)
	TransactionFee     float64   `json:"transaction_fee"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// SmartContract represents a deployed smart contract
type SmartContract struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	DeployerID      int64     `json:"deployer_id"`
	ContractType    string    `json:"contract_type"`   // "component_registry", "nft", "rights_management"
	ABI             string    `json:"abi"`             // Contract ABI (JSON)
	ByteCode        string    `json:"byte_code"`       // Compiled bytecode
	SourceCode      string    `json:"source_code"`     // Solidity source
	Network         string    `json:"network"`         // "ethereum", "polygon", "arbitrum"
	NetworkID       int64     `json:"network_id"`      // Chain ID
	Status          string    `json:"status"`          // "deployed", "deprecated", "testing"
	Version         string    `json:"version"`
	TransactionHash string    `json:"transaction_hash"`
	DeployedAt      time.Time `json:"deployed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NFTMint represents an NFT minted for a component
type NFTMint struct {
	ID              int64     `json:"id"`
	ComponentID     int64     `json:"component_id"`
	CreatorID       int64     `json:"creator_id"`
	TokenID         string    `json:"token_id"`        // NFT token ID
	ContractAddress string    `json:"contract_address"` // NFT contract address
	MetadataURI     string    `json:"metadata_uri"`    // IPFS/Arweave URI
	MintAddress     string    `json:"mint_address"`    // Creator's wallet address
	OwnerAddress    string    `json:"owner_address"`   // Current owner's wallet
	Status          string    `json:"status"`          // "minted", "listed", "sold", "burned"
	MintPrice       float64   `json:"mint_price"`      // Original mint price
	RoyaltyRate     float64   `json:"royalty_rate"`    // Royalty percentage
	TransactionHash string    `json:"transaction_hash"`
	ChainID         int64     `json:"chain_id"`
	MintedAt        time.Time `json:"minted_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NFTListing represents an NFT listing for sale
type NFTListing struct {
	ID               int64     `json:"id"`
	NFTMintID        int64     `json:"nft_mint_id"`
	TokenID          string    `json:"token_id"`
	ListingAddress   string    `json:"listing_address"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`       // "ETH", "USDC", "MATIC"
	Status           string    `json:"status"`         // "active", "sold", "cancelled"
	Marketplace      string    `json:"marketplace"`    // "opensea", "rarible", "custom"
	MarketplaceURL   string    `json:"marketplace_url"`
	ExpiresAt        *time.Time `json:"expires_at"`
	TransactionHash  string    `json:"transaction_hash"`
	ListedAt         time.Time `json:"listed_at"`
	SoldAt           *time.Time `json:"sold_at"`
	BuyerAddress     *string    `json:"buyer_address"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RightsVerification verifies component rights on blockchain
type RightsVerification struct {
	ID                int64     `json:"id"`
	ComponentID       int64     `json:"component_id"`
	VerifierAddress   string    `json:"verifier_address"`   // Address that verified
	VerificationType  string    `json:"verification_type"`  // "ownership", "usage_rights", "commercial"
	IsVerified        bool      `json:"is_verified"`
	VerificationData  string    `json:"verification_data"`  // JSON metadata
	TransactionHash   string    `json:"transaction_hash"`
	ExpiresAt         *time.Time `json:"expires_at"`        // For time-limited rights
	RevokedAt         *time.Time `json:"revoked_at"`
	VerifiedAt        time.Time `json:"verified_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// WalletAddress represents a blockchain wallet
type WalletAddress struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Address        string    `json:"address"`      // Ethereum/Polygon address
	ChainID        int64     `json:"chain_id"`
	Network        string    `json:"network"`      // "ethereum", "polygon", "arbitrum"
	IsVerified     bool      `json:"is_verified"`
	VerificationNonce string `json:"verification_nonce"` // For signing verification
	SignedMessage  string    `json:"signed_message"`
	IsPrimary      bool      `json:"is_primary"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ComponentTransfer tracks ownership transfer of components
type ComponentTransfer struct {
	ID              int64     `json:"id"`
	ComponentID     int64     `json:"component_id"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	TransactionHash string    `json:"transaction_hash"`
	TransactionFee  float64   `json:"transaction_fee"`
	TransferPrice   *float64  `json:"transfer_price"`
	RoyaltiesPaid   float64   `json:"royalties_paid"`
	CreatorRoyalty  *float64  `json:"creator_royalty"`
	Status          string    `json:"status"`     // "pending", "confirmed", "failed"
	BlockNumber     int64     `json:"block_number"`
	TransferredAt   time.Time `json:"transferred_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// ComponentRights manages usage and commercial rights
type ComponentRights struct {
	ID                  int64     `json:"id"`
	ComponentID         int64     `json:"component_id"`
	OwnerAddress        string    `json:"owner_address"`
	CanUseCommercially  bool      `json:"can_use_commercially"`
	CanModify           bool      `json:"can_modify"`
	CanResell           bool      `json:"can_resell"`
	CanLicense          bool      `json:"can_license"`
	CanBurn             bool      `json:"can_burn"`
	ExclusiveOwner      bool      `json:"exclusive_owner"`   // True if exclusive rights
	RightsExpiresAt     *time.Time `json:"rights_expires_at"`
	UsageCount          int64     `json:"usage_count"`
	MaxUsageAllowed     *int64    `json:"max_usage_allowed"`
	RoyaltyPercentage   float64   `json:"royalty_percentage"`
	RoyaltyRecipient    string    `json:"royalty_recipient"` // Wallet address
	TransactionHash     string    `json:"transaction_hash"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// BlockchainTransaction tracks all blockchain transactions
type BlockchainTransaction struct {
	ID              int64     `json:"id"`
	TransactionHash string    `json:"transaction_hash"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	Method          string    `json:"method"`       // "mint", "transfer", "burn", "verify"
	Value           float64   `json:"value"`
	GasUsed         int64     `json:"gas_used"`
	GasPrice        float64   `json:"gas_price"`
	Status          string    `json:"status"`       // "pending", "confirmed", "failed"
	ErrorMessage    *string   `json:"error_message"`
	BlockNumber     int64     `json:"block_number"`
	Timestamp       time.Time `json:"timestamp"`
	Network         string    `json:"network"`
	ChainID         int64     `json:"chain_id"`
	RelatedEntity   string    `json:"related_entity"` // "component", "nft", "rights"
	RelatedID       int64     `json:"related_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// BlockchainConfig holds blockchain configuration
type BlockchainConfig struct {
	ID                  int64  `json:"id"`
	Network             string `json:"network"`  // "ethereum", "polygon", "arbitrum"
	ChainID             int64  `json:"chain_id"`
	RpcEndpoint         string `json:"rpc_endpoint"`
	NetworkName         string `json:"network_name"`
	CurrencySymbol      string `json:"currency_symbol"` // "ETH", "MATIC", "ARB"
	ContractRegistryAddress string `json:"contract_registry_address"`
	ContractNFTAddress  string `json:"contract_nft_address"`
	ContractRightsAddress string `json:"contract_rights_address"`
	GasLimit            int64  `json:"gas_limit"`
	GasPrice            float64 `json:"gas_price"`
	IsEnabled           bool   `json:"is_enabled"`
	IsTestnet           bool   `json:"is_testnet"`
}

// ComponentRegistry represents blockchain component registry metadata
type ComponentRegistry struct {
	ID              int64      `json:"id"`
	ComponentID     int64      `json:"component_id"`
	BlockchainID    string     `json:"blockchain_id"`      // Blockchain-side ID
	ContentHash     string     `json:"content_hash"`       // SHA-256 of component
	MetadataHash    string     `json:"metadata_hash"`
	RegistryAddress string     `json:"registry_address"`
	CreatorAddress  string     `json:"creator_address"`
	Verified        bool       `json:"verified"`
	VerifiedAt      *time.Time `json:"verified_at"`
	VerifierAddress *string    `json:"verifier_address"`
	Version         int32      `json:"version"`
	PreviousHash    *string    `json:"previous_hash"`
	TransactionHash string     `json:"transaction_hash"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// SmartContractInteraction tracks calls to smart contracts
type SmartContractInteraction struct {
	ID              int64     `json:"id"`
	SmartContractID int64     `json:"smart_contract_id"`
	CallerAddress   string    `json:"caller_address"`
	Method          string    `json:"method"`
	Parameters      string    `json:"parameters"`     // JSON serialized params
	ReturnValue     string    `json:"return_value"`   // JSON serialized return
	Status          string    `json:"status"`         // "pending", "succeeded", "failed"
	ErrorMessage    *string   `json:"error_message"`
	TransactionHash string    `json:"transaction_hash"`
	BlockNumber     int64     `json:"block_number"`
	GasUsed         int64     `json:"gas_used"`
	CalledAt        time.Time `json:"called_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// TokenMetadata represents NFT metadata
type TokenMetadata struct {
	ID              int64  `json:"id"`
	TokenID         string `json:"token_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageURI        string `json:"image_uri"`
	ExternalURL     string `json:"external_url"`
	Attributes      string `json:"attributes"`     // JSON serialized attributes
	AnimationURL    *string `json:"animation_url"`
	YouTube         *string `json:"youtube_url"`
	RoyaltyInfo     string `json:"royalty_info"`  // JSON with percentage and recipient
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
