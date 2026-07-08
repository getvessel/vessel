const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = process.env.PORT || 80;

const server = http.createServer((req, res) => {
  if (req.url === "/" || req.url === "/install.sh") {
    res.writeHead(200, { "Content-Type": "text/x-shellscript" });
    fs.createReadStream(path.join(__dirname, "install.sh")).pipe(res);
  } else if (req.url === "/upgrade.sh") {
    res.writeHead(200, { "Content-Type": "text/x-shellscript" });
    fs.createReadStream(path.join(__dirname, "upgrade.sh")).pipe(res);
  } else if (req.url === "/version") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ version: "0.1.0-alpha", latest: true }));
  } else {
    res.writeHead(404);
    res.end("Not found");
  }
});

server.listen(PORT, () => {
  console.log(`📦 get-vessel script delivery server running on port ${PORT}`);
});
