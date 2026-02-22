//go:build legacy
// +build legacy

package service

import (
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// ABTestingService manages A/B testing experiments
type ABTestingService struct {
	experiments map[string]*ABExperiment
	variants    map[string][]*Variant
	assignments map[string]*UserAssignment
	results     map[string]*ExperimentResults
	mu          sync.RWMutex
}

// ABExperiment represents an A/B test experiment
type ABExperiment struct {
	ID             string
	Name           string
	Description    string
	RestaurantID   int64
	Status         string // "draft", "running", "paused", "completed"
	Objective      string
	HypothesisText string
	Target         string      // "all_users", "percentage", "segment"
	TargetValue    int         // percentage or segment ID
	StartDate      time.Time
	EndDate        *time.Time
	MinSampleSize  int64
	StatisticalSig float64 // 0.95, 0.99, etc.
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Variant represents a test variant
type Variant struct {
	ID             string
	ExperimentID   string
	Name           string
	Description    string
	TrafficPercent int
	IsControl      bool // True for control group
	Properties     map[string]interface{}
	CreatedAt      time.Time
}

// UserAssignment tracks variant assignment for users
type UserAssignment struct {
	ID           string
	ExperimentID string
	UserID       int64
	VariantID    string
	AssignedAt   time.Time
	Attributes   map[string]interface{}
}

// ExperimentResults contains statistical results of an experiment
type ExperimentResults struct {
	ExperimentID     string
	Status           string // "running", "inconclusive", "winner", "loser"
	WinnerVariantID  *string
	Variants         map[string]*VariantStats
	OverallConversion float64
	StatisticalSig   float64
	ConfidenceLevel  float64
	SampleSize       int64
	CalculatedAt     time.Time
}

// VariantStats contains statistics for a variant
type VariantStats struct {
	VariantID        string
	Participants     int64
	Conversions      int64
	ConversionRate   float64
	AverageValue     float64
	StandardDev      float64
	LiftVsControl    float64
	IsWinner         bool
}

// NewABTestingService creates a new A/B testing service
func NewABTestingService() *ABTestingService {
	return &ABTestingService{
		experiments: make(map[string]*ABExperiment),
		variants:    make(map[string][]*Variant),
		assignments: make(map[string]*UserAssignment),
		results:     make(map[string]*ExperimentResults),
	}
}

// CreateExperiment creates a new A/B experiment
func (s *ABTestingService) CreateExperiment(experiment *ABExperiment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if experiment.ID == "" {
		experiment.ID = fmt.Sprintf("exp_%d", time.Now().UnixNano())
	}
	experiment.Status = "draft"
	experiment.CreatedAt = time.Now()
	experiment.UpdatedAt = time.Now()

	s.experiments[experiment.ID] = experiment
	s.variants[experiment.ID] = make([]*Variant, 0)

	fmt.Printf("DEBUG: A/B experiment created: %s\n", experiment.Name)
	return nil
}

// AddVariant adds a variant to an experiment
func (s *ABTestingService) AddVariant(experimentID string, variant *Variant) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.experiments[experimentID]; !exists {
		return fmt.Errorf("experiment not found")
	}

	if variant.ID == "" {
		variant.ID = fmt.Sprintf("var_%d", time.Now().UnixNano())
	}
	variant.ExperimentID = experimentID
	variant.CreatedAt = time.Now()

	s.variants[experimentID] = append(s.variants[experimentID], variant)
	fmt.Printf("DEBUG: Variant added to experiment %s: %s\n", experimentID, variant.Name)

	return nil
}

// StartExperiment starts an experiment
func (s *ABTestingService) StartExperiment(experimentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	experiment, exists := s.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment not found")
	}

	// Validate variants
	variants, exists := s.variants[experimentID]
	if !exists || len(variants) < 2 {
		return fmt.Errorf("experiment must have at least 2 variants")
	}

	experiment.Status = "running"
	experiment.StartDate = time.Now()
	experiment.UpdatedAt = time.Now()

	s.results[experimentID] = &ExperimentResults{
		ExperimentID: experimentID,
		Status:       "running",
		Variants:     make(map[string]*VariantStats),
		CalculatedAt: time.Now(),
	}

	fmt.Printf("DEBUG: Experiment started: %s\n", experimentID)
	return nil
}

// AssignUserToVariant assigns a user to a variant
func (s *ABTestingService) AssignUserToVariant(experimentID string, userID int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	experiment, exists := s.experiments[experimentID]
	if !exists {
		return "", fmt.Errorf("experiment not found")
	}

	if experiment.Status != "running" {
		return "", fmt.Errorf("experiment is not running")
	}

	// Check if user already assigned
	assignmentKey := fmt.Sprintf("%s_%d", experimentID, userID)
	if assignment, exists := s.assignments[assignmentKey]; exists {
		return assignment.VariantID, nil
	}

	variants := s.variants[experimentID]
	if len(variants) == 0 {
		return "", fmt.Errorf("no variants available")
	}

	// Deterministic assignment based on user ID and experiment ID
	variantID := s.selectVariant(userID, experimentID, variants)

	assignment := &UserAssignment{
		ID:           fmt.Sprintf("assign_%d", time.Now().UnixNano()),
		ExperimentID: experimentID,
		UserID:       userID,
		VariantID:    variantID,
		AssignedAt:   time.Now(),
		Attributes:   make(map[string]interface{}),
	}

	s.assignments[assignmentKey] = assignment

	fmt.Printf("DEBUG: User %d assigned to variant %s in experiment %s\n", userID, variantID, experimentID)
	return variantID, nil
}

// RecordConversion records a conversion for an experiment
func (s *ABTestingService) RecordConversion(experimentID string, userID int64, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	assignmentKey := fmt.Sprintf("%s_%d", experimentID, userID)
	assignment, exists := s.assignments[assignmentKey]
	if !exists {
		return fmt.Errorf("user not assigned to experiment")
	}

	results, exists := s.results[experimentID]
	if !exists {
		return fmt.Errorf("experiment results not found")
	}

	variantID := assignment.VariantID
	variantStats, exists := results.Variants[variantID]
	if !exists {
		variantStats = &VariantStats{
			VariantID: variantID,
		}
		results.Variants[variantID] = variantStats
	}

	variantStats.Conversions++
	variantStats.AverageValue = (variantStats.AverageValue*float64(variantStats.Conversions-1) + value) / float64(variantStats.Conversions)

	fmt.Printf("DEBUG: Conversion recorded for user %d in variant %s (value: %.2f)\n", userID, variantID, value)
	return nil
}

// GetAssignment gets the variant assignment for a user
func (s *ABTestingService) GetAssignment(experimentID string, userID int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assignmentKey := fmt.Sprintf("%s_%d", experimentID, userID)
	assignment, exists := s.assignments[assignmentKey]
	if !exists {
		return "", fmt.Errorf("user not assigned")
	}

	return assignment.VariantID, nil
}

// CalculateResults calculates statistical results for an experiment
func (s *ABTestingService) CalculateResults(experimentID string) (*ExperimentResults, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	experiment, exists := s.experiments[experimentID]
	if !exists {
		return nil, fmt.Errorf("experiment not found")
	}

	results, exists := s.results[experimentID]
	if !exists {
		return nil, fmt.Errorf("results not found")
	}

	// Calculate metrics for each variant
	totalParticipants := int64(0)
	totalConversions := int64(0)

	for _, variantStats := range results.Variants {
		totalParticipants += variantStats.Participants
		totalConversions += variantStats.Conversions

		variantStats.ConversionRate = float64(variantStats.Conversions) / float64(variantStats.Participants)
	}

	results.OverallConversion = float64(totalConversions) / float64(totalParticipants)
	results.SampleSize = totalParticipants

	// Find control variant
	variants := s.variants[experimentID]
	var controlVariantID string
	for _, v := range variants {
		if v.IsControl {
			controlVariantID = v.ID
			break
		}
	}

	// Calculate lift vs control
	if controlStats, exists := results.Variants[controlVariantID]; exists && controlStats.Participants > 0 {
		for _, variantStats := range results.Variants {
			if variantStats.VariantID != controlVariantID {
				lift := (variantStats.ConversionRate - controlStats.ConversionRate) / controlStats.ConversionRate
				variantStats.LiftVsControl = lift * 100
			}
		}
	}

	// Determine statistical significance using chi-square test
	if totalParticipants >= experiment.MinSampleSize {
		s.performChiSquareTest(results, experiment.StatisticalSig)
	}

	results.CalculatedAt = time.Now()

	fmt.Printf("DEBUG: Results calculated for experiment %s\n", experimentID)
	return results, nil
}

// GetResults gets the current results of an experiment
func (s *ABTestingService) GetResults(experimentID string) (*ExperimentResults, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results, exists := s.results[experimentID]
	if !exists {
		return nil, fmt.Errorf("results not found")
	}

	return results, nil
}

// EndExperiment ends an experiment
func (s *ABTestingService) EndExperiment(experimentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	experiment, exists := s.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment not found")
	}

	now := time.Now()
	experiment.EndDate = &now
	experiment.Status = "completed"
	experiment.UpdatedAt = now

	fmt.Printf("DEBUG: Experiment ended: %s\n", experimentID)
	return nil
}

// ListExperiments lists all experiments
func (s *ABTestingService) ListExperiments(restaurantID int64) []*ABExperiment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	experiments := make([]*ABExperiment, 0)
	for _, exp := range s.experiments {
		if exp.RestaurantID == restaurantID {
			experiments = append(experiments, exp)
		}
	}

	return experiments
}

// Helper functions

func (s *ABTestingService) selectVariant(userID int64, experimentID string, variants []*Variant) string {
	// Create deterministic hash based on user ID and experiment ID
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%s", userID, experimentID)))
	hashValue := int64(0)
	for i := 0; i < 8; i++ {
		hashValue = hashValue*256 + int64(hash[i])
	}

	// Use hash to select variant based on traffic percentages
	randValue := hashValue % 100
	cumulativePercent := 0

	for _, variant := range variants {
		cumulativePercent += variant.TrafficPercent
		if randValue < int64(cumulativePercent) {
			return variant.ID
		}
	}

	return variants[0].ID
}

func (s *ABTestingService) performChiSquareTest(results *ExperimentResults, targetSigLevel float64) {
	if len(results.Variants) < 2 {
		return
	}

	variants := make([]*VariantStats, 0, len(results.Variants))
	for _, v := range results.Variants {
		variants = append(variants, v)
	}

	// Calculate chi-square statistic
	expectedConversionRate := results.OverallConversion
	chiSquare := 0.0

	for _, variantStats := range variants {
		expected := float64(variantStats.Participants) * expectedConversionRate
		observed := float64(variantStats.Conversions)

		if expected > 0 {
			chiSquare += math.Pow(observed-expected, 2) / expected
		}
	}

	// Simple approximation: if chi-square > 3.84, significant at 95% confidence
	// This is a simplified approach; production would use proper statistical tests
	significanceThreshold := 3.84
	if chiSquare > significanceThreshold {
		results.Status = "winner"
		results.ConfidenceLevel = 0.95

		// Find winner (highest conversion rate)
		var winner *VariantStats
		for _, v := range variants {
			if winner == nil || v.ConversionRate > winner.ConversionRate {
				winner = v
			}
		}

		if winner != nil {
			winner.IsWinner = true
			winnerID := winner.VariantID
			results.WinnerVariantID = &winnerID
		}
	} else {
		results.Status = "inconclusive"
		results.ConfidenceLevel = 0.0
	}

	fmt.Printf("DEBUG: Chi-square test completed (value: %.2f, status: %s)\n", chiSquare, results.Status)
}
