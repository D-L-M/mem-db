import { expect } from 'chai';
import * as request from 'sync-request';
import * as btoa from 'btoa';


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
                "engine": "MemDB",
                "state": "active",
                "version": "0.0.1"
            }
        );

    });


});
