"use strict";

module.exports = {
  routes: [
    {
      method: "GET",
      path: "/internal/repository-configuration:id",
      handler: "repository-configuration.internalFind",
      config: {
        policies: [],
        middlewares: [],
      },
    },
  ],
};
