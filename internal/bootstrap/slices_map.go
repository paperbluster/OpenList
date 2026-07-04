package bootstrap

import (
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"
)

// InitSlicesMap populates conf.SlicesMap from the settings database.
// It reads comma-separated setting values and splits them into slices.
// Must be called after InitDB() and data.InitData().
func InitSlicesMap() {
	// Settings keys that contain comma-separated lists
	sliceKeys := []string{
		conf.ImageTypes,
		conf.VideoTypes,
		conf.AudioTypes,
		conf.TextTypes,
		conf.ProxyTypes,
		conf.ProxyIgnoreHeaders,
		conf.IgnoreDirectLinkParams,
	}

	for _, key := range sliceKeys {
		val := strings.TrimSpace(setting.GetStr(key))
		if val != "" {
			parts := strings.Split(val, ",")
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			conf.SlicesMap[key] = parts
		} else {
			conf.SlicesMap[key] = []string{}
		}
	}

	// Register callback so SlicesMap stays in sync when settings change
	op.RegisterSettingChangingCallback(InitSlicesMap)
}
