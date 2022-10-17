"use strict";
const { MongoClient, ObjectID } = require("mongodb");
let Redis = require("ioredis"),
_ = require('lodash'),
uuid = require("uuid"),
Promise = require("bluebird"),
SESSION= require('../enums/session'),
crypto = require("crypto");
const url  = process.env["MONGODB"];//"mongodb://adfin:123456@localhost:27017/?maxPoolSize=20&w=majority";
 
const client = new MongoClient(url) //process.env["MONGODB"]);

let config = {
    secret: Math.random(),
    sessionKey: 'sid',
    sidPrefix: '',
    keepAlive: true,
    expires: 30 * 60 * 1000,
    path: "/",
    httpOnly: true,
    secure: false,
    csrf: {
        able:false,
        key: SESSION.CSRF,
        ignoreMethod: ['GET', 'HEAD', 'OPTIONS'],
        value: function(req){
            return req.param(SESSION.CSRF) || req.headers['csrf-token'] || req.headers['xsrf-token'] || req.headers['x-csrf-token'] || req.headers['x-xsrf-token'];
        }
    },
    redisNodes: [],
    
};

let run = async function() {
    try {
      // Connect the client to the server (optional starting in v4.7)
      await client.connect();
      // Establish and verify connection
      await client.db("admin").command({ ping: 1 });
      console.log("Connected successfully to Mongodb server");
    } finally {
      // Ensures that the client will close when you finish/error
      await client.close();
    }
  }
  

let getDomain = function (req) {
    return req.protocal + '://' + req.headers.host;
}
let generateCSRF = () => {return uuid.v4().replace(/[\/\+\=\-]/g,'')}

let sign = function (val,secret) {
    let now = Date.now()+Math.random()
    secret = crypto.createHash('md5').update(val+'.'+now).digest('hex')
    let secretStr1 = crypto.createHmac('sha1',secret).update(secret+'.'+now).digest('base64'),
        secretStr2 = crypto.createHmac('sha1',secret).update(val+'.'+now).digest('base64'),
        secretStr3 = crypto.createHmac('sha256',secret).update(secretStr1+secretStr2).digest('base64'),
        secretLen = Math.round(secretStr3.length/3)
    let secretStrArr = [secretStr3.substr(0,secretLen),secretStr1,secretStr3.substr(secretLen,secretLen),secretStr2,secretStr3.substr(secretLen*2)];
    return secretStrArr.join('').replace(/[\/\+\=\-]/g, '')
};

let getRandomNode = function(){
    let number = 8 + Math.round(Math.random() *8);
    let node = []
    while(number--){
        node.push(Math.floor(Math.random() *256))
    }
    return node
}
let generate = function (req, res) {
    let session = {};
    let time = new Date().getTime();
    session.id = config.sidPrefix + sign(uuid.v1({
            node: getRandomNode(),
            clockseq: 0x1234,
            msecs: new Date('1991-04-11').getTime(),
            nsecs: 1122
        }), config.secret);
    //session.id = config.sidPrefix + uuid.v4().replace(/[\/\+\=\-]/g, '');
    session.expires = time + config.expires;
    session[config.csrf.key] = generateCSRF();
    req.session = session;
    writeHead(req, res);
    return session;
};
let writeHead = function (req, res) {
    if (req.session && req.session.id) {
        res.cookie(config.sessionKey, req.session.id, {
            httpOnly: config.httpOnly,
            secure: config.secure,
            path: config.path,
            signed: req.secret ? true : false
        });
        // 判断是否开启csrf
        //if (config.csrf.able) {
        //    res.cookie(config.csrf.key, req.session[config.csrf.key], {
        //        httpOnly: config.httpOnly,
        //        secure: config.secure,
        //        path: config.path,
        //        signed: req.secret ? true : false
        //    });
        //}
    }
};
let checkCSRF = function(req, csrf) {
    let referer = req.headers.referer || getDomain(req), method = (req.method || 'get').toUpperCase();
    if (config.csrf.ignoreMethods.indexOf(method) >= 0 || (new RegExp(req.headers.host).test(referer) && config.csrf.value(req) === csrf)) {
        return true;
    } else {
        return false;
    }
};

 
let setCSRFLocal = function(req, res) {
    if (config.csrf.able) {
        res.locals[SESSION.CSRF] = req.csrfToken();
    }
};

 
let redisCluster, notConn = true;
 
