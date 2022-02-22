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
                target: 'http://192.168.93.194:9999/',
                changeOrigin: true,
            }
        }
    }
}