module.exports = {
    transpileDependencies: [],
    devServer: {
        port: 4400,
        open: true,
        proxy: {
            '/api': {
                target: 'http://wslhost:9999/',
                changeOrigin: true,
            }
        }
    }
}