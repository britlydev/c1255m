package version

import "runtime/debug"

var (
	Version = ""
	Commit  = ""
	Date    = ""
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
}

func Get() Info {
	info := Info{Version: Version, Commit: Commit, Date: Date}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	info.GoVersion = bi.GoVersion
	if info.Version == "" {
		info.Version = bi.Main.Version
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if info.Commit == "" {
				info.Commit = s.Value
			}
		case "vcs.time":
			if info.Date == "" {
				info.Date = s.Value
			}
		}
	}

	return info
}
