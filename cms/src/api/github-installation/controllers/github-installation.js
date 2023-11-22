"use strict";

/**
 * A set of functions called "actions" for `github-installation`
 */

const { entityService } = require("@strapi/strapi").factories;

module.exports = {
  handleInstallationWebhook: async (ctx, next) => {
    console.log(ctx.request.body);
    switch (ctx.request.body.action) {
      case "created":
        return await installationCreateHandler(ctx.request.body, next);

      case "deleted":
        return await installationDeleteHandler(ctx.request.body, next);

      case "added":
        break;

      case "removed":
        return await installationRemoveHandler(ctx.request.body, next);
    }
  },
};

const installationAddHandler = async (body, next) => {};

const installationRemoveHandler = async (body, next) => {};

const installationCreateHandler = async (body, next) => {
  const { installation, repositories } = body;

  const { id, account } = installation;
  try {
    const existingRepositories = await strapi.entityService.findMany(
      "api::repository.repository",
      {
        filters: {
          full_name: { $in: repositories.map((repo) => repo.full_name) },
        },
      }
    );

    const newRepositories = repositories.filter((repo) => {
      return !existingRepositories.find(
        (existingRepo) => existingRepo.full_name === repo.full_name
      );
    });

    const newRepositoryEntities = newRepositories.map((repo) => {
      return {
        full_name: repo.full_name,
        name: repo.name,
        private: repo.private,
        repository_id: repo.id.toString(),
      };
    });
    const repoIDs = existingRepositories.map((repo) => repo.id);
    for (const entity of newRepositoryEntities) {
      const repo = await strapi.entityService.create(
        "api::repository.repository",
        { data: entity }
      );
      repoIDs.push(repo.id);
    }

    const installationEntity = await strapi.entityService.create(
      "api::installation.installation",
      {
        data: {
          username: account.login,
          installation_id: id.toString(),
          repositories: repoIDs,
        },
      }
    );
  } catch (error) {
    console.log(error);
    return error;
  }

  return { success: true };
};

const installationDeleteHandler = async (body, next) => {
  try {
    const { installation } = body;
    const { id } = installation;
    const existingInstallation = await strapi.entityService.findMany(
      "api::installation.installation",
      {
        filters: {
          installation_id: id.toString(),
        },
      }
    );
    if (existingInstallation.length === 0) {
      return { success: true };
    }
    const installationToDelete = existingInstallation[0];

    await strapi.entityService.delete(
      "api::installation.installation",
      installationToDelete.id
    );
  } catch (error) {
    console.log(error);
    return error;
  }
  return { success: true };
};
