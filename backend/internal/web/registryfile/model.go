package registryfile

import (
	"context"
	"encoding/json"
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/web/balancer"
	"logbook/internal/web/logger"
	"logbook/models"
	"os"
	"sync"
	"time"
)

type FileReader struct {
	l        *logger.Logger
	filepath string
	params   ServiceParams
	reload   time.Duration
	content  []models.Instance
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

type ServiceParams struct {
	Port int
	Tls  bool
}

// interface
var _ balancer.InstanceSource = &FileReader{}

func NewFileReader(filepath string, deplycfg *deployment.Config, params ServiceParams) *FileReader {
	ctx, cancel := context.WithCancel(context.Background())
	fr := &FileReader{
		l:        logger.NewLogger("FileReader"),
		filepath: filepath,
		ctx:      ctx,
		cancel:   cancel,
		reload:   deplycfg.RegistryFile.UpdatePeriod,
		params:   params,
	}
	fr.readfile()
	go fr.tick()
	return fr
}

func (r *FileReader) readfile() error {
	h, err := os.Open(r.filepath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	err = json.NewDecoder(h).Decode(&(r.content))
	if err != nil {
		return fmt.Errorf("decoding file: %w", err)
	}
	for i := 0; i < len(r.content); i++ { // overwrite
		r.content[i].Port = r.params.Port
		r.content[i].Tls = r.params.Tls
	}
	return nil
}

func (r *FileReader) Stop() {
	r.cancel()
}

func (r *FileReader) tick() {
	t := time.NewTicker(r.reload)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			r.mu.Lock()
			if err := r.readfile(); err != nil {
				r.l.Println(fmt.Errorf("error: reloading discovery file: %w", err))
			}
			r.mu.Unlock()
		case <-r.ctx.Done():
			return
		}
	}
}

func (df *FileReader) Instances() ([]models.Instance, error) {
	df.mu.RLock()
	defer df.mu.RUnlock()
	return df.content, nil
}
