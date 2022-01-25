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
                target: 'http://172.19.5.115:9999/',
                changeOrigin: true,
            }
        }
    }
}