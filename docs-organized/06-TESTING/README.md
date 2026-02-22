# Testing Documentation

## Overview

Testing strategy and execution reports for the POS SaaS platform.

## Documents

- **[PHASE_1_5_INTEGRATION_TESTS.md](PHASE_1_5_INTEGRATION_TESTS.md)** - Integration test suite covering phases 1-5
- **[PERFORMANCE_COMPARISON_ANALYSIS.md](PERFORMANCE_COMPARISON_ANALYSIS.md)** - Performance benchmarks and comparison analysis

## Running Tests

### Backend
```bash
cd backend
go test ./...
go vet ./...
```

### Frontend
```bash
cd frontend
pnpm test          # Run all tests
pnpm build         # Verify builds
```
