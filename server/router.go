package server

import (
	"github.com/OpenListTeam/OpenList/v4/cmd/flags"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/sign"
	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/OpenListTeam/OpenList/v4/server/handles"
	"github.com/OpenListTeam/OpenList/v4/server/middlewares"
	"github.com/OpenListTeam/OpenList/v4/server/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	e.ContextWithFallback = true
	if !utils.SliceContains([]string{"", "/"}, conf.URL.Path) {
		e.GET("/", func(c *gin.Context) {
			c.Redirect(302, conf.URL.Path)
		})
	}
	Cors(e)
	g := e.Group(conf.URL.Path)
	if conf.Conf.Scheme.HttpPort != -1 && conf.Conf.Scheme.HttpsPort != -1 && conf.Conf.Scheme.ForceHttps {
		e.Use(middlewares.ForceHttps)
	}
	g.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	g.GET("/favicon.ico", handles.Favicon)
	g.GET("/robots.txt", handles.Robots)
	g.GET("/manifest.json", static.ManifestJSON)
	g.GET("/i/:link_name", handles.Plist)
	common.SecretKey = []byte(conf.Conf.JwtSecret)
	g.Use(middlewares.StoragesLoaded)
	if conf.Conf.MaxConnections > 0 {
		g.Use(middlewares.MaxAllowed(conf.Conf.MaxConnections))
	}
	WebDav(g.Group("/dav"))

	downloadLimiter := middlewares.DownloadRateLimiter(stream.ClientDownloadLimit)
	signCheck := middlewares.Down(sign.Verify)
	g.GET("/d/*path", middlewares.PathParse, signCheck, downloadLimiter, handles.Down)
	g.GET("/p/*path", middlewares.PathParse, signCheck, downloadLimiter, handles.Proxy)
	g.HEAD("/d/*path", middlewares.PathParse, signCheck, handles.Down)
	g.HEAD("/p/*path", middlewares.PathParse, signCheck, handles.Proxy)

	g.GET("/sd/:sid", middlewares.EmptyPathParse, middlewares.SharingIdParse, downloadLimiter, handles.SharingDown)
	g.GET("/sd/:sid/*path", middlewares.PathParse, middlewares.SharingIdParse, downloadLimiter, handles.SharingDown)
	g.HEAD("/sd/:sid", middlewares.EmptyPathParse, middlewares.SharingIdParse, handles.SharingDown)
	g.HEAD("/sd/:sid/*path", middlewares.PathParse, middlewares.SharingIdParse, handles.SharingDown)

	api := g.Group("/api")
	auth := api.Group("", middlewares.Auth(false))

	api.POST("/auth/login", handles.Login)
	api.POST("/auth/login/hash", handles.LoginHash)
	auth.GET("/me", handles.CurrentUser)
	auth.POST("/me/update", handles.UpdateCurrent)
	auth.GET("/auth/logout", handles.LogOut)

	// no need auth
	public := api.Group("/public")
	public.Any("/settings", handles.PublicSettings)
	public.Any("/archive_extensions", handles.PublicArchiveExtensions)

	_fs(auth.Group("/fs"))
	fsAndShare(api.Group("/fs", middlewares.Auth(true)))
	_task(auth.Group("/task", middlewares.AuthNotGuest))
	_sharing(auth.Group("/share", middlewares.AuthNotGuest))
	admin(auth.Group("/admin", middlewares.AuthAdmin))
	if flags.Debug || flags.Dev {
		debug(g.Group("/debug"))
	}
	static.Static(g, func(handlers ...gin.HandlerFunc) {
		e.NoRoute(handlers...)
	})
}

