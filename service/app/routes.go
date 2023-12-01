package app

func (a *App) Run() error {
	a.SetupRoutes()
	return a.server.Run()
}

func (a *App) SetupRoutes() {

	a.server.POST("/repository-configuration", a.HandleStrapiWebhook)

}
