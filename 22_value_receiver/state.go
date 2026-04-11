package value_receiver

// UpdateState tracks the active firmware version and maintains a history
// of previous versions for rollback support.
//
// Apply transitions to a new version.
// Rollback reverts to the most recent previous version.
// Current returns the active version.
// History returns all previous versions, oldest first.

type UpdateState struct {
    current string
    history []string
}

func NewUpdateState(initial string) *UpdateState {
    return &UpdateState{current: initial}
}

// Apply records the current version in history and sets newVersion as active.
func (s *UpdateState) Apply(newVersion string) {
    s.history = append(s.history, s.current)
    s.current = newVersion
}

func (s *UpdateState) Current() string {
    return s.current
}

func (s *UpdateState) History() []string {
    return s.history
}

// Rollback reverts to the most recent historical version.
// Does nothing if history is empty.
func (s *UpdateState) Rollback() {
    if len(s.history) == 0 {
        return
    }
    last := len(s.history) - 1
    s.current = s.history[last]
    s.history = s.history[:last]
}
