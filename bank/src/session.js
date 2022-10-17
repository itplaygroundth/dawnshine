const Client = require('./client');
const EventEmitter = require('events');

module.exports = class Session extends EventEmitter {
    state;
    accountNumber;
    accountType;
    pin;
    _client;
    _promise;

    constructor(state, accountNumber, accountType, pin) {
        super();
        this.state = state;
        this.accountNumber = accountNumber;
        this.accountType = accountType;
        this.pin = pin;

        this.reset();
    }

    async sync(callback) {
        if (this._promise) {
            try {
                await this._promise;
            } catch (e) {
            }
        }

        return this._promise = callback();
    }

    async _refresh() {
        this._client = new Client(this.state);
        await this._client.init();
        await this._client.verifyPin(this.pin);

        this.emit('STATE_UPDATED', this.state);
    }

    async call(callback) {
        return this.sync(async () => {
            if (!this._client) {
                await this._refresh();
            }

            try {
                return await callback(this._client);
            } catch (e) {
                if (e.message.includes('ERR3006')) {
                    await this._refresh();

                    return await callback(this._client);
                } else {
                    throw e;
                }
            }
        });
    }

    reset() {
        this._activityList = [];
        this._bankInfoList = [];
        this._bankRefreshedAt = 0;
        this._inquireForTransferMoneyList = [];
    }

    async getBankInfoList() {
        const now = Date.now();
        if (now > this._bankRefreshedAt + (30 * 1000)) {
            this._bankInfoList = await this.call(client => client.getBankInfoList());
            this._bankRefreshedAt = now;
        }

        return this._bankInfoList;
    }
}