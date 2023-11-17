'use strict';

/**
 * repository router
 */

const { createCoreRouter } = require('@strapi/strapi').factories;

module.exports = createCoreRouter('api::repository.repository');
