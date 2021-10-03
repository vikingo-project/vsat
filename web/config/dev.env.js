"use strict";
const merge = require("webpack-merge");
const prodEnv = require("./prod.env");

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  VUE_APP_GATE_URL: process.env.VUE_APP_GATE_URL
    ? `"${process.env.VUE_APP_GATE_URL}"`
    : '"https://api.vikingo.org/api/gate/"'
});
