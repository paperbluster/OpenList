package drivers

import (
	// 本地与虚拟存储
	_ "github.com/OpenListTeam/OpenList/v4/drivers/alias"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/crypt"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/local"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/strm"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/url_tree"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/virtual"

	// 标准协议
	_ "github.com/OpenListTeam/OpenList/v4/drivers/ftp"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/sftp"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/smb"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/webdav"

	// 对象存储 / 开放平台
	_ "github.com/OpenListTeam/OpenList/v4/drivers/alist_v3"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/azure_blob"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/cloudreve"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/cloudreve_v4"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/openlist"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/openlist_share"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/s3"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/uss"

	// PikPak
	_ "github.com/OpenListTeam/OpenList/v4/drivers/pikpak"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/pikpak_share"

	// 基础设施 (被上述驱动内部依赖)
	_ "github.com/OpenListTeam/OpenList/v4/drivers/autoindex"
	_ "github.com/OpenListTeam/OpenList/v4/drivers/chunk"
)

// All do nothing,just for import
// same as _ import
func All() {
}
