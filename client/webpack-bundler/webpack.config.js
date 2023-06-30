const path = require('path');

module.exports = {
    entry: {
        // 'wallet_connect': './src/wallet_connect/index.js',
        'cosmos_kit': './src/cosmos_kit/pages/index.tsx',
        // 'cosmos_kit': './src/cosmos_kit/src/App.js',
    },
    module: {
        rules: [
            // `js` and `jsx` files are parsed using `babel`
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ["babel-loader"],
            },
            // `ts` and `tsx` files are parsed using `ts-loader`
            {
                test: /\.(ts|tsx)$/,
                loader: "ts-loader"
            }
        ],
    },
    resolve: {
        extensions: ["*", ".js", ".jsx", ".ts", ".tsx"],
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'dist'),
        library: {
            name: 'walletConnect',
            type: 'var',
        },
    },
    mode: 'production',
};