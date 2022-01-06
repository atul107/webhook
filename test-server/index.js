var express = require("express");
var app = express();
var bodyParser = require("body-parser");
app.use(bodyParser.json());
app.post("/", function (req, res) {
  console.log(req.headers);
  console.log(req.body);
  res.status(200).send("Success");
});
app.listen(3000, () => console.log("Example app listening on port 3000!"));
