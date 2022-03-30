var express = require('express');
var app = express();
var stringReplace = require('string-replace-middleware');

var KC_URL = "http://localhost:8080/auth";


var BACKEND_BASE_URL = "http://localhost:3000";

app.use(stringReplace({
  'BACKEND_BASE_URL': BACKEND_BASE_URL,
  'KC_URL': KC_URL
}));

app.use(express.static('.'))

app.get('/', function (req, res) {
  res.render('index.html');
});

app.get('/client.js', function (req, res) {
  res.render('client.js');
});

app.listen(8000);
