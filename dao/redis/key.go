package redis

const (
	Prefix             = "webapp:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	KeyPostVotedZSetPF = "post:voted"
	KeyCommunitySetPF  = "community:"
)

func GetRedisKey(key string) string {
	return Prefix + key

}
