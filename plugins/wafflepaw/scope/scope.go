package scope

import (
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/infinitybotlist/sysmanage-web/plugins/wafflepaw/types"
)

type Scope struct {
	// This always links to the current integration
	Integration *types.IntegrationConfig

	// Export mutex
	exportMutex sync.Mutex

	// Variables
	variables map[string]map[string]any

	// Global index for variables
	globalIndex map[string]string
}

func (s *Scope) Export(k string, value any) {
	// Ensure thread safety
	s.exportMutex.Lock()
	defer s.exportMutex.Unlock()

	// Check if exportable
	if !slices.Contains(s.Integration.Export, k) {
		return
	}

	if s.variables == nil {
		s.variables = make(map[string]map[string]any)
	}

	if s.globalIndex == nil {
		s.globalIndex = make(map[string]string)
	}

	s.variables[s.Integration.Name][k] = value

	_, ok := s.globalIndex[k]

	// Only if its not already in the global index can we add it to global index
	// for faster access
	//
	// Note that the global index is not actually used right now but may be used in the future
	if !ok {
		s.globalIndex[k] = s.Integration.Name
	}
}

// Returns an exported value from the scope
//
// Note that all keys must be of the format <integration>.<key> (e.g. ping.response_time)
func (s *Scope) GetExported(k string) (Export, error) {
	splitKey := strings.SplitN(k, ".", 2)

	if len(splitKey) != 2 {
		return Export{}, errors.New("invalid key format")
	}

	integration := splitKey[0]
	key := splitKey[1]

	if s.variables == nil {
		return Export{}, errors.New("no variables")
	}

	_, ok := s.variables[integration]

	if !ok {
		return Export{}, errors.New("integration not found")
	}

	_, ok = s.variables[integration][key]

	if !ok {
		return Export{}, errors.New("key not found")
	}

	return Export{
		value: s.variables[integration][key],
	}, nil
}

type Export struct {
	value any
}

func CastExport[T any](e Export, def T) T {
	casted, ok := e.value.(T)

	if !ok {
		return def
	}

	return casted
}

func (e Export) String() string {
	s, ok := e.value.(string)

	if !ok {
		return ""
	}

	return s
}

func (e Export) Int() int {
	i, ok := e.value.(int)

	if !ok {
		return 0
	}

	return i
}

func (e Export) Int64() int64 {
	i, ok := e.value.(int64)

	if !ok {
		return 0
	}

	return i
}

func (e Export) Float64() float64 {
	f, ok := e.value.(float64)

	if !ok {
		return 0
	}

	return f
}
