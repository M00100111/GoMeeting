package structs

type MemberStatus struct {
	MemberId     uint64
	Status       uint64 //0正常1网络不佳2掉线3禁言
	MicStatus    uint64 //0关1开
	CameraStatus uint64
	ScreenStatus uint64
}
