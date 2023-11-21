'use strict';

/**
 * installation controller
 */

const { createCoreController } = require('@strapi/strapi').factories;

module.exports = createCoreController('api::installation.installation');
