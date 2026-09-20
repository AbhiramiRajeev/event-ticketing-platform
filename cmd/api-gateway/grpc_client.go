package gateway

import (
	eventpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
	registrationpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
	userpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
)

type Clients struct {
	UserClient         userpb.UserServiceClient
	EventClient        eventpb.EventServiceClient
	RegistrationClient registrationpb.RegistrationServiceClient
}