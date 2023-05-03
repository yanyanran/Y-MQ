package yerMQ

type Stats struct {
	Topics    []TopicStats
	Producers []ClientStats
}

type TopicStats struct {
	TopicName string
	Channels  []ChannelStats
	Stopped   bool
}

type ChannelStats struct {
	ChannelName string
	Clients     []ClientStats
	Stopped     bool
}

type ClientStats interface {
	String() string
}
