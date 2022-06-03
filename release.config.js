module.exports = {
  "branches": [
    { "name": "main" },
    { "name": "beta", "channel": "beta", "prerelease": "beta" },
  ],
  "plugins": [
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ]
}