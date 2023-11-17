const crypto = require("crypto");

module.exports = ({ env }) => ({
  "users-permissions": {
    config: {
      jwtSecret:
        env("ADMIN_JWT_SECRET") || crypto.randomBytes(16).toString("base64"),
    },
  },
});
