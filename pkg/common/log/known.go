package log

const (
	XRequestIDKey = "X-Request-ID"

	XIdKey = "X-ID"

	XUidKey = "X-Uid"

	XUsernameKey = "X-Username"

	XEmailKey = "X-Email"
)

// Kafka topic
const (
	OfflinePrivateTopic = "offline_private_message"
	OfflineGroupTopic   = "offline_group_message"

	ChannelOfflineTopic      = "channel-offline-private"
	ChannelGroupOfflineTopic = "channel-offline-group"
	ChannelNodeTopic         = "channel-node"
)

// Redis Bitmap
const (
	RedisBitmapUserLoggedKey = "bitmap:user:online"
)
