'use strict';

/**
 * repository-configuration service
 */

const { createCoreService } = require('@strapi/strapi').factories;

module.exports = createCoreService('api::repository-configuration.repository-configuration');
