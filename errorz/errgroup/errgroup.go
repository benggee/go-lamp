package errgroup

import (
	"context"
	"sync"
)

type Group struct {
	wg        sync.WaitGroup
	cancel    func()
	groupOnce sync.Once
	err       error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	context, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, context
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.err = err
			if g.cancel != nil {
				g.cancel()
			}
		}
	}()
}