package typings

type User struct {
	Id     int
	Name   string
	ApiKey string
}

type Feed struct {
	Id     int
	Name   string
	Url    string
	UserId int
}

type FeedFollow struct {
	Id     int
	FeedId int
}
