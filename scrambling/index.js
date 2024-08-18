const express = require("express");
const crypto = require("crypto");
const app = express();

const port =
  require("dotenv").config({ path: `.env.${process.env.NODE_ENV}` }).parsed
    .PORT || 3999;

const cstimer = require("cstimer_module");
cstimer.setSeed(crypto.randomBytes(64).toString("hex"));

app.get("/api/v0/scramble/:scramblingCode", (req, res) => {
  const scramblingCode = req.params.scramblingCode;
  const numScrambles = parseInt(req.query.numScrambles);

  let scrambles = [];
  for (let i = 0; i < numScrambles; i++) {
    switch (scramblingCode) {
      case "555wca":
        scrambles.push(cstimer.getScramble(scramblingCode, 60));
        break;
      case "666wca":
        scrambles.push(cstimer.getScramble(scramblingCode, 80));
        break;
      case "777wca":
        scrambles.push(cstimer.getScramble(scramblingCode, 100));
        break;
      case "mgmp":
        scrambles.push(cstimer.getScramble(scramblingCode, 70));
        break;
      case "555bld":
        scrambles.push(cstimer.getScramble(scramblingCode, 60));
        break;
      default:
        scrambles.push(cstimer.getScramble(scramblingCode));
        break;
    }
  }

  console.log(`/api/v0/scramble/${scramblingCode}: `, scrambles);
  console.log("========================================");

  res.send(scrambles);
});

app.get("/api/v0/view/:scramblingCode/:imgType", (req, res) => {
  const scramblingCode = req.params.scramblingCode;
  const imgType = req.params.imgType;
  const scramble = req.query.scramble;

  console.log(`/api/v0/view/${scramblingCode}/${imgType}: `, imgType);
  console.log("========================================");

  res.send(cstimer.getImage(scramble, scramblingCode));
});

app.listen(port, () => {
  console.log(`Scrambling service listening on port ${port}.`);
});
