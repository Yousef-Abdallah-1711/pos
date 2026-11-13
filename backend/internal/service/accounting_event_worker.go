package service

import (
	"context"
	"log"
	"time"
)

type AccountingEventProcessor interface {
	ProcessNextAccountingEvent(context.Context) (bool, error)
	ReconcileAccountingEvents(context.Context) error
}

// AccountingEventWorker replays durable accounting intents until the process stops.
type AccountingEventWorker struct {
	processor AccountingEventProcessor
	interval  time.Duration
}

func NewAccountingEventWorker(processor AccountingEventProcessor) *AccountingEventWorker {
	return &AccountingEventWorker{processor: processor, interval: 15 * time.Second}
}

func (w *AccountingEventWorker) Start(ctx context.Context) {
	if w.processor == nil {
		log.Printf("[AccountingRecovery] worker not started: processor is missing")
		return
	}
	go func() {
		w.runOnce(ctx)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[AccountingRecovery] worker stopped")
				return
			case <-ticker.C:
				w.runOnce(ctx)
			}
		}
	}()
}

func (w *AccountingEventWorker) runOnce(ctx context.Context) {
	if err := w.processor.ReconcileAccountingEvents(ctx); err != nil {
		log.Printf("[AccountingRecovery] reconciliation error: %v", err)
		return
	}
	for ctx.Err() == nil {
		processed, err := w.processor.ProcessNextAccountingEvent(ctx)
		if err != nil {
			log.Printf("[AccountingRecovery] processing error: %v", err)
		}
		if !processed || err != nil {
			return
		}
	}
}
