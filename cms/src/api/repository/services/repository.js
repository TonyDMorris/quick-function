'use strict';

/**
 * repository service
 */

const { createCoreService } = require('@strapi/strapi').factories;

module.exports = createCoreService('api::repository.repository');
