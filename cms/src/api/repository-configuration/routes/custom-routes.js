"use strict";

module.exports = {
  routes: [
    {
      method: "GET",
      path: "/internal/repository-configurations/:id",
      handler: "repository-configuration.internalFind",
      config: {
        policies: [],
        middlewares: [],
      },
    },
  ],
};
