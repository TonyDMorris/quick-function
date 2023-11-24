"use strict";

/**
 * installation router
 */

const { createCoreRouter } = require("@strapi/strapi").factories;

module.exports = createCoreRouter("api::installation.installation", {
  routes: [
    {
      method: "POST",
      path: "/installation",
      handler: "installation.handleInstallationWebhook",
      config: {
        policies: [],
        middlewares: [],
      },
    },
  ],
});
