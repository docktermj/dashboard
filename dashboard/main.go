package dashboard

import (
	"context"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Dashboard interface {
	Serve(ctx context.Context) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6204xxxx".
const ComponentId = 6204

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates.
var IdMessages = map[int]string{
	2000: "Entry: %+v",
	2001: "SENZING_ENGINE_CONFIGURATION_JSON: %v",
	2002: "Enabling all services",
	2003: "Server listening on port %v",
	3001: "Error reading file: %s",
	4001: "Call to net.Listen(tcp, %s) failed.",
	5001: "Failed to serve.",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
