package structs

type MemberStatus struct {
	UserId   string
	Username string
	Sex      uint64
	Email    string

	UserStatus uint64
	UserType   uint64

	MicStatus    uint64 //0关1开
	CameraStatus uint64
	ScreenStatus uint64
}
