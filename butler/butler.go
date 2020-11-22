package butler

import (
	"context"
	"sync"
)

type (
	butler struct {
		tasks []task
	}

	task interface {
		Action(ctx context.Context) error
		Notify(ctx context.Context) error
		Rest(ctx context.Context) error
	}
)

func CallButler() *butler {
	return &butler{}
}

func (b *butler) AddTask(ctx context.Context, t task) {
	b.tasks = append(b.tasks, t)
}

func (b *butler) StartWorking(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, t := range b.tasks {
		t := t
		wg.Add(1)
		go func() {
			for {
				// TODO: error handling
				if err := t.Action(ctx); err != nil {
					wg.Done()
					break
				}
				if err := t.Notify(ctx); err != nil {
					wg.Done()
					break
				}
				if err := t.Rest(ctx); err != nil {
					wg.Done()
					break
				}
			}
		}()
	}
	wg.Wait()
	return nil
}
