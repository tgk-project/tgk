package tgk

type CommonHID interface {
	Begin() error
	End() error
	Up(keycode uint16) error
	Down(keycode uint16) error
	sendReport(report []byte) error
	ReleaseAll() error
}