func admin(g *gin.RouterGroup) {
	user := g.Group("/user")
	user.GET("/list", handles.ListUsers)
	user.GET("/get", handles.GetUser)
	user.POST("/create", handles.CreateUser)
	user.POST("/update", handles.UpdateUser)
	user.POST("/delete", handles.DeleteUser)
	user.POST("/del_cache", handles.DelUserCache)

	storage := g.Group("/storage")
	storage.GET("/list", handles.ListStorages)
	storage.GET("/get", handles.GetStorage)
	storage.POST("/create", handles.CreateStorage)
	storage.POST("/update", handles.UpdateStorage)
	storage.POST("/delete", handles.DeleteStorage)
	storage.POST("/enable", handles.EnableStorage)
	storage.POST("/disable", handles.DisableStorage)
	storage.POST("/load_all", handles.LoadAllStorages)

	driver := g.Group("/driver")
	driver.GET("/list", handles.ListDriverInfo)
	driver.GET("/names", handles.ListDriverNames)
	driver.GET("/info", handles.GetDriverInfo)

	setting := g.Group("/setting")
	setting.GET("/get", handles.GetSetting)
	setting.GET("/list", handles.ListSettings)
	setting.POST("/save", handles.SaveSettings)
	setting.POST("/delete", handles.DeleteSetting)
	setting.POST("/default", handles.DefaultSettings)
	setting.POST("/reset_token", handles.ResetToken)

	index := g.Group("/index")
	index.GET("/progress", handles.GetIndexProgress)
	index.POST("/build", handles.BuildSearchIndex)
	index.POST("/stop", handles.StopIndexBuild)

	scan := g.Group("/scan")
	scan.GET("/progress", handles.GetManualScanProgress)
	scan.POST("/start", handles.StartManualScan)
	scan.POST("/stop", handles.StopManualScan)

	// retain /admin/task API to ensure compatibility with legacy automation scripts
	_task(g.Group("/task"))
}

func fsAndShare(g *gin.RouterGroup) {
	g.Any("/list", handles.FsListSplit)
	g.Any("/get", handles.FsGetSplit)
}

func _fs(g *gin.RouterGroup) {
	g.Any("/search", handles.Search)
	g.Any("/other", handles.FsOther)
	g.Any("/dirs", handles.FsDirs)
	g.POST("/mkdir", handles.FsMkdir)
	g.POST("/rename", handles.FsRename)
	g.POST("/batch_rename", handles.FsBatchRename)
	g.POST("/regex_rename", handles.FsRegexRename)
	g.POST("/move", handles.FsMove)
	g.POST("/recursive_move", handles.FsRecursiveMove)
	g.POST("/copy", handles.FsCopy)
	g.POST("/remove", handles.FsRemove)
	g.POST("/remove_empty_directory", handles.FsRemoveEmptyDirectory)
	uploadLimiter := middlewares.UploadRateLimiter(stream.ClientUploadLimit)
	g.PUT("/put", middlewares.FsUp, uploadLimiter, handles.FsStream)
	g.PUT("/form", middlewares.FsUp, uploadLimiter, handles.FsForm)
	g.POST("/link", middlewares.AuthAdmin, handles.Link)
	// Direct upload (client-side upload to storage)
	g.POST("/get_direct_upload_info", middlewares.FsUp, handles.FsGetDirectUploadInfo)
}

func _task(g *gin.RouterGroup) {
	handles.SetupTaskRoute(g)
}

func _sharing(g *gin.RouterGroup) {
	g.Any("/list", handles.ListSharings)
	g.GET("/get", handles.GetSharing)
	g.POST("/create", handles.CreateSharing)
	g.POST("/update", handles.UpdateSharing)
	g.POST("/delete", handles.DeleteSharing)
	g.POST("/enable", handles.SetEnableSharing(false))
	g.POST("/disable", handles.SetEnableSharing(true))
}

func Cors(r *gin.Engine) {
	config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	config.AllowOrigins = conf.Conf.Cors.AllowOrigins
	config.AllowHeaders = conf.Conf.Cors.AllowHeaders
	config.AllowMethods = conf.Conf.Cors.AllowMethods
	r.Use(cors.New(config))
}

