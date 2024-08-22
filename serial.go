package tgk

type SerialHIDMoc struct {
}

func NewSerialHIDMoc() *SerialHIDMoc {
	return &SerialHIDMoc{}
}

func (b *SerialHIDMoc) Begin() error {
	println("tgk: begin")
	return nil
}
func (b *SerialHIDMoc) End() error {
	println("tgk: end")
	return nil
}
func (b *SerialHIDMoc) Up(keycode uint16) error {
	println("[UP] keycode: %d", keycode)
	return nil
}
func (b *SerialHIDMoc) Down(keycode uint16) error {
	println("[DOWN] keycode: %d", keycode)
	return nil
}
func (b *SerialHIDMoc) sendReport(report []byte) error {
	println("[SEND REPORT] report: %v", report)
	return nil
}
func (b *SerialHIDMoc) ReleaseAll() error {
	println("[RELEASE ALL]")
	return nil
}
