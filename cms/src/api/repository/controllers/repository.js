"use strict";

/**
 * repository controller
 */

const { createCoreController } = require("@strapi/strapi").factories;

module.exports = createCoreController(
  "api::repository.repository",
  ({ strapi }) => ({
    async find(ctx) {
      return {
        success: false,
        message: "Not implemented",
      };

      // some more custom logic
    },
    async findOne(ctx) {
      return {
        success: false,
        message: "Not implemented",
      };

      // some more custom logic
    },
    async create(ctx) {
      return {
        success: false,
        message: "Not implemented",
      };

      // some more custom logic
    },
    async update(ctx) {
      return {
        success: false,
        message: "Not implemented",
      };

      // some more custom logic
    },
  })
);
