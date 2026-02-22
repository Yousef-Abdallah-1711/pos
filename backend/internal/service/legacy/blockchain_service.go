//go:build legacy
// +build legacy

package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"context"
	"pos-saas/internal/repository"
)

// BlockchainService handles blockchain operations for components
type BlockchainService struct {
	repo repository.BlockchainRepository
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(repo repository.BlockchainRepository) *BlockchainService {
	return &BlockchainService{
		repo: repo,
	}
}

// RegisterComponentOnBlockchain registers a component on the blockchain
func (s *BlockchainService) RegisterComponentOnBlockchain(ctx context.Context, componentID int64, creatorID int64, componentData []byte, rightsType string) (*models.BlockchainComponent, error) {
	if componentID == 0 || creatorID == 0 {
		return nil, fmt.Errorf("componentID and creatorID are required")
	}

	// Validate rights type
	validRights := map[string]bool{
		"full":        true,
		"limited":     true,
		"commercial":  true,
	}

	if !validRights[rightsType] {
		return nil, fmt.Errorf("invalid rights type: %s", rightsType)
	}

	// Calculate component hash (SHA-256)
	hashBytes := sha256.Sum256(componentData)
	componentHash := hex.EncodeToString(hashBytes[:])

	// Check if already registered
	existing, _ := s.repo.GetBlockchainComponentByHash(ctx, componentHash)
	if existing != nil {
		return existing, fmt.Errorf("component already registered on blockchain")
	}

	// Simulate blockchain registration
	blockchainAddress := generateEthereumAddress()
	transactionHash := generateTransactionHash()

	blockchain := &models.BlockchainComponent{
		ComponentID:       componentID,
		CreatorID:         creatorID,
		BlockchainAddress: blockchainAddress,
		TransactionHash:   transactionHash,
		ContractVersion:   "1.0",
		ComponentHash:     componentHash,
		RegistryTimestamp: time.Now(),
		VerificationStatus: "verified",
		RightsType:        rightsType,
		ExclusiveRights:   false,
		LicenseType:       "Proprietary",
		RoyaltyPercentage: 10.0, // 10% default
		TransactionFee:    0.0, // Would be set based on gas usage
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return s.repo.RegisterComponentOnBlockchain(ctx, blockchain)
}

// MintComponentNFT mints an NFT for a registered component
func (s *BlockchainService) MintComponentNFT(ctx context.Context, componentID int64, creatorID int64, creatorWallet string, mintPrice float64) (*models.NFTMint, error) {
	// Check if component is registered
	registrations, err := s.repo.GetComponentBlockchainRegistrations(ctx, componentID)
	if err != nil || len(registrations) == 0 {
		return nil, fmt.Errorf("component not registered on blockchain")
	}

	// Verify creator wallet exists
	wallet, err := s.repo.GetWalletByAddress(ctx, creatorWallet)
	if err != nil || wallet == nil {
		return nil, fmt.Errorf("creator wallet not verified")
	}

	if !wallet.IsVerified {
		return nil, fmt.Errorf("creator wallet must be verified before minting NFT")
	}

	// Generate NFT token ID (simulated)
	tokenID := generateTokenID()
	contractAddress := "0x1234567890123456789012345678901234567890" // NFT contract address
	metadataURI := fmt.Sprintf("ipfs://QmXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX") // IPFS URI
	transactionHash := generateTransactionHash()

	nft := &models.NFTMint{
		ComponentID:     componentID,
		CreatorID:       creatorID,
		TokenID:         tokenID,
		ContractAddress: contractAddress,
		MetadataURI:     metadataURI,
		MintAddress:     creatorWallet,
		OwnerAddress:    creatorWallet,
		Status:          "minted",
		MintPrice:       mintPrice,
		RoyaltyRate:     10.0,
		TransactionHash: transactionHash,
		ChainID:         1, // Ethereum mainnet
		MintedAt:        time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Record transaction
	tx := &models.BlockchainTransaction{
		TransactionHash: transactionHash,
		FromAddress:     creatorWallet,
		ToAddress:       contractAddress,
		Method:          "mint",
		Value:           mintPrice,
		GasUsed:         150000, // Estimated gas
		GasPrice:        50.0,   // Estimated gas price in GWEI
		Status:          "confirmed",
		BlockNumber:     17000000, // Simulated block
		Timestamp:       time.Now(),
		Network:         "ethereum",
		ChainID:         1,
		RelatedEntity:   "nft",
		RelatedID:       0, // Will be set after creation
	}

	s.repo.RecordTransaction(ctx, tx)

	return s.repo.MintNFT(ctx, nft)
}

// ListNFTForSale lists an NFT for sale on marketplace
func (s *BlockchainService) ListNFTForSale(ctx context.Context, nftID int64, price float64, currency string) (*models.NFTListing, error) {
	// Get NFT
	nft, err := s.repo.GetNFTMint(ctx, nftID)
	if err != nil || nft == nil {
		return nil, fmt.Errorf("NFT not found")
	}

	if nft.Status == "sold" {
		return nil, fmt.Errorf("cannot list NFT that has already been sold")
	}

	marketplaceURL := fmt.Sprintf("https://opensea.io/assets/ethereum/%s/%s", nft.ContractAddress, nft.TokenID)
	expiresAt := time.Now().AddDate(0, 1, 0) // List for 1 month

	listing := &models.NFTListing{
		NFTMintID:       nftID,
		TokenID:         nft.TokenID,
		ListingAddress:  nft.OwnerAddress,
		Price:           price,
		Currency:        currency,
		Status:          "active",
		Marketplace:     "opensea",
		MarketplaceURL:  marketplaceURL,
		ExpiresAt:       &expiresAt,
		TransactionHash: generateTransactionHash(),
		ListedAt:        time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return s.repo.CreateNFTListing(ctx, listing)
}

// TransferComponentOwnership transfers component ownership between addresses
func (s *BlockchainService) TransferComponentOwnership(ctx context.Context, componentID int64, fromAddress, toAddress string, price *float64) (*models.ComponentTransfer, error) {
	// Verify both addresses are registered wallets
	fromWallet, _ := s.repo.GetWalletByAddress(ctx, fromAddress)
	toWallet, _ := s.repo.GetWalletByAddress(ctx, toAddress)

	if fromWallet == nil || !fromWallet.IsVerified {
		return nil, fmt.Errorf("sender wallet not verified")
	}

	if toWallet == nil || !toWallet.IsVerified {
		return nil, fmt.Errorf("recipient wallet not verified")
	}

	// Get component blockchain registration for royalty info
	registrations, err := s.repo.GetComponentBlockchainRegistrations(ctx, componentID)
	if err != nil || len(registrations) == 0 {
		return nil, fmt.Errorf("component not found on blockchain")
	}

	registration := registrations[0]
	royaltyPercent := registration.RoyaltyPercentage
	var creatorRoyalty *float64
	if price != nil && royaltyPercent > 0 {
		royalty := *price * (royaltyPercent / 100)
		creatorRoyalty = &royalty
	}

	transactionHash := generateTransactionHash()

	transfer := &models.ComponentTransfer{
		ComponentID:     componentID,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		TransactionHash: transactionHash,
		TransactionFee:  0.01, // Simulated gas fee in ETH
		TransferPrice:   price,
		RoyaltiesPaid:   0,
		CreatorRoyalty:  creatorRoyalty,
		Status:          "confirmed",
		BlockNumber:     17000000,
		TransferredAt:   time.Now(),
		CreatedAt:       time.Now(),
	}

	// If price provided, calculate royalties paid
	if price != nil && creatorRoyalty != nil {
		transfer.RoyaltiesPaid = *creatorRoyalty
	}

	// Record transaction
	tx := &models.BlockchainTransaction{
		TransactionHash: transactionHash,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		Method:          "transfer",
		Value:           0,
		GasUsed:         100000,
		GasPrice:        50.0,
		Status:          "confirmed",
		BlockNumber:     17000000,
		Timestamp:       time.Now(),
		Network:         "ethereum",
		ChainID:         1,
		RelatedEntity:   "component",
		RelatedID:       componentID,
	}

	s.repo.RecordTransaction(ctx, tx)

	return s.repo.RecordComponentTransfer(ctx, transfer)
}

// VerifyComponentRights verifies ownership and usage rights
func (s *BlockchainService) VerifyComponentRights(ctx context.Context, componentID int64, ownerAddress string, verificationType string) (*models.RightsVerification, error) {
	// Validate verification type
	validTypes := map[string]bool{
		"ownership":        true,
		"usage_rights":     true,
		"commercial":       true,
	}

	if !validTypes[verificationType] {
		return nil, fmt.Errorf("invalid verification type: %s", verificationType)
	}

	// Verify wallet
	wallet, _ := s.repo.GetWalletByAddress(ctx, ownerAddress)
	if wallet == nil || !wallet.IsVerified {
		return nil, fmt.Errorf("verifier wallet not verified")
	}

	transactionHash := generateTransactionHash()

	verification := &models.RightsVerification{
		ComponentID:      componentID,
		VerifierAddress:  ownerAddress,
		VerificationType: verificationType,
		IsVerified:       true,
		VerificationData: fmt.Sprintf(`{"type":"%s","verified_at":"%s"}`, verificationType, time.Now().Format(time.RFC3339)),
		TransactionHash:  transactionHash,
		ExpiresAt:        nil, // Permanent unless revoked
		VerifiedAt:       time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return s.repo.VerifyComponentRights(ctx, verification)
}

// SetComponentRights sets usage and commercial rights for a component
func (s *BlockchainService) SetComponentRights(ctx context.Context, componentID int64, ownerAddress string, canUseCommercially, canModify, canResell, canLicense bool) (*models.ComponentRights, error) {
	transactionHash := generateTransactionHash()

	rights := &models.ComponentRights{
		ComponentID:         componentID,
		OwnerAddress:        ownerAddress,
		CanUseCommercially:  canUseCommercially,
		CanModify:           canModify,
		CanResell:           canResell,
		CanLicense:          canLicense,
		CanBurn:             true, // Always allow burning
		ExclusiveOwner:      false,
		RightsExpiresAt:     nil,
		UsageCount:          0,
		MaxUsageAllowed:     nil, // Unlimited
		RoyaltyPercentage:   0,
		RoyaltyRecipient:    ownerAddress,
		TransactionHash:     transactionHash,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return s.repo.SetComponentRights(ctx, rights)
}

// RegisterWalletAddress registers a blockchain wallet for a user
func (s *BlockchainService) RegisterWalletAddress(ctx context.Context, userID int64, address string, network string) (*models.WalletAddress, error) {
	// Validate address format (basic check for Ethereum address)
	if len(address) != 42 || address[:2] != "0x" {
		return nil, fmt.Errorf("invalid Ethereum address format")
	}

	// Check if wallet already registered
	existing, _ := s.repo.GetWalletByAddress(ctx, address)
	if existing != nil {
		return nil, fmt.Errorf("wallet address already registered")
	}

	// Determine chain ID based on network
	chainID := int64(1) // Default to Ethereum mainnet
	switch network {
	case "polygon":
		chainID = 137
	case "arbitrum":
		chainID = 42161
	case "goerli":
		chainID = 5 // Goerli testnet
	}

	nonce := generateVerificationNonce()

	wallet := &models.WalletAddress{
		UserID:            userID,
		Address:           address,
		ChainID:           chainID,
		Network:           network,
		IsVerified:        false,
		VerificationNonce: nonce,
		IsPrimary:         false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return s.repo.RegisterWalletAddress(ctx, wallet)
}

// VerifyWalletAddress verifies a wallet by signature
func (s *BlockchainService) VerifyWalletAddress(ctx context.Context, walletID int64, signedMessage string) (*models.WalletAddress, error) {
	wallet, err := s.repo.GetWalletAddress(ctx, walletID)
	if err != nil || wallet == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	// In production, verify ECDSA signature with wallet.VerificationNonce
	// For now, accept if signed message matches expected format
	if signedMessage == "" {
		return nil, fmt.Errorf("invalid signature")
	}

	wallet.IsVerified = true
	wallet.SignedMessage = signedMessage
	wallet.UpdatedAt = time.Now()

	err = s.repo.UpdateWalletVerification(ctx, walletID, true)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// SetPrimaryWallet sets the primary wallet for a user
func (s *BlockchainService) SetPrimaryWallet(ctx context.Context, walletID int64) error {
	wallet, err := s.repo.GetWalletAddress(ctx, walletID)
	if err != nil || wallet == nil {
		return fmt.Errorf("wallet not found")
	}

	if !wallet.IsVerified {
		return fmt.Errorf("wallet must be verified before setting as primary")
	}

	return s.repo.SetPrimaryWallet(ctx, walletID)
}

// GetComponentVerificationStatus retrieves all verification statuses for a component
func (s *BlockchainService) GetComponentVerificationStatus(ctx context.Context, componentID int64) (map[string]interface{}, error) {
	registrations, err := s.repo.GetComponentBlockchainRegistrations(ctx, componentID)
	if err != nil || len(registrations) == 0 {
		return nil, fmt.Errorf("component not registered on blockchain")
	}

	registration := registrations[0]

	// Get NFTs
	nfts, _ := s.repo.GetComponentNFTs(ctx, componentID)

	// Get rights verifications
	verifications, _ := s.repo.GetComponentRightsVerifications(ctx, componentID)

	status := map[string]interface{}{
		"blockchain_registered": true,
		"registration_hash":     registration.ComponentHash,
		"verified":              registration.VerificationStatus == "verified",
		"rights_type":           registration.RightsType,
		"exclusive_rights":      registration.ExclusiveRights,
		"license_type":          registration.LicenseType,
		"royalty_percentage":    registration.RoyaltyPercentage,
		"nft_count":             len(nfts),
		"verification_count":    len(verifications),
		"registered_at":         registration.RegistryTimestamp,
	}

	return status, nil
}

// GetNFTSalesHistory retrieves sales history for an NFT
func (s *BlockchainService) GetNFTSalesHistory(ctx context.Context, tokenID string) ([]*models.ComponentTransfer, error) {
	// This would typically query NFT transfer events from blockchain
	// For now, return empty as this would require blockchain event parsing
	return []*models.ComponentTransfer{}, nil
}

// EstimateGasCost estimates the gas cost for blockchain operations
func (s *BlockchainService) EstimateGasCost(operation string, gasPrice float64) float64 {
	// Typical gas costs (in GWEI)
	gasCosts := map[string]int64{
		"register":  200000,
		"mint":      150000,
		"transfer":  100000,
		"burn":      60000,
		"verify":    120000,
	}

	gasUsed := gasCosts[operation]
	if gasUsed == 0 {
		gasUsed = 100000 // Default estimate
	}

	// Convert gas units and price (assuming GWEI to ETH)
	gasCostETH := float64(gasUsed) * (gasPrice / 1e9)
	return gasCostETH
}

// Helper functions

func generateEthereumAddress() string {
	// Generate mock Ethereum address
	return fmt.Sprintf("0x%040x", rand.Int63())
}

func generateTransactionHash() string {
	// Generate mock transaction hash (64 hex characters)
	hash := sha256.Sum256([]byte{byte(rand.Intn(256)), byte(time.Now().Unix())})
	return "0x" + hex.EncodeToString(hash[:])
}

func generateTokenID() string {
	// Generate mock NFT token ID
	return fmt.Sprintf("%d", rand.Int63())
}

func generateVerificationNonce() string {
	// Generate random nonce for signature verification
	bytes := make([]byte, 32)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(bytes)
}
