const path = require('path');

module.exports = {
    entry: './src/index.js',
    output: {
        filename: 'wallet_connect.js',
        path: path.resolve(__dirname, 'dist'),
        library: {
            name: 'walletConnect',
            type: 'var',
        },
    },
    mode: 'production',
};