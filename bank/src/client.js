const axios = require('axios');
const { CookieJar } = require('tough-cookie');
const { wrapper } = require('axios-cookiejar-support');
const uuid = require('uuid');
const CryptoMessage = require("./crypto");
const dayjs = require("dayjs");
const { KBANK } = require("./banks");
const { keyBy } = require("lodash/collection");

module.exports = class Client {
    encrypter;
    client;

    constructor(state = {}) {
        this.state = state;
        this.client = wrapper(axios.create({
            jar: new CookieJar(),
        }));
    }

    async init() {
        this.sessionId = uuid.v4().replaceAll('-', '');

        const exchangeResponse = await this._exchange();

        await this._receiveMobileNo(exchangeResponse['t1']);

        const checkAuthenIDAndProfileResponse = await this._checkAuthenIDAndProfile();

        this.sessionId = checkAuthenIDAndProfileResponse['sessionId'];
        this.hashV = checkAuthenIDAndProfileResponse['hashV'];

        if (checkAuthenIDAndProfileResponse['getProfileStatusResult']['quickBalanceFlag'] === 'Y') {
            this.callMobileUtilityInquireQuickBalanceHome(this.sessionId, this.hashV); // fake
        }
    }

    async _exchange() {
        const keyPair = await CryptoMessage.generateKey();

        const response = await this.callSecurityExchangeKeyAndConfigV2(this.sessionId, {
            'appName': 'KPLUS_Victoria',
            'appVersion': '5.15.3',
            'configModifyDateTime': this.state['configModifyDateTime'] ?? '20220119084252603',
            'dynamicMenuModifyDateTime': this.state['dynamicMenuModifyDateTime'] ?? '20210522231551703',
            'forceUpdatedDate': '20180420160313',
            'labelModifyDateTime': this.state['labelModifyDateTime'] ?? '20220118182153423',
            'language': 'T',
            'menuModifyDateTime': this.state['menuModifyDateTime'] ?? '20220118182153423',
            'messageModifyDateTime': this.state['messageModifyDateTime'] ?? '20220118182153423',
            'osVersion': '11',
            'platform': 'android',
            'publicKeyA': await CryptoMessage.exportPublicKey(keyPair.publicKey),
        });

        [
            'configModifyDateTime',
            'dynamicMenuModifyDateTime',
            'labelModifyDateTime',
            'menuModifyDateTime',
            'messageModifyDateTime'
        ].forEach(key => {
            const value = response[key] ?? null;
            if (value) {
                this.state[key] = value;
            }
        }
        );

        const secretKey = await CryptoMessage.deriveSecretKey(keyPair.privateKey, await CryptoMessage.importPublicKey(response['publicKeyC']));

        this.encrypter = new CryptoMessage(secretKey, Buffer.from(this.sessionId.substring(6, 22)));

        return response;
    }

    _receiveMobileNo(t1) {
        return this.callSecurityReceiveMobileNo(this.sessionId, {
            'networkType': '2',
            't1': t1,
        });
    }

    _checkAuthenIDAndProfile(hashV) {
        return this.callSecurityCheckAuthenIDAndProfile(this.sessionId, hashV, {
            'modelName': 'M2007J20CT',
            'db1': this.state['db1'],
            'dka3': this.state['dka3'],
            'dm1': this.state['dm1'],
            'wifiKey': this.state['wifiKey'],
            'modelType': '2',
            'token': this.state['token'],
            'menuModifyDateTime': this.state['menuModifyDateTime'] ?? '20220118182153423',
            'detectDetail': '',
        });
    }

    async verifyPin(pin) {
        let verifyPinResponse = await this.callSecurityVerifyPin(this.sessionId, {
            'pin': pin,
            'dm1': this.state['dm1'],
            'dka3': this.state['dka3'],
            'channelDetail': '',
            'latitude': '0',
            'longitude': '0',
            'listAccountFlag': 'Y',
            'needUpdateSessionFlag': 'Y',
        });

        this.state['db1'] = verifyPinResponse['db1'];

        this.loggedInSessionId = verifyPinResponse['sessionId'];

        return verifyPinResponse;
    }

    getInquiryAccountBalance(accountNo, accountType) {
        return this.callMobileUtilityInquiryAccountBalance(this.loggedInSessionId, {
            'accountNo': accountNo,
            'accountType': accountType,
        });
    }

    getAccountActivityList(accountNumber, startRecord = 1, numberOfRecord = 50, nextPageId = null, subtractMouth = 0) {
        const date = dayjs().subtract(subtractMouth, 'month');

        return this.callAccountActivityListAccountActivityV2(this.loggedInSessionId, {
            'accountNumber': accountNumber,
            'startRecord': startRecord,
            'numberOfRecord': numberOfRecord,
            'startDate': date.startOf('month').format('DDMMYYYY'),
            'endDate': date.endOf('month').format('DDMMYYYY'),
            'nextPageId': nextPageId,
        });
    }

    getAccountActivityDetail(accountNumber, activity) {
        return this.callMobileUtilityShowAccountActivityDetail(this.loggedInSessionId, {
            'accountNumber': accountNumber,
            'transactionNumber': activity['transactionNumber'],
            'rqUid': activity['rqUid'],
            'transactionUxDate': activity['transactionUxDate'],
            'sourceSystemId': activity['sourceSystemId'],
        });
    }

    async inquireForTransferMoney(handle, fromAccount, toAccount, amount, transferType, targetBankCode, sourceBankCode = KBANK) {
        const data = await this.callMoneyTransferInquireForTransferMoneyORFT(this.loggedInSessionId, {
            'fromAccount': fromAccount,
            'toAccount': toAccount,
            'amount': amount,
            'transferType': transferType,
            'qrOrigin': '',
            'rtpId': '',
            'targetBankCode': targetBankCode,
            'sourceBankCode': sourceBankCode,
        });
        handle['kbankInternalSessionId'] = data['kbankInternalSessionId'];
        handle['fromAccount'] = fromAccount;
        handle['toAccount'] = toAccount;
        handle['amount'] = amount;
        handle['transferType'] = transferType;
        handle['targetBankCode'] = targetBankCode;
        handle['sourceBankCode'] = sourceBankCode;

        return data;
    }

    transferMoney(handle) {
        return this.callMoneyTransferTransferMoneyORFT(this.loggedInSessionId, {
            'kbankInternalSessionId': handle['kbankInternalSessionId'],
            'transferTemplateId': '',
            'latitude': '0',
            'longitude': '0',
            'additionalNote': handle['additionalNote'] ?? '',
            'sendAdditionalNoteFlag': 'N',
            'categoryCode': handle['categoryCode'] ?? '99',
            'fromAccount': handle['fromAccount'],
            'toAccount': handle['toAccount'],
            'amount': handle['amount'],
            'transferType': handle['transferType'],
            'qrOrigin': '',
            'rtpId': '',
            'targetBankCode': handle['targetBankCode'],
            'sourceBankCode': handle['sourceBankCode'],
        });
    }

    async getBankInfoList() {
        const response = await this.callMobileUtilityListBankDestination(this.loggedInSessionId, {
            'startRecord': '1',
            'numberOfRecord': '-1',
        });

        return keyBy(response['bankInfoList'], v => v.targetBankCode);
    }

    scanQr(rawQrBarcode, rawType = 'QR', qrOrigin = 'QUICKPAY') {
        return this.callPaymentUtilityScanQrBarcode(this.loggedInSessionId ?? this.sessionId, {
            'rawQrBarcode': rawQrBarcode,
            'rawType': rawType,
            'qrOrigin': qrOrigin,
            'latitude': '0',
            'longitude': '0',
        });
    }

    callSecurityExchangeKeyAndConfigV2(sessionId, data) {
        return this.call({
            path: 'security/exchangeKeyAndConfigV2',
            command: 'SECURITY',
            data,
            sessionId,
        });
    }

    callSecurityReceiveMobileNo(sessionId, data) {
        return this.call({
            https: false,
            path: 'security/receiveMobileNo',
            command: 'SECURITY',
            data,
            sessionId,
        });
    }

    callSecurityCheckAuthenIDAndProfile(sessionId, hashV, data) {
        return this.call({
            path: 'security/checkAuthenIDAndProfile',
            command: 'SECURITY',
            hashV,
            sessionId,
            data,
        });
    }

    callSecurityVerifyPin(sessionId, data) {
        return this.call({
            path: 'security/verifyPin',
            command: 'FOOTER_BANKING',
            sessionId,
            data,
        });
    }

    callMobileUtilityInquireQuickBalanceHome(sessionId, hashV) {
        return this.call({
            path: 'mobileUtility/inquireQuickBalanceHome',
            command: 'SECURITY',
            hashV,
            sessionId,
        });
    }

    callMobileUtilityInquiryAccountBalance(sessionId, data) {
        return this.call({
            path: 'mobileUtility/inquiryAccountBalance',
            command: 'FOOTER_BANKING',
            sessionId,
            data,
        });
    }

    callAccountActivityListAccountActivityV2(sessionId, data) {
        return this.call({
            path: 'accountActivity/listAccountActivityV2',
            command: 'FOOTER_BANKING',
            sessionId,
            data,
        });
    }

    callMobileUtilityShowAccountActivityDetail(sessionId, data) {
        return this.call({
            path: 'mobileUtility/showAccountActivityDetail',
            command: 'FOOTER_BANKING',
            sessionId,
            data,
        });
    }

    callMobileUtilityListBankDestination(sessionId, data) {
        return this.call({
            path: 'mobileUtility/listBankDestination',
            command: 'TRF_OTHER',
            sessionId,
            data,
        });
    }

    callMoneyTransferInquireForTransferMoneyORFT(sessionId, data) {
        return this.call({
            path: 'moneyTransfer/inquireForTransferMoneyORFT',
            command: 'TRF_OTHER',
            sessionId,
            data,
        });
    }

    callMoneyTransferTransferMoneyORFT(sessionId, data) {
        return this.call({
            path: 'moneyTransfer/transferMoneyORFT',
            command: 'TRF_OTHER',
            sessionId,
            data,
        });
    }

    callPaymentUtilityScanQrBarcode(sessionId, data) {
        return this.call({
            path: 'paymentUtility/scanQrBarcode',
            command: 'DUMMY_SCAN_BARCODE',
            sessionId,
            data,
        });
    }

    async call({ method = "POST", path, https = true, data = {}, command, sessionId, hashV = '' }) {
        if (this.encrypter) {
            data = await this.encrypter.encrypt(JSON.stringify(data));
        }

        console.log(path);
        const response = await this.client.request({
            method,
            url: `${https ? 'https' : 'http'}://rt10.kasikornbank.com/kplus-service/${path}`,
            data: {
                requestBody: {
                    clientData: data,
                },
                requestHeader: {
                    command,
                    hashV,
                    sessionId
                }
            }
        });

        return this._processResponse(response);
    }

    async _processResponse(response) {
        const data = this.encrypter && typeof response.data === "string" ? JSON.parse(await this.encrypter.decrypt(response.data)) : response.data;
        const responseHeader = data['responseHeader'];

        if (responseHeader['status'] !== 'S') {
            console.log(responseHeader);
            throw new Error(`${responseHeader['messageCode'] ?? ''}, ${responseHeader['displayText'] ?? ''}`);
        }

        return data['responseBody'] ?? [];
    }
}