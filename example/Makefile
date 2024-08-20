.PHONY: build flash flash-monitor clean
build: 
	mkdir -p build
	tinygo build -target=feather-nrf52840 --size short -o build/firmware.uf2 .
flash:
	tinygo flash -target=feather-nrf52840 --size short . 
flash-monitor:
	tinygo flash -target=feather-nrf52840 --size short -monitor
clean:
	rm -rf build