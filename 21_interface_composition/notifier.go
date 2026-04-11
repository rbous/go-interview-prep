package interface_composition

// A notification system built with interface composition.
//
// Formatter formats a raw message into a final string.
// Sender delivers messages and records what was sent.
// PrefixNotifier composes both: it formats then sends.

// Formatter formats a raw message.
type Formatter interface {
    Format(msg string) string
}

// Sender sends a message and tracks sent history.
type Sender interface {
    Send(msg string) error
    Sent() []string
}

// PrefixFormatter prepends a prefix to every message.
type PrefixFormatter struct {
    Prefix string
}

func (f *PrefixFormatter) Format(msg string) string {
    return f.Prefix + ": " + msg
}

// MemorySender records all sent messages in memory.
type MemorySender struct {
    messages []string
}

func (s *MemorySender) Send(msg string) error {
    s.messages = append(s.messages, msg)
    return nil
}

func (s *MemorySender) Sent() []string {
    return s.messages
}

// PrefixNotifier composes a PrefixFormatter and a MemorySender.
// Notify should format the message and send it.
type PrefixNotifier struct {
    *PrefixFormatter
    *MemorySender
}

// NewPrefixNotifier returns a PrefixNotifier that prepends prefix to all messages.
func NewPrefixNotifier(prefix string) *PrefixNotifier {
    return &PrefixNotifier{
        PrefixFormatter: &PrefixFormatter{prefix},
        MemorySender: &MemorySender{},
    }
}

// Notify formats msg and delivers it via the embedded sender.
func (n *PrefixNotifier) Notify(msg string) error {
    formatted := n.Format(msg)
    return n.Send(formatted)
}
