package model

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	UserID       string
	AtExpires    int64
	RtExpires    int64
}
