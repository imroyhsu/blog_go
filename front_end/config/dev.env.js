'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
    NODE_ENV: '"development"',
    API_DOMAIN: '"http://www.royhsu.cn"',
    API_PORT: '"8080"'
})
