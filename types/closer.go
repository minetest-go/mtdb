package types

// Closer is a type that signals the implementation that we are finished
// consuming the data/performing work and the async operation can be stopped.
//
// Used by the block.BlockRepository.Iterator implementation in order to perform
// an early stop.
type Closer interface {
	Close() error
}

// WhenDone is a simple channel-based implementation of the Closer interface
type WhenDone chan bool

func (w WhenDone) Close() error {
	w <- true
	return nil
}
