package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// tarFile is an entry of a test archive: a file with its content and mode,
// a folder (a name that ends in "/"), or a link.
type tarFile struct {
	name, content     string
	mode              int64
	symlink, hardlink string
}

func tarExe(name, content string) tarFile { return tarFile{name: name, content: content, mode: 0o755} }

func tarRegular(name, content string) tarFile {
	return tarFile{name: name, content: content, mode: 0o644}
}

func tarFolder(name string) tarFile { return tarFile{name: name} }

func tarSymlink(name, to string) tarFile { return tarFile{name: name, symlink: to} }

func tarHardlink(name, to string) tarFile { return tarFile{name: name, hardlink: to} }

// appTar is a tar.gz of files, in their order.
func appTar(t *testing.T, files ...tarFile) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	w := tar.NewWriter(gz)
	for _, f := range files {
		hdr := &tar.Header{Name: f.name, Mode: f.mode, Size: int64(len(f.content)), Typeflag: tar.TypeReg}
		switch {
		case strings.HasSuffix(f.name, "/"):
			hdr.Typeflag, hdr.Size, hdr.Mode = tar.TypeDir, 0, 0o755
		case f.symlink != "":
			hdr.Typeflag, hdr.Size, hdr.Mode, hdr.Linkname = tar.TypeSymlink, 0, 0o755, f.symlink
		case f.hardlink != "":
			hdr.Typeflag, hdr.Size, hdr.Linkname = tar.TypeLink, 0, f.hardlink
		}
		if err := w.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if hdr.Typeflag == tar.TypeReg {
			w.Write([]byte(f.content))
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	gz.Close()
	return buf.Bytes()
}

func writeTar(t *testing.T, files ...tarFile) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "x.tar.gz")
	if err := os.WriteFile(file, appTar(t, files...), 0o644); err != nil {
		t.Fatal(err)
	}
	return file
}

// The release's archive: one akad-vpn folder, whose contents go into
// the folder, executables staying executable. GNU tar's "./" entry first
// changes nothing.
func TestUntarUnpacksTheAppFolder(t *testing.T) {
	file := writeTar(t,
		tarFolder("./"),
		tarFolder("akad-vpn/"),
		tarExe("akad-vpn/akad-vpn", "app"),
		tarRegular("akad-vpn/lib/libapp.so", "lib"),
		tarExe("akad-vpn/core/akad-vpn-core", "core"),
	)
	out := filepath.Join(t.TempDir(), "out")
	if err := untar(file, out); err != nil {
		t.Fatal(err)
	}
	if read(t, filepath.Join(out, "akad-vpn")) != "app" || read(t, filepath.Join(out, "core", "akad-vpn-core")) != "core" || read(t, filepath.Join(out, "lib", "libapp.so")) != "lib" {
		t.Fatal("wrong layout")
	}
	if runtime.GOOS != "windows" {
		info, _ := os.Stat(filepath.Join(out, "akad-vpn"))
		lib, _ := os.Stat(filepath.Join(out, "lib", "libapp.so"))
		if info.Mode().Perm()&0o111 == 0 || lib.Mode().Perm()&0o111 != 0 {
			t.Fatalf("modes %v %v", info.Mode(), lib.Mode())
		}
	}
	if _, err := os.Stat(out + ".part"); err == nil {
		t.Fatal("the unpacking folder is left")
	}
}

// A macOS bundle's frameworks link within themselves; the links stay links.
func TestUntarKeepsLinksWithinTheApp(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("links need Windows' developer mode")
	}
	fw := "AkadVPN.app/Contents/Frameworks/FlutterMacOS.framework/"
	file := writeTar(t,
		tarExe("AkadVPN.app/Contents/MacOS/AkadVPN", "app"),
		tarSymlink(fw+"FlutterMacOS", "Versions/Current/FlutterMacOS"),
		tarSymlink(fw+"Versions/Current", "A"),
		tarExe(fw+"Versions/A/FlutterMacOS", "flutter"),
	)
	out := filepath.Join(t.TempDir(), "out")
	if err := untar(file, out); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(out, "Contents", "Frameworks", "FlutterMacOS.framework", "FlutterMacOS")
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a link: %v %v", info, err)
	}
	if read(t, link) != "flutter" {
		t.Fatal("the link leads elsewhere")
	}
}

