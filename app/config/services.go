// Since both the services and their consumers needs to know URL paths, single source-of-truth is needed.
// TODO: Keep this up to date as new services and endpoints defined.

package config

import (
	. "logbook/internal/web/paths"
)

var Site Domain // assigned after config read
var ApiGateway = Gateway{Site, "/api/v1.0.0"}

var (
	Document     = Service{ApiGateway, "/document"}
	DocumentList = Endpoint{Document, "/list/{root}", GET}
)

var (
	Task     = Service{ApiGateway, "/tasks"}
	TaskList = Endpoint{Task, "/list/{root}", GET}
)

var (
	Tag         = Service{ApiGateway, "tags"}
	TagCreation = Endpoint{Tag, "/", POST}
	TagAssign   = Endpoint{Tag, "/assign", POST}
)
