package webhook

import (
	"context"
	"sync"
)

// Router manages event routing to handlers
type Router struct {
	routes map[string]map[string][]Handler // provider -> eventType -> handlers
	mu     sync.RWMutex
}

// NewRouter creates a new webhook router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string][]Handler),
	}
}

// Register registers a handler for a provider and event type
func (r *Router) Register(provider string, eventType string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.routes[provider] == nil {
		r.routes[provider] = make(map[string][]Handler)
	}

	r.routes[provider][eventType] = append(r.routes[provider][eventType], handler)
}

// Route routes an event to registered handlers
func (r *Router) Route(ctx context.Context, event *Event) error {
	r.mu.RLock()
	providerRoutes, ok := r.routes[event.Provider]
	r.mu.RUnlock()

	if !ok {
		return nil // No routes for this provider
	}

	handlers, ok := providerRoutes[event.Type]
	if !ok {
		return nil // No handlers for this event type
	}

	// Execute handlers
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// AsyncRouter executes handlers asynchronously
type AsyncRouter struct {
	router     *Router
	maxWorkers int
	queue      chan *routeJob
	wg         sync.WaitGroup
}

type routeJob struct {
	ctx   context.Context
	event *Event
	done  chan error
}

// NewAsyncRouter creates an async webhook router
func NewAsyncRouter(router *Router, maxWorkers int) *AsyncRouter {
	ar := &AsyncRouter{
		router:     router,
		maxWorkers: maxWorkers,
		queue:      make(chan *routeJob, maxWorkers*2),
	}

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		ar.wg.Add(1)
		go ar.worker()
	}

	return ar
}

// worker processes routing jobs
func (ar *AsyncRouter) worker() {
	defer ar.wg.Done()

	for job := range ar.queue {
		err := ar.router.Route(job.ctx, job.event)
		job.done <- err
		close(job.done)
	}
}

// Route routes an event asynchronously
func (ar *AsyncRouter) Route(ctx context.Context, event *Event) error {
	job := &routeJob{
		ctx:   ctx,
		event: event,
		done:  make(chan error, 1),
	}

	ar.queue <- job
	return <-job.done
}

// Stop stops the async router
func (ar *AsyncRouter) Stop() {
	close(ar.queue)
	ar.wg.Wait()
}

// DeadLetterQueue handles failed webhook deliveries
type DeadLetterQueue struct {
	failed []*FailedEvent
	mu     sync.RWMutex
}

// FailedEvent represents a failed webhook event
type FailedEvent struct {
	Event         *Event
	Error         error
	Attempts      int
	FirstFailedAt string
	LastFailedAt  string
}

// NewDeadLetterQueue creates a new dead letter queue
func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		failed: make([]*FailedEvent, 0),
	}
}

// Add adds a failed event to the queue
func (dlq *DeadLetterQueue) Add(event *Event, err error) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	dlq.failed = append(dlq.failed, &FailedEvent{
		Event:    event,
		Error:    err,
		Attempts: 1,
	})
}

// GetAll returns all failed events
func (dlq *DeadLetterQueue) GetAll() []*FailedEvent {
	dlq.mu.RLock()
	defer dlq.mu.RUnlock()

	result := make([]*FailedEvent, len(dlq.failed))
	copy(result, dlq.failed)
	return result
}

// Remove removes a failed event from the queue
func (dlq *DeadLetterQueue) Remove(eventID string) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	for i, failed := range dlq.failed {
		if failed.Event.ID == eventID {
			dlq.failed = append(dlq.failed[:i], dlq.failed[i+1:]...)
			return
		}
	}
}
