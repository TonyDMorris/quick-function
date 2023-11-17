'use strict';

/**
 * repository controller
 */

const { createCoreController } = require('@strapi/strapi').factories;

module.exports = createCoreController('api::repository.repository');
