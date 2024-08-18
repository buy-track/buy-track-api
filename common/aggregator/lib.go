package aggregator

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"

	"time"
)

type Correlation struct {
	CorrelationId string
	Data          any
}

type RespondCorrelation struct {
	CorrelationId string
	Status        bool
	Data          any
}

type CorrelationManager struct {
	pending map[string]chan *RespondCorrelation
	count   uint64
}

func (c *CorrelationManager) Add(data any) (*Correlation, <-chan *RespondCorrelation) {
	id := uuid.New().String()
	correlation := &Correlation{
		CorrelationId: id,
		Data:          data,
	}
	receiver := make(chan *RespondCorrelation, 1)
	c.pending[id] = receiver
	c.count++

	return correlation, receiver
}

func (c *CorrelationManager) Respond(id string, data *RespondCorrelation) {
	if r, ok := c.pending[id]; ok {
		defer func() {
			close(r)
			delete(c.pending, id)
			c.count--
		}()
		r <- data
	} else {
		log.Errorf("cannot found id : %v", data)
	}
}

func NewCorrelationManager() *CorrelationManager {
	t := time.Now()
	rand.Seed(uint64(t.Unix()))
	return &CorrelationManager{
		pending: make(map[string]chan *RespondCorrelation, 1000),
		count:   0,
	}
}