let refreshRedis = function(req, res) {
    let session = req.session || {}, id = req.session && req.session.id;
    session.expires = Date.now() + config.expires;
    writeHead(req, res);
    return id ? redisCluster.expire(id, config.expires / 1000) : new Promise((resolve, reject) => {resolve({});});
};

// 保存
let save = function (req, res, next) {
    let id = req.session && req.session.id;
    setCSRFLocal(req, res);
   
    if (!notConn && id) {
        let json = JSON.stringify(req.session);
        redisCluster.hset(id, 'session', json)
            .then(() => { return refreshRedis(req, res) })
            .then(() => { next(); })
            .catch((err) => { console.error(err) && next();});
    } else {
         
        next();
    }
};

 
let reset = function (req, res, next) {
    let id = req.cookies[config.sessionKey];
 
    if (!notConn && id) {
        res.clearCookie(config.sessionKey, {
            httpOnly: config.httpOnly,
            secure: config.secure,
            path: config.path,
            signed: req.secret ? true : false
        });
        //if (config.csrf.able) {
        //    res.clearCookie(config.csrf.key, {
        //        httpOnly: config.httpOnly,
        //        secure: config.secure,
        //        path: config.path,
        //        signed: req.secret ? true : false
        //    });
        //}
        generate(req, res);
        redisCluster.hdel(id, 'session').then(() => { next(); }).catch((err) => { console.error(err) && next(); });
    } else {
        generate(req, res);
        next();
    }
};

let exportsFns = {
    session: (req, res, next) => {
    
        if (!config.redisNodes || config.redisNodes.length == 0) {
            console.info('not redis nodes info');
            return (req, res, next) => { next(); };
        }
       //     console.log(req.cookies)
        req.csrfToken = function() {
            return this.session && this.session[config.csrf.key];
        };
        let id = req.cookies[config.sessionKey];
        if (!id) {
            req.session = generate(req, res);
            save(req, res, next);
        } else if (notConn) {
            console.error('redisCluster not connection!');
            if (!req.session) {
                req.session = generate(req, res);
            }
            save(req, res, next);
        } else {
            redisCluster.hget(id, 'session').then((reply) => {
                let time = Date.now();
                let expires = time + config.expires, session = {id: id, expires: expires};
                session[config.csrf.key] = generateCSRF();
                if (reply) {
                    session = JSON.parse(reply);
                } else {
                    session = generate(req, res);
        
                    var requestType = req.headers['X-Requested-With'] || req.headers['x-requested-with'];
                    if (null !== requestType && "XMLHttpRequest" === requestType && config.csrf.value(req)) {
                        session[config.csrf.key] = config.csrf.value(req);
                    }
                }
                req.session = session;
                setCSRFLocal(req, res);
                // csrf verification
                if (config.csrf.able && session[config.csrf.key] && !checkCSRF(req, session[config.csrf.key])) {
                    var err = new Error('CSRF verification failed, Request aborted.');
                    err.status = 403;
                    next(err);
                } else if (config.keepAlive) {
                    save(req, res, next);
                } else {
                    next();
                }
            }).catch((err) => {
                    console.error(err);
                req.session = generate(req, res);
                save(req, res, next);
            });
        }
    },

    save: (req, res, next) => {
        if (next) {
 
            save(req, res, next);
        } else {
            return new Promise((resolve, reject) => {
                save(req, res, () => { resolve({})});
        });
        }
        reset: reset
    } 

 
};
module.exports = function(options) {
    _.merge(config, options);
    // const clusterOptions={
    //     enableReadyCheck: true,
    //     retryDelayOnClusterDown:300,
    //     retryDelayOnFailover:1000,
    //     retryDelayOnTryAgain:3000,
    //     slotRefreshTimeout:200000000000000,
    //     clusterRetryStrategy: (times)=>Math.min(times*1000,10000),
    //     dnsLookup: (address,callback)=>callback(null,address),
    //     scaleReads: 'slave',
    //     showFriendlyErrorStack: true   

    // }
    run().catch(console.dir);
    
    redisCluster = new Redis.Cluster(config.redisNodes);

    redisCluster.on('error', (err) => { //console.error(err); 
        notConn = true; });
    redisCluster.on('connect', () => { //console.info('redisCluster connect'); 
        notConn = false; });
    exportsFns.Cluster = redisCluster;
    
    // return function(req,res,next) { 
    //     next()
    // } 
   
    return exportsFns;
};