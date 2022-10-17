const express = require('express');
const app = express();
const port = process.env.PORT;

app.get('/', (req, res) => {    res.json({msg:"Hello World"}) });
app.listen(port, () => {console.log('server runing on '+port) });      