package usecase

import (
	"fmt"
	"math"

	"pos-saas/internal/domain"
)

const maxJournalAmount = 9_999_999_999_999.99
const maxJournalAmountCents int64 = 999_999_999_999_999

// AutoPostingValidationError reports a journal that cannot safely be posted.
type AutoPostingValidationError struct {
	Code    string
	Message string
}

func (err *AutoPostingValidationError) Error() string {
	return err.Message
}

func validateAutoPostingLines(lines []domain.JournalEntryLine) ([]domain.JournalEntryLine, error) {
	if len(lines) < 2 {
		return nil, &AutoPostingValidationError{
			Code:    "insufficient_journal_lines",
			Message: "auto-posting requires at least two valid journal lines",
		}
	}

	validated := make([]domain.JournalEntryLine, len(lines))
	var totalDebitCents, totalCreditCents int64
	for index, line := range lines {
		debitCents, err := journalAmountToCents(line.DebitAmount, index, "debit")
		if err != nil {
			return nil, err
		}
		creditCents, err := journalAmountToCents(line.CreditAmount, index, "credit")
		if err != nil {
			return nil, err
		}
		if debitCents > 0 && creditCents > 0 {
			return nil, invalidJournalLine(index, "a line cannot contain both a debit and a credit amount")
		}
		if debitCents == 0 && creditCents == 0 {
			return nil, invalidJournalLine(index, "a line must contain a positive debit or credit amount")
		}
		if totalDebitCents > maxJournalAmountCents-debitCents || totalCreditCents > maxJournalAmountCents-creditCents {
			return nil, &AutoPostingValidationError{
				Code:    "journal_total_out_of_range",
				Message: "auto-posting totals exceed the journal database amount limit",
			}
		}

		totalDebitCents += debitCents
		totalCreditCents += creditCents
		lineCopy := line
		if lineCopy.LineOrder == 0 {
			lineCopy.LineOrder = index + 1
		}
		validated[index] = lineCopy
	}

	if totalDebitCents != totalCreditCents {
		debitTotal := float64(totalDebitCents) / 100
		creditTotal := float64(totalCreditCents) / 100
		difference := math.Abs(debitTotal - creditTotal)
		return nil, &AutoPostingValidationError{
			Code: "unbalanced_journal",
			Message: fmt.Sprintf(
				"auto-posting rejected: debit total %.2f and credit total %.2f differ by %.2f; correct the source amounts before retrying",
				debitTotal,
				creditTotal,
				difference,
			),
		}
	}

	return validated, nil
}

func journalAmountToCents(amount float64, lineIndex int, side string) (int64, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, invalidJournalLine(lineIndex, fmt.Sprintf("%s amount must be finite", side))
	}
	if amount < 0 {
		return 0, invalidJournalLine(lineIndex, fmt.Sprintf("%s amount cannot be negative", side))
	}
	if amount > maxJournalAmount {
		return 0, invalidJournalLine(lineIndex, fmt.Sprintf("%s amount exceeds the journal database amount limit", side))
	}
	if roundedAmount(amount) != amount {
		return 0, invalidJournalLine(lineIndex, fmt.Sprintf("%s amount has more than two decimal places supported by the journal schema", side))
	}
	return int64(math.Round(amount * 100)), nil
}

func invalidJournalLine(lineIndex int, reason string) *AutoPostingValidationError {
	return &AutoPostingValidationError{
		Code:    "invalid_journal_line",
		Message: fmt.Sprintf("auto-posting rejected journal line %d: %s", lineIndex+1, reason),
	}
}
