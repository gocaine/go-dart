package hardware

// An InputEvent holds board input event.
type InputEvent struct {
	Sector     int
	Multiplier int
}

// InputEventProducer is waiting for hardware event and push then in the chan.
type InputEventProducer interface {
	Produce(inputEventChannel chan InputEvent)
	Shutdown()
}

// InputEventConsumer can consume InputEvent
type InputEventConsumer interface {
	// Prepare() error
	Consume(event InputEvent)
	Shutdown()
}
