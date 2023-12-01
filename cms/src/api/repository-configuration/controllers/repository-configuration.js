"use strict";

/**
 * repository-configuration controller
 */

const { createCoreController } = require("@strapi/strapi").factories;
const { entityService } = require("@strapi/strapi").factories;
module.exports = createCoreController(
  "api::repository-configuration.repository-configuration",
  ({ strapi }) => ({
    async create(ctx) {
      const user = ctx.state.user;
      const newConfig = ctx.request.body.data;
      console.log(newConfig);
      const existingInstallations = await strapi.entityService
        .findMany("api::installation.installation", {
          filters: {
            username: ctx.state.user.username,
          },
          populate: {
            repositories: {
              filters: {
                id: newConfig.repository,
              },
            },
          },
        })
        .catch((err) => {
          console.log(err);
          throw err;
        });

      if (existingInstallations.length === 0) {
        return { success: false, message: "Installation not found" };
      }
      if (existingInstallations[0].repositories.length === 0) {
        return {
          success: false,
          message: "Repository not found as part of this installation",
        };
      }

      const existingConfig = await strapi.entityService
        .findMany("api::repository-configuration.repository-configuration", {
          filters: {
            installation: {
              id: {
                $eq: existingInstallations[0].id,
              },
            },
            repository: {
              id: {
                $eq: existingInstallations[0].repositories[0].id,
              },
            },
          },
        })
        .catch((err) => {
          console.log(err);
          throw err;
        });

      if (existingConfig.length > 0) {
        return {
          success: false,
          message: "Repository configuration already exists",
        };
      }

      ctx.request.body.data = {
        ...newConfig,
        installation: existingInstallations[0].id,
        repository: existingInstallations[0].repositories[0].id,
      };

      return await super.create(ctx);
    },

    async update(ctx) {
      const user = ctx.state.user;
      const newConfig = ctx.request.body.data;
      const configID = ctx.params.id;

      const existingInstallations = await strapi.entityService
        .findMany("api::installation.installation", {
          filters: {
            username: ctx.state.user.username,
          },
          populate: {
            repository_configurations: {
              populate: {
                repository: true,
              },
              filters: {
                id: configID,
              },
            },
          },
        })
        .catch((err) => {
          console.log(err);
          throw err;
        });

      if (existingInstallations.length === 0) {
        return { success: false, message: "Installation not found" };
      }
      if (existingInstallations[0].repository_configurations.length === 0) {
        return {
          success: false,
          message:
            "Repository configuration not found as part of this installation",
        };
      }

      ctx.request.body.data = {
        ...newConfig,
        installation: existingInstallations[0].id,
        repository:
          existingInstallations[0].repository_configurations[0].repository.id,
        last_generation:
          existingInstallations[0].repository_configurations[0].last_generation,
      };

      return await super.update(ctx);
    },

    async delete(ctx) {
      const user = ctx.state.user;
      const configID = ctx.params.id;

      const existingInstallations = await strapi.entityService
        .findMany("api::installation.installation", {
          filters: {
            username: ctx.state.user.username,
          },
          populate: {
            repository_configurations: {
              populate: {
                repository: true,
              },
              filters: {
                id: configID,
              },
            },
          },
        })
        .catch((err) => {
          console.log(err);
          throw err;
        });

      if (existingInstallations.length === 0) {
        return { success: false, message: "Installation not found" };
      }
      if (existingInstallations[0].repository_configurations.length === 0) {
        return {
          success: false,
          message:
            "Repository configuration not found as part of this installation",
        };
      }

      ctx.params.id = existingInstallations[0].repository_configurations[0].id;

      return await super.delete(ctx);
    },

    async find() {
      return {
        success: false,
        message: "Not implemented",
      };
    },

    async findOne() {
      return {
        success: false,
        message: "Not implemented",
      };
    },

    async internalFind(ctx) {
      const id = ctx.params.id;
      const repositoryConfiguration = await strapi.entityService
        .findOne("api::repository-configuration.repository-configuration", id, {
          populate: {
            installation: true,
            repository: true,
          },
        })
        .catch((err) => {
          console.log(err);
          throw err;
        });

      if (!repositoryConfiguration) {
        return {
          success: false,
          message: "Repository configuration not found",
        };
      }

      return repositoryConfiguration;
    },
  })
);
