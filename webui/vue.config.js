module.exports = {
    devServer: {
        port: 4400,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            '/api': {
                target: 'http://172.17.124.37:9999/',
                changeOrigin: true,
            }
        }
    }
}