package pkg

import (
	"context"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

func (o *Orchestrator[K, V]) defaultSaveOptions() protocols.SaveOptions {
	ctx := context.Background()

	options := protocols.SaveOptions{
		Context:       ctx,
		HowWillItSave: protocols.Sequential,
		Targets:       o.standardOrder,
	}
	return options
}

func (o *Orchestrator[K, V]) defaultGetOptions() protocols.GetOptions {
	ctx := context.Background()
	options := protocols.GetOptions{
		Context:      ctx,
		HowWillItGet: protocols.Cache,
		Targets:      o.standardOrder,
	}
	return options
}

func (o *Orchestrator[K, V]) defaultDeleteOptions() protocols.DeleteOptions {
	ctx := context.Background()
	options := protocols.DeleteOptions{
		Context:         ctx,
		Targets:         o.standardOrder,
		HowWillItDelete: protocols.SequentialDelete,
	}
	return options
}
