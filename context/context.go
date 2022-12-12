// context/context.go
// Draw contexts.

package context

import "time"

type DrawContext interface {
	PreDraw(time.Duration)
}
