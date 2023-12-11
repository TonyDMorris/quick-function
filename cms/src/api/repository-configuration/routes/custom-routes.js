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
    {
      method: "GET",
      path: "/internal/repository-configurations",
      handler: "repository-configuration.internalFindMany",
      config: {
        policies: [],
        middlewares: [],
      },
    },
    {
      method: "PUT",
      path: "/internal/repository-configurations/:id",
      handler: "repository-configuration.internalUpdate",
      config: {
        policies: [],
        middlewares: [],
      },
    },
  ],
};
