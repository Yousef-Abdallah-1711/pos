package repository

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentProductEditsCannotUndoConfirmedOrderStock(t *testing.T) {
	fixture := newStockFixture(t, 20)
	products := NewProductRepository(fixture.db)
	snapshot, err := products.GetProductByID(fixture.tenant, fixture.restaurant, fixture.product)
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]int64, 10)
	for index := range ids {
		ids[index] = fixture.order(t, fmt.Sprintf("concurrent-metadata-%d", index), 2)
	}
	start := make(chan struct{})
	results := make(chan error, 20)
	var workers sync.WaitGroup
	for index, id := range ids {
		workers.Add(2)
		go func(orderID int64) {
			defer workers.Done()
			<-start
			results <- fixture.transition(orderID, "confirmed")
		}(id)
		go func(index int) {
			defer workers.Done()
			copy := *snapshot
			copy.NameEn = fmt.Sprintf("Concurrent metadata %d", index)
			<-start
			_, err := products.UpdateProduct(&copy)
			results <- err
		}(index)
	}
	close(start)
	workers.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatalf("concurrent order/metadata operation: %v", err)
		}
	}
	fixture.assertStock(t, 0)
}
