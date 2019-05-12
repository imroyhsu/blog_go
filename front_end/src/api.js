var domain = process.env.API_DOMAIN
var port = process.env.API_PORT

global.signin = domain + ":" + port + "/login"
global.saveArticle = domain + ":" + port + "/articles"
global.updateArticle = domain + ":" + port + "/articles"
global.deleteArticle = domain + ":" + port + "/articles"
global.articleDetail = domain + ":" + port + "/articles/details"
global.articleList = domain + ":" + port + "/articles"
