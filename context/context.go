// context/context.go
// Draw contexts.

package context

type DrawContext interface {
	PreDraw()
	PostDraw()
}
