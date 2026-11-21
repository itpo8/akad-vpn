package desktop

import (
	"bytes"
	"image/png"
	"testing"
)

// Every tray state has an icon in both formats. AkadVPN's supplied brand
// artwork uses a solid navy tile, rather than the upstream transparent tile
// and its center-aligned upright bar.
func TestTrayIcons(t *testing.T) {
	for _, name := range []string{iconRunning, iconRunningAsleep, iconSet, iconSetAsleep, iconWaiting, iconError, iconOff} {
		data, err := trayFiles.ReadFile("tray/" + name + ".png")
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if b := img.Bounds(); b.Dx() != 32 || b.Dy() != 32 {
			t.Fatalf("%s: size %v", name, b)
		}
		// An icon must contain visible contrasting artwork, regardless of
		// which part of the user's mark occupies its center pixel.
		visible, bright := 0, 0
		for y := 0; y < 32; y++ {
			for x := 0; x < 32; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if a < 0x8000 {
					continue
				}
				visible++
				if r+g+b > 3*0x7000 {
					bright++
				}
			}
		}
		if visible == 0 || bright == 0 {
			t.Fatalf("%s: icon must contain visible contrasting artwork", name)
		}
		ico, err := trayFiles.ReadFile("tray/" + name + ".ico")
		if err != nil || !bytes.Equal(ico[:4], []byte{0, 0, 1, 0}) {
			t.Fatalf("%s: not an ICO file: %v", name, err)
		}
		if len(trayIcon(name)) == 0 {
			t.Fatalf("%s: no icon for this system", name)
		}
	}
}
