"use strict";

/**
 * A set of functions called "actions" for `github-installation`
 */

const { entityService } = require("@strapi/strapi").factories;
const crypto = require("crypto");

const secret = process.env.GITHUB_WEBHOOK_SECRET;

// For these headers, a sigHashAlg of sha1 must be used instead of sha256
// GitHub: X-Hub-Signature
// Gogs:   X-Gogs-Signature
const sigHeaderName = "x-hub-signature";
const sigHashAlg = "sha1";

module.exports = {
  handleInstallationWebhook: async (ctx, next) => {
    const test = process.env.TEST;
    if (!test) {
      const data = JSON.stringify(ctx.request.body);
      const headerData = JSON.stringify(ctx.request.headers);

      const header = ctx.request.headers[sigHeaderName] || "";

      const sig = Buffer.from(header, "utf8");

      const hmac = crypto.createHmac(sigHashAlg, secret);

      const digest = Buffer.from(
        `${sigHashAlg}=${hmac.update(data).digest("hex")}`,
        "utf8"
      );

      if (
        sig.length !== digest.length ||
        !crypto.timingSafeEqual(digest, sig)
      ) {
        strapi.log.debug(
          `Request body digest (${digest}) did not match ${sigHeaderName} (${sig})`
        );
        return ctx.throw(403, "Invalid signature");
      }
      strapi.log.debug(ctx.request.body);
    }
    switch (ctx.request.body.action) {
      case "created":
        return await installationCreateHandler(ctx.request.body, next);

      case "deleted":
        return await installationDeleteHandler(ctx.request.body, next);

      case "added":
        return await installationAddHandler(ctx.request.body, next);

      case "removed":
        return await installationRemoveHandler(ctx.request.body, next);
    }
  },
};

const installationAddHandler = async (body, next) => {
  const { installation, repositories_added } = body;

  const { id, account } = installation;
  try {
    const existingRepositories = await strapi.entityService.findMany(
      "api::repository.repository",
      {
        filters: {
          full_name: { $in: repositories_added.map((repo) => repo.full_name) },
        },
      }
    );

    const newRepositories = repositories_added.filter((repo) => {
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

    const existingInstallation = await strapi.entityService.findMany(
      "api::installation.installation",
      {
        filters: {
          installation_id: id.toString(),
        },
        populate: ["repositories"],
      }
    );

    if (existingInstallation.length === 0) {
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
    } else {
      const installationEntity = existingInstallation[0];
      const newRepoIDs = [
        ...installationEntity.repositories
          .map((repo) => repo.id)
          .filter((id) => {
            return !repoIDs.includes(id);
          }),
        ...repoIDs,
      ];
      await strapi.entityService.update(
        "api::installation.installation",
        installationEntity.id,
        {
          data: {
            repositories: newRepoIDs,
          },
        }
      );
    }
  } catch (error) {
    console.log(error);
    return error;
  }

  return { success: true };
};

const installationRemoveHandler = async (body, next) => {
  const { installation, repositories_removed } = body;

  const { id, account } = installation;
  try {
    const existingRepositories = await strapi.entityService.findMany(
      "api::repository.repository",
      {
        filters: {
          full_name: {
            $in: repositories_removed.map((repo) => repo.full_name),
          },
        },
      }
    );

    const existingRepositoryIDs = existingRepositories.map((repo) => repo.id);

    const existingInstallation = await strapi.entityService.findMany(
      "api::installation.installation",
      {
        filters: {
          installation_id: id.toString(),
        },
        populate: ["repositories"],
      }
    );

    if (existingInstallation.length === 0) {
      return { success: true };
    } else {
      const installationEntity = existingInstallation[0];
      const newRepoIDs = [
        ...installationEntity.repositories
          .map((repo) => repo.id)
          .filter((id) => {
            return !existingRepositoryIDs.includes(id);
          }),
      ];
      await strapi.entityService.update(
        "api::installation.installation",
        installationEntity.id,
        {
          data: {
            repositories: newRepoIDs,
          },
        }
      );
    }
  } catch (error) {
    console.log(error);
    return error;
  }

  return { success: true };
};

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
