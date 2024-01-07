package api

import (
	"sync"

	"go.uber.org/dig"
)

var (
	binder *dig.Container
	once   sync.Once
)

func New() *dig.Container {
	once.Do(func() {
		binder = dig.New()

		provideApp(binder)
		provideController(binder)
		provideUseCase(binder)
		provideRepository(binder)
		provideWrapper(binder)
	})

	return binder
}
