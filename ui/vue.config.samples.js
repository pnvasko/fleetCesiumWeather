module.exports = {
    pages: {
        'index': {
            entry: './src/main.ts',
            template: 'public/index.html',
            title: 'Welcome to my vue generator project',
            chunks: ['chunk-vendors', 'chunk-common', 'index']
        },
        'bad': {
            entry: './src/error-instance.ts',
            template: 'public/bad.html',
            title: 'Error page',
            chunks: ['chunk-vendors', 'chunk-common', 'index']
        },
        /* Disabled - Only one time
        'googleVerify': {
          entry: './src/error-instance.ts',
          template: 'public/somelink.html',
          title: 'Error page',
          chunks: ['chunk-vendors', 'chunk-common', 'index']
        },
        */

    },
    devServer: {
        'port': 3000
    },
    css: {
        sourceMap: false
    },
    pwa: {
        name: 'My App',
        themeColor: '#4DBA87',
        msTileColor: '#000000',
        appleMobileWebAppCapable: 'yes',
        appleMobileWebAppStatusBarStyle: 'black',
    },
}
// ---------------------------------------------------------------------------------------------------------------------
// vue.config.js
module.exports = {
    configureWebpack: {
        plugins: [new MyAwesomeWebpackPlugin()]
    }
}
// ---------------------------------------------------------------------------------------------------------------------
module.exports = {
    configureWebpack: {
        plugins: [
            new MyAwesomeWebpackPlugin()
        ]
    }
}

// In my case i was doing below and was facing the same error
const path = require('path')
const PrerenderSPAPlugin = require('prerender-spa-plugin')

module.exports = {
    configureWebpack: {
        plugins: [
            new PrerenderSPAPlugin({
                // Required - The path to the webpack-outputted app to prerender.
                staticDir: path.join(__dirname, 'dist'),
                // Required - Routes to render.
                routes: ['/', '/page-path', '/another-page'],
            })
        ]
    },
    lintOnSave: false
}

// ---------------------------------------------------------------------------------------------------------------------

const StylelintPlugin = require('stylelint-webpack-plugin')

module.exports = {

    assetsDir: 'asset',

    configureWebpack: config => {
        config.entry = '@/wrapper/main.js'
    },

    chainWebpack: config => {
        config.plugins.delete('prefetch')

        config.plugin('stylelint').use(StylelintPlugin, [
            {
                files: '**/*.s?(a|c)ss'
            }
        ])
    },

    lintOnSave: undefined,
    runtimeCompiler: true
}
