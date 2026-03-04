package usecase

import (
	"errors"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

type AccountingUseCase struct {
	accountRepo *repository.AccountRepository
}

func NewAccountingUseCase(accountRepo *repository.AccountRepository) *AccountingUseCase {
	return &AccountingUseCase{accountRepo: accountRepo}
}

func (uc *AccountingUseCase) ListAccounts(tenantID, restaurantID int, tree bool, includeInactive bool) ([]domain.Account, error) {
	if tree {
		return uc.accountRepo.GetAccountTree(tenantID, restaurantID, includeInactive)
	}
	return uc.accountRepo.ListAccounts(tenantID, restaurantID, includeInactive)
}

func (uc *AccountingUseCase) GetAccountByID(tenantID, restaurantID, id int) (*domain.Account, error) {
	return uc.accountRepo.GetAccountByID(tenantID, restaurantID, id)
}

func (uc *AccountingUseCase) CreateAccount(account *domain.Account) (int, error) {
	if err := uc.prepareAccount(account, 0); err != nil {
		return 0, err
	}
	account.IsSystem = false
	return uc.accountRepo.CreateAccount(account)
}

func (uc *AccountingUseCase) UpdateAccount(account *domain.Account) error {
	existing, err := uc.accountRepo.GetAccountByID(account.TenantID, account.RestaurantID, account.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("account not found")
	}

	account.IsSystem = existing.IsSystem
	account.CurrentBalance = existing.CurrentBalance
	if err := uc.prepareAccount(account, account.ID); err != nil {
		return err
	}

	return uc.accountRepo.UpdateAccount(account)
}

func (uc *AccountingUseCase) DeleteAccount(tenantID, restaurantID, id int) error {
	account, err := uc.accountRepo.GetAccountByID(tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("account not found")
	}
	if account.IsSystem {
		return errors.New("system accounts cannot be deleted")
	}

	hasChildren, err := uc.accountRepo.HasChildAccounts(tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("account cannot be deleted while it has child accounts")
	}

	hasEntries, err := uc.accountRepo.HasJournalEntries(tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if hasEntries {
		return errors.New("account cannot be deleted because it has journal entries")
	}

	return uc.accountRepo.DeleteAccount(tenantID, restaurantID, id)
}

func (uc *AccountingUseCase) prepareAccount(account *domain.Account, currentID int) error {
	account.AccountCode = strings.TrimSpace(account.AccountCode)
	account.AccountName = strings.TrimSpace(account.AccountName)
	account.AccountNameAr = strings.TrimSpace(account.AccountNameAr)
	account.Description = strings.TrimSpace(account.Description)
	account.DescriptionAr = strings.TrimSpace(account.DescriptionAr)
	account.AccountType = strings.ToLower(strings.TrimSpace(account.AccountType))
	account.AccountSubtype = strings.ToLower(strings.TrimSpace(account.AccountSubtype))
	account.NormalBalance = strings.ToLower(strings.TrimSpace(account.NormalBalance))
	account.Currency = strings.ToUpper(strings.TrimSpace(account.Currency))

	if account.AccountCode == "" {
		return errors.New("account_code is required")
	}
	if account.AccountName == "" {
		return errors.New("account_name is required")
	}
	if account.AccountType == "" {
		return errors.New("account_type is required")
	}
	if !isValidAccountType(account.AccountType) {
		return errors.New("account_type must be one of: asset, liability, equity, revenue, expense")
	}
	if account.NormalBalance == "" {
		account.NormalBalance = defaultNormalBalance(account.AccountType)
	}
	if account.NormalBalance != domain.NormalBalanceDebit && account.NormalBalance != domain.NormalBalanceCredit {
		return errors.New("normal_balance must be debit or credit")
	}
	if account.Currency == "" {
		account.Currency = "EGP"
	}

	account.Level = 1
	if account.ParentID != nil {
		if *account.ParentID == currentID && currentID > 0 {
			return errors.New("account cannot be its own parent")
		}

		parent, err := uc.accountRepo.GetAccountByID(account.TenantID, account.RestaurantID, *account.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent account not found")
		}
		account.Level = parent.Level + 1
	}

	return nil
}

func isValidAccountType(accountType string) bool {
	switch accountType {
	case domain.AccountTypeAsset, domain.AccountTypeLiability, domain.AccountTypeEquity, domain.AccountTypeRevenue, domain.AccountTypeExpense:
		return true
	default:
		return false
	}
}

func defaultNormalBalance(accountType string) string {
	switch accountType {
	case domain.AccountTypeLiability, domain.AccountTypeEquity, domain.AccountTypeRevenue:
		return domain.NormalBalanceCredit
	default:
		return domain.NormalBalanceDebit
	}
}
