const { defineConfig } = require('cypress')

module.exports = defineConfig({
    video: true,
    watchForFileChanges: false,
    e2e: {
        chromeWebSecurity: false,
        videoCompression: false,
        viewportWidth: 1920,
        viewportHeight: 1080
    }
})
