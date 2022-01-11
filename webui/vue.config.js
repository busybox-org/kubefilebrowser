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
                target: 'http://172.24.175.30:9999/',
                changeOrigin: true,
            }
        }
    }
}