module.exports = {
  routes: [
    {
      method: "POST",
      path: "/github-installation",
      handler: "github-installation.handleInstallationWebhook",
      config: {
        policies: [],
        middlewares: [],
      },
    },
  ],
};
