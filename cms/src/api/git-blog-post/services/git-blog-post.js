'use strict';

/**
 * git-blog-post service
 */

const { createCoreService } = require('@strapi/strapi').factories;

module.exports = createCoreService('api::git-blog-post.git-blog-post');
