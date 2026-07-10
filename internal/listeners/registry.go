package listeners

import (
	"vessel.dev/vessel/internal/dispatch"
	"vessel.dev/vessel/internal/events"
)

// RegisterAll binds all domain events to their respective handler functions
func RegisterAll(bus *events.EventBus, dispatcher *dispatch.DispatcherService) {
	registerDeploymentListeners(bus, dispatcher)
	// Additional listener groups can be registered here
}

func registerDeploymentListeners(bus *events.EventBus, dispatcher *dispatch.DispatcherService) {
	bus.Subscribe(events.ProjectDeployed, func(payload interface{}) {
		// Example: Convert payload and use dispatcher to send notifications
		// e.g. dispatcher.Send(...)
	})
}
