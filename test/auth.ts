import { expect } from 'chai';
import * as request from 'sync-request';
import * as btoa from 'btoa';
import * as fs from 'fs';
import * as os from 'os';
import * as crypto from 'crypto';


describe('Authentication', function()
{


    this.timeout(5000);


    it('fails with no authorisation headers', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999').getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with non-Basic authorisation', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Bearer ' + btoa('root:badpassword')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with the incorrect password', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:badpassword')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with the incorrect username', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('badusername:password')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with the incorrect username and password', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('badusername:badpassword')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('succeeds with the correct username and password', () =>
    {

        let authedResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(authedResponse).to.deep.equal(
            {
                'engine': 'MemDB',
                'state': 'active',
                'version': '0.0.1'
            }
        );

    });


    it('fails with missing HMAC hash', function()
    {

        let secretKey = fs.readFileSync(os.homedir() + '/.memdb/.key');
        let body      = '';
        let hmacNonce = Date.now();
        let hmac      = crypto.createHmac('sha512', secretKey);
        let hmacAuth  = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');

        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-nonce': hmacNonce}}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with missing HMAC nonce', function()
    {

        let secretKey = fs.readFileSync(os.homedir() + '/.memdb/.key');
        let body      = '';
        let hmacNonce = Date.now();
        let hmac      = crypto.createHmac('sha512', secretKey);
        let hmacAuth  = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');

        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-auth': hmacAuth}}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with incorrect HMAC hash', function()
    {

        let secretKey      = fs.readFileSync(os.homedir() + '/.memdb/.key');
        let body           = '';
        let hmacNonce      = Date.now();
        let hmac           = crypto.createHmac('sha512', secretKey);
        let hmacAuth       = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');
        let badHmacAuth    = hmacAuth + 'x';

        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-auth': badHmacAuth, 'x-hmac-nonce': hmacNonce}}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with incorrect HMAC nonce', function()
    {

        let secretKey = fs.readFileSync(os.homedir() + '/.memdb/.key');
        let body      = '';
        let hmacNonce = Date.now();
        let hmac      = crypto.createHmac('sha512', secretKey);
        let hmacAuth  = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');
        let badNonce  = hmacNonce + 'x';

        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-auth': hmacAuth, 'x-hmac-nonce': badNonce}}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('fails with incorrect secret key', function()
    {

        let secretKey = 'abcdefghijklmnopqrstuvwxyz012345';
        let body      = '';
        let hmacNonce = Date.now();
        let hmac      = crypto.createHmac('sha512', secretKey);
        let hmacAuth  = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');

        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-auth': hmacAuth, 'x-hmac-nonce': hmacNonce}}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Not authorised',
                    'success': false
                }
            );

        }

    });


    it('works with HMAC message headers', function()
    {

        let secretKey = fs.readFileSync(os.homedir() + '/.memdb/.key');
        let body      = '';
        let hmacNonce = Date.now();
        let hmac      = crypto.createHmac('sha512', secretKey);
        let hmacAuth  = hmac.update(new Buffer(body + hmacNonce, 'utf-8')).digest('hex');

        let authedResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999', {'headers': {'x-hmac-auth': hmacAuth, 'x-hmac-nonce': hmacNonce}}).getBody().toString('utf8'));

        expect(authedResponse).to.deep.equal(
            {
                'engine': 'MemDB',
                'state': 'active',
                'version': '0.0.1'
            }
        );

    });


});
