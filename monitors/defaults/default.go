/*
Package defaults imports all Monitor packages so that they
register with the global monitor registry. This package can be imported in the
main package to automatically register all of the standard supported Heartbeat
modules.
*/
package defaults

import (
	_ "github.com/lijingwei9060/infobeat/monitors/active/hw"
)
