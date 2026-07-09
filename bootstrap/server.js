const http = require("node:http");
const fs = require("node:fs");
const path = require("node:path");

const PORT = process.env.PORT || 80;
const SCRIPTS = new Set(["/install.sh", "/upgrade.sh"]);

const server = http.createServer((req, res) => {
  if (req.url === "/version") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ version: "0.1.0", latest: true }));
    return;
  }

  const file = path.join(
    __dirname,
    req.url === "/" ? "install.sh" : path.basename(req.url),
  );
  const name = path.basename(file);

  if (!SCRIPTS.has("/" + name) || !fs.existsSync(file)) {
    res.writeHead(404);
    res.end("Not found");
    return;
  }

  res.writeHead(200, { "Content-Type": "text/x-shellscript" });
  fs.createReadStream(file).pipe(res);
});

server.listen(PORT, () => {
  console.log(`📦 install server running on port ${PORT}`);
});
