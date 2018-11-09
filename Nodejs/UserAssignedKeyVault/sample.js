var rp = require('request-promise');
var KeyVault = require('azure-keyvault');
var keyVaultName = process.env.KEY_VAULT_NAME;
var secretName = process.env.SECRET_NAME;
var secretVersion = process.env.SECRET_VERSION;

var options = {
    uri: 'http://169.254.169.254/metadata/identity/oauth2/token',
    qs: {
        'api-version': '2018-02-01',
        'resource': 'https://vault.azure.net',
    },
    headers: {
        'Metadata': 'true'
    },
    json: true // Automatically parses the JSON string in the response
};

rp(options).then(function (tokenResponse) {
    var authenticator = function (challenge, callback) {
        var authorizationValue = tokenResponse.token_type + ' ' + tokenResponse.access_token;
        return callback(null, authorizationValue);
    };

    var credentials = new KeyVault.KeyVaultCredentials(authenticator);
    var client = new KeyVault.KeyVaultClient(credentials);

    client.getSecret(`https://${keyVaultName}.vault.azure.net`, secretName, secretVersion).then((secretBundle) => {
        console.log(`The secret value is: ${secretBundle.value}`)
    }).catch((err) => {
        console.log(err)
    });
}).catch(function (err) {
    console.log(err)
});