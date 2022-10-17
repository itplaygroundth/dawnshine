"use strict";
const Session = require('./lib/session')
const fs = require('fs');
const express = require('express');
const bodyParser = require('body-parser')
const {body, validationResult, param, check} = require('express-validator');
const {forEach} = require('lodash/collection');
const multer = require('multer');
const Jimp = require('jimp');
const jsQr = require('jsqr');
const NodeCache = require('node-cache');
const AsyncPool = require('tiny-async-pool');
//const Session = require('./src/session');
const cookies = require("cookie-parser");


const app = express();
const port = process.env.PORT;

const config = require('./config');
const cache = new NodeCache();

const state = JSON.parse(fs.readFileSync(config.stateFile));

// const session = new Session(state, config.accountNo, config.accountType, config.pin);
 
// session.on('STATE_UPDATED', async state => {
//     await fs.promises.writeFile(config.stateFile, JSON.stringify(state, null, 4));
// });
app.use(cookies());
app.use(bodyParser.urlencoded({extended: false}));

const nodes = process.env.REDIS_CLUSTER.split(",").map((url) => {
    const [host, port] = url.split(":");
    return { host, port };
  }); 
const session = Session({
    sessionKey: 'sid',
    sidPrefix: 'sid-prefix-',
    expires: 30 * 60 * 1000,
    redisNodes: nodes
    // redisNodes:[{
    //     "host": "whoami.lab.lan",
    //     "port": '6377',
    // },
    // {
    //     "host": "whoami.lab.lan",
    //     "port": '6378',
        
    // },
    // {
    //     "host": "whoami.lab.lan",
    //     "port": '6379',
        
    // }    
    // ],
})


app.use(session.session)
app.use(session.save)
//app.use(Cluster.save)
 
app.get('/',async(req,res) => {
    try {
        res.json(req.session)
    }
    catch(e){
        res.status(400).json({error: e.message});
    }
})

// app.get('/balance', async (req, res) => {
//     try {
//         res.json(await session.call((client) => {
//             return client.getInquiryAccountBalance(session.accountNumber, session.accountType);
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.get('/activities', async (req, res) => {
//     try {
//         res.json(await session.call(async (client) => {
//             const response = await client.getAccountActivityList(session.accountNumber);

//             forEach(response['activityList'] ?? [], activity => {
//                 session._activityList[activity['rqUid']] = activity;
//             });

//             return response;
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.get('/custom/activities', async (req, res) => {
//     try {
//         res.json(await session.call(async (client) => {
//             const response = await client.getAccountActivityList(session.accountNumber);
//             const tasks = [];
//             forEach(response['activityList'] ?? [], activity => {
//                 session._activityList[activity['rqUid']] = activity;

//                 const cacheKey = `detail:${activity['rqUid']}`;
//                 const detail = cache.get(cacheKey);
//                 if (detail === undefined) {
//                     if (['CR', 'DR'].includes(activity['transactionType']) && activity['allowShowDetailFlag'] === 'Y') {
//                         tasks.push(async () => {
//                             let detail;
//                             try {
//                                 detail = await client.getAccountActivityDetail(session.accountNumber, activity);
//                             } catch (e) {
//                                 detail = null;
//                             }

//                             activity['_detail'] = detail;
//                             cache.set(cacheKey, detail);
//                         });
//                     }
//                 } else {
//                     activity['_detail'] = detail;
//                 }
//             });

//             await AsyncPool(2, tasks, v => v());

//             return response;
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.get('/activity-detail/:rqUid', param('rqUid').exists().bail().custom(async v => {
//     if (!session._activityList[v]) {
//         return Promise.reject('ไม่พบข้อมูล');
//     }
// }), async (req, res) => {
//     const errors = validationResult(req);
//     if (!errors.isEmpty()) {
//         return res.status(400).json({errors: errors.array()});
//     }

//     try {
//         const activity = session._activityList[req.params['rqUid']];
//         res.json(await session.call((client) => {
//             return client.getAccountActivityDetail(session.accountNumber, activity);
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.get('/bank-info-list', async (req, res) => {
//     try {
//         res.json(await session.getBankInfoList());
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.post('/inquire-for-transfer-money', [body('amount').exists().bail().isFloat({min: 0.01}), body('toAccount').exists(), body('toBankCode').exists().bail().custom(async v => {
//     const bankInfoList = await session.getBankInfoList();
//     if (!bankInfoList[v]) {
//         return Promise.reject('ไม่พบธนาคารปลายทาง');
//     }
// }),], async (req, res) => {
//     const errors = validationResult(req);
//     if (!errors.isEmpty()) {
//         return res.status(400).json({errors: errors.array()});
//     }

//     try {
//         const body = req.body;
//         res.json(await session.call(async (client) => {
//             const bankInfo = (await session.getBankInfoList())[body['toBankCode']];

//             const handle = {};
//             const response = await client.inquireForTransferMoney(handle, session.accountNumber, body['toAccount'], body['amount'], bankInfo['transferType'], bankInfo['targetBankCode']);
//             session._inquireForTransferMoneyList[response['kbankInternalSessionId']] = handle;

//             return response;
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.post('/transfer-money/:kbankInternalSessionId', param('kbankInternalSessionId').exists().bail().custom(async v => {
//     if (!session._inquireForTransferMoneyList[v]) {
//         return Promise.reject('ไม่พบข้อมูล');
//     }
// }), async (req, res) => {
//     const errors = validationResult(req);
//     if (!errors.isEmpty()) {
//         return res.status(400).json({errors: errors.array()});
//     }

//     try {
//         res.json(await session.call(async (client) => {
//             const handle = session._inquireForTransferMoneyList[req.params['kbankInternalSessionId']];

//             return client.transferMoney(handle);
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.post('/scan-qrcode/:raw', async (req, res) => {
//     try {
//         res.json(await session.call(async (client) => {
//             return client.scanQr(req.params['raw']);
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

// app.post('/scan-qrcode', check('image').custom((v, {req, res}) => {
//     const upload = multer().single('image');

//     return new Promise((resolve, reject) => {
//         upload(req, res, async function (err) {
//             if (req.fileValidationError) {
//                 return reject(req.fileValidationError);
//             } else if (err) {
//                 return reject(err);
//             }

//             try {
//                 const jimp = await Jimp.read(req.file.buffer);
//                 const qr = jsQr(jimp.bitmap.data, jimp.bitmap.width, jimp.bitmap.height);
//                 if (qr) {
//                     req.body.qrcode = qr.data;
//                 } else {
//                     reject('ไม่พบ qrcode');
//                 }
//             } catch (e) {
//                 return reject('รูปภาพไม่ถูกต้อง');
//             }

//             return resolve();
//         });
//     });
// }), async (req, res) => {
//     const errors = validationResult(req);
//     if (!errors.isEmpty()) {
//         return res.status(400).json({errors: errors.array()});
//     }

//     try {
//         res.json(await session.call(async (client) => {
//             return client.scanQr(req.body.qrcode);
//         }));
//     } catch (e) {
//         res.status(400).json({error: e.message});
//     }
// });

app.listen(port, () => {
    console.log(`App is running at port: ${port}`);
});