// +build linux

package filewatch

import "context"

// RunEvent is just a facarde
func (e *Handler) RunEvent(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {
	e.Run(ctx, exit, exitError)
}
