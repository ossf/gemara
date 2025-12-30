// Schema lifecycle: experimental | stable | deprecated
@status("stable")
package schemas

import "time"

@go(gemara)

// Contact represents contact information used across multiple layers
#Contact: {
	// The contact person's name.
	name: string
	// The entity with which the contact is affiliated, such as a school or employer.
	affiliation?: string @go(Affiliation,type=*string)
	// A preferred email address to reach the contact.
	email?: #Email @go(Email,type=*Email)
	// A social media handle or profile for the contact.
	social?: string @go(Social,type=*string)
}

// Actor represents an entity (human or tool) that can perform actions in evaluations.
#Actor: {
	// Id uniquely identifies the actor.
	id: string
	// Name provides the name of the actor.
	name: string
	// Type specifies the type of entity interacting in the workflow.
	type: #ActorType @go(Type)
	// Version specifies the version of the actor (if applicable, e.g., for tools).
	version?: string
	// Description provides additional context about the actor.
	description?: string
	// Uri provides a general URI for the actor information.
	uri?: =~"^https?://[^\\s]+$"
	// Contact provides contact information for the actor.
	contact?: #Contact @go(Contact)
}

// ActorType specifies what entity is interacting in the workflow.
#ActorType: "Human" | "Software" | "Software-Assisted" @go(-)

// Email represents a validated email address pattern
#Email: =~"^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"

// Datetime represents an ISO 8601 formatted datetime string
#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime,format="date-time")

// Date represents a date string (ISO 8601 date format)
#Date: time.Format("2006-01-02") @go(Date,format="date")

// Category represents a category used for applicability or classification
#Category: {
	id:          string
	title:       string
	description: string
}

// Family represents a logical grouping of guidelines or controls which share a common purpose or function
#Family: {
	id:          string
	title:       string
	description: string
}
