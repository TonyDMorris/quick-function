"use strict";

/**
 * installation controller
 */

const { createCoreController } = require("@strapi/strapi").factories;

module.exports = createCoreController(
  "api::installation.installation",
  ({ strapi }) => ({
    async find(ctx) {
      const user = ctx.state.user;
      ctx.query = {
        ...ctx.query,
        filters: {
          ...ctx.query.filters,
          username: user.username,
        },
      };

      const { data, meta } = await super.find(ctx);

      // some more custom logic
      meta.date = Date.now();

      return { data, meta };
    },
  })
);