func TestUntarRefusesNamesOutsideTheFolderAndLinks(t *testing.T) {
	for _, bad := range [][]tarFile{
		{tarRegular("../evil", "x")},
		{tarRegular("akad-vpn/../../evil", "x")},
		// Inside the archive, but outside the folder once its top is cut.
		{tarRegular("akad-vpn/../evil", "x")},
		{tarRegular("/abs", "x")},
		{tarSymlink("akad-vpn/core/passwd", "/etc/passwd")},
		{tarSymlink("akad-vpn/core/up", "../../etc")},
		{tarHardlink("akad-vpn/core/hard", "akad-vpn/akad-vpn")},
		// Each link stays inside on its own; followed, the second leads
		// out: l is the app's folder, so l/.. is its parent.
		{tarSymlink("akad-vpn/deep/deeper/l", "../.."), tarSymlink("akad-vpn/out", "deep/deeper/l/..")},
		// A link made inside another link that leads to its own folder.
		{tarSymlink("akad-vpn/here", "."), tarSymlink("akad-vpn/here/up", "..")},
		// APFS ignores case: "Here" is the link "here".
		{tarSymlink("akad-vpn/here", "."), tarSymlink("akad-vpn/out", "Here/..")},
	} {
		file := writeTar(t, append([]tarFile{tarExe("akad-vpn/akad-vpn", "app")}, bad...)...)
		if err := untar(file, filepath.Join(t.TempDir(), "out")); err == nil {
			t.Errorf("%q: unpacked", bad[len(bad)-1].name)
		}
	}
}

// A Linux release downloads and installs as the Windows one does: its
// tar.gz for this processor, checked against the signed list.
func TestALinuxReleaseInstalls(t *testing.T) {
	thePackage = linuxPackage
	t.Cleanup(func() { thePackage = windowsPackage })
	r := newRelease(t, "1.2.3")
	r.files[AssetName("1.2.3")] = appTar(t,
		tarExe("akad-vpn/akad-vpn", "app 1.2.3"),
		tarExe("akad-vpn/core/akad-vpn-core", "core 1.2.3"),
		tarExe("akad-vpn/core/akad-vpn-tun", "tun 1.2.3"),
	)
	r.sign()
	r.serve(t)
	if !strings.HasSuffix(AssetName("1.2.3"), "-Linux-"+runtime.GOARCH+".tar.gz") {
		t.Fatalf("asset %s", AssetName("1.2.3"))
	}
	appDir, dataDir := appFolder(t, "1.0.0", false)
	write(t, filepath.Join(appDir, "akad-vpn"), "app 1.0.0")
	p := installed(t, appDir, dataDir)
	if read(t, filepath.Join(p.Package, "core", "akad-vpn-core")) != "core 1.2.3" || read(t, filepath.Join(p.Package, "akad-vpn")) != "app 1.2.3" {
		t.Fatalf("unpacked into %s", p.Package)
	}
}

// A macOS release is the app's bundle; its Contents replaces the bundle's.
func TestAMacReleaseInstalls(t *testing.T) {
	thePackage = macPackage
	t.Cleanup(func() { thePackage = windowsPackage })
	r := newRelease(t, "1.2.3")
	files := []tarFile{
		tarRegular("AkadVPN.app/Contents/Info.plist", "plist 1.2.3"),
		tarExe("AkadVPN.app/Contents/MacOS/AkadVPN", "app 1.2.3"),
		tarExe("AkadVPN.app/Contents/MacOS/akad-vpn-core", "core 1.2.3"),
		tarExe("AkadVPN.app/Contents/MacOS/akad-vpn-tun", "tun 1.2.3"),
	}
	if runtime.GOOS != "windows" {
		files = append(files, tarSymlink("AkadVPN.app/Contents/Frameworks/App.framework/App", "Versions/Current/App"))
	}
	r.files[AssetName("1.2.3")] = appTar(t, files...)
	r.sign()
	r.serve(t)
	if AssetName("1.2.3") != "AkadVPN-1.2.3-macOS.tar.gz" {
		t.Fatalf("asset %s", AssetName("1.2.3"))
	}
	appDir, dataDir := appFolder(t, "1.0.0", false)
	p := installed(t, appDir, dataDir)
	contents := filepath.Join(p.Package, "Contents")
	if read(t, filepath.Join(contents, "MacOS", "akad-vpn-core")) != "core 1.2.3" || read(t, filepath.Join(contents, "MacOS", "AkadVPN")) != "app 1.2.3" {
		t.Fatalf("unpacked into %s", p.Package)
	}
	if entries, _ := os.ReadDir(p.Package); len(entries) != 1 {
		t.Fatalf("the package holds more than Contents: %v", entries)
	}
}

// The app's folder holds the core: core/ on Windows and Linux, the bundle
// on macOS.
func TestAppDirIsTheFolderThatHoldsTheCore(t *testing.T) {
	t.Cleanup(func() { thePackage = windowsPackage })
	for _, c := range []struct {
		pkg       appPackage
		exe, want string
	}{
		{windowsPackage, "/Users/a/AkadVPN/core/akad-vpn-core.exe", "/Users/a/AkadVPN"},
		{linuxPackage, "/home/a/.local/opt/akad-vpn/core/akad-vpn-core", "/home/a/.local/opt/akad-vpn"},
		{macPackage, "/Applications/AkadVPN.app/Contents/MacOS/akad-vpn-core", "/Applications/AkadVPN.app"},
	} {
		thePackage = c.pkg
		if got := AppDir(filepath.FromSlash(c.exe)); got != filepath.FromSlash(c.want) {
			t.Errorf("%s: %s", c.exe, got)
		}
	}
}
