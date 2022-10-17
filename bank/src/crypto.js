const {webcrypto} = require('crypto');

module.exports = class CryptoMessage {
    constructor(secretKey, iv) {
        this.secretKey = secretKey;
        this.iv = iv;
    }

    async encrypt(message) {
        const encrypted = await webcrypto.subtle.encrypt(
            {
                name: "AES-CBC",
                iv: this.iv
            },
            this.secretKey,
            Buffer.from(message)
        );

        return Buffer.from(encrypted).toString('base64');
    }

    async decrypt(message) {
        const decrypted = await webcrypto.subtle.decrypt(
            {
                name: "AES-CBC",
                iv: this.iv
            },
            this.secretKey,
            Buffer.from(message, 'base64')
        );

        return Buffer.from(decrypted).toString();
    }

    static generateKey() {
        return webcrypto.subtle.generateKey(
            {
                name: "ECDH",
                namedCurve: "P-256"
            },
            false,
            ["deriveKey"]
        );
    }

    static deriveSecretKey(privateKey, publicKey) {
        return webcrypto.subtle.deriveKey(
            {
                name: "ECDH",
                public: publicKey
            },
            privateKey,
            {
                name: "AES-CBC",
                length: 256
            },
            true,
            ["encrypt", "decrypt"]
        );
    }

    static async exportPublicKey(publicKey) {
        return Buffer.from(await webcrypto.subtle.exportKey('spki', publicKey)).toString('hex');
    }

    static importPublicKey(publicKey) {
        return webcrypto.subtle.importKey(
            'spki',
            Buffer.from(publicKey, 'hex'),
            {
                name: "ECDH",
                namedCurve: "P-256"
            },
            false,
            ["deriveKey"],
        );
    }
}
