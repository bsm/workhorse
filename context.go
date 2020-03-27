package workhorse

import "context"

type contextKey struct{ name string }

func (k *contextKey) String() string {
	return "workhorse context value " + k.name
}

var taskNameCtxKey = &contextKey{"TaskName"}

// TaskName extracts the task name from the context.
func TaskName(ctx context.Context) string {
	if v := ctx.Value(taskNameCtxKey); v != nil {
		return v.(string)
	}
	return ""
}
