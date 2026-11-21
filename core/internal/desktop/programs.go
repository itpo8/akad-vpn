package desktop

import (
	"slices"
	"strings"

	"github.com/itpo8/akad-vpn/core/internal/service"
	"github.com/itpo8/akad-vpn/core/internal/tunapi"
)

// ownPrograms are this app's own executables, which never skip the VPN by
// choice.
var ownPrograms = map[string]bool{
	"akad-vpn.exe": true, "akad-vpn-core.exe": true, "akad-vpn-tun.exe": true,
	"akad-vpn": true, "akad-vpn-core": true, "akad-vpn-tun": true,
}

// runningPrograms are the programs running on this computer, for the TUN
// settings' choice of programs that skip the VPN (GET /api/tun/programs):
// each once by its file name, without the system's own and this app's, in
// the order of their names. listProcesses gives their executables' paths.
func runningPrograms() ([]service.RunningProgram, error) {
	paths, err := listProcesses()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	out := make([]service.RunningProgram, 0, len(paths))
	for _, path := range paths {
		name := path[strings.LastIndexAny(path, `/\`)+1:]
		key := strings.ToLower(name)
		if name == "" || seen[key] || ownPrograms[key] || tunapi.SystemProgram(path) {
			continue
		}
		seen[key] = true
		out = append(out, service.RunningProgram{Name: name, Path: path})
	}
	slices.SortFunc(out, func(a, b service.RunningProgram) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	return out, nil
}
