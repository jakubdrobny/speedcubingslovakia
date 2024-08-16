const express = require("express");
const crypto = require("crypto");
const app = express();
const port = 3456;

const cstimer = require("cstimer_module");
cstimer.setSeed(crypto.randomBytes(64).toString("hex"));

app.get("/api/v0/scramble/:scramblingCode", (req, res) => {
  const scramblingCode = req.params.scramblingCode;
  const numScrambles = parseInt(req.query.numScrambles);
  console.log(scramblingCode, numScrambles);

  let scrambles = [];
  for (let i = 0; i < numScrambles; i++)
    scrambles.push(cstimer.getScramble(scramblingCode));

  res.send(scrambles);
});

app.get("/api/v0/view/:scramblingCode/:imgType", (req, res) => {
  const scramblingCode = req.params.scramblingcode;
  const scramble = req.query.scramble;

  res.send(cstimer.getImage(scramble, scramblingCode));
});

app.listen(port, () => {
  console.log(`Scrambling service listening on port ${port}.`);
});
