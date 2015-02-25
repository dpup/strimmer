// Hacks index.html to point at the compiled JS file instead of the development
// mode JSX file.
var fs = require('fs')
var INDEX_FILE = 'web/index.html'
var contents = fs.readFileSync(INDEX_FILE, {encoding: 'utf8'})
var out = contents.replace(/<!-- START SCRIPT -->[\S\s.]*<!-- END SCRIPT -->/gm, '<script src="app.js"></script>')
fs.writeFileSync(INDEX_FILE + '.bak', contents, {encoding: 'utf8'})
fs.writeFileSync(INDEX_FILE, out, {encoding: 'utf8'})
