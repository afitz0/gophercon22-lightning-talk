package constants

const (
	DEBUG = true

	// "Happy path" statuses
	ORDER_RECIEVED   = "received"
	ORDER_PLACED     = "placed"
	ORDER_INPROGRESS = "fulfilling"
	ORDER_FULFILLED  = "fulfilled"
	ORDER_ARCHIVED   = "archived"

	// Error-state statuses
	E_ORDER_DROPPED   = "dropped"
	E_ORDER_LOST      = "lost"
	E_ORDER_BACKORDER = "backordered"
	E_ORDER_DUPLICATE = "duped"
	E_TIMEOUT         = "timeout"
)
