package app

func (a *App) Run() error {
	a.setupRoutes()
	go a.startWorkerPool()
	go a.cron.StartAsync()
	err := a.loadSchedules()
	if err != nil {
		return err
	}
	return a.server.Run()
}

func (a *App) setupRoutes() {

	a.server.POST("/repository-configuration", a.HandleStrapiWebhook)

}
