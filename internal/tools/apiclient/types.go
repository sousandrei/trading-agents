package apiclient

import "context"

type Fetch func(ctx context.Context, args map[string]any) (map[string]any, error)
