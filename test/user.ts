import { expect } from 'chai';
import * as request from 'sync-request';
import * as btoa from 'btoa';
import * as sleep from 'sleep-sync';


describe('User management', function()
{


    this.timeout(5000);


    it('can update a user\'s password', () =>
    {

        /*
         * Update the password
         */
        let document       = {'username': 'root', 'password': 'password2', 'action': 'update'};
        let updateResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody().toString('utf8'));

        expect(updateResponse).to.deep.equal(
            {
                'success': true,
                'message': 'User will be created or updated'
            }
        );

        sleep(250);

        /*
         * Ensure that the old password no longer works
         */
        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody();
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

        /*
         * Ensure that the new password works
         */
        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password2')}}).getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'engine': 'MemDB',
                'state': 'active',
                'version': '0.0.1'
            }
        );

        /*
         * Reset the password
         */
        let resetDocument = {'username': 'root', 'password': 'password', 'action': 'update'};

        request('PUT', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password2')}, 'json': resetDocument});

        sleep(250);

    });


    it('can create and delete a new user', function()
    {

        /*
         * Create the user
         */
        let document       = {'username': 'foo', 'password': 'bar', 'action': 'create'};
        let createResponse = JSON.parse(request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody().toString('utf8'));

        expect(createResponse).to.deep.equal(
            {
                'success': true,
                'message': 'User will be created or updated'
            }
        );

        sleep(250);

        /*
         * Ensure that the new user's credentials work
         */
        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('foo:bar')}}).getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'engine': 'MemDB',
                'state': 'active',
                'version': '0.0.1'
            }
        );

        /*
         * Delete the user
         */
        let deleteDocument = {'username': 'foo', 'action': 'delete'};
        let deleteResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': deleteDocument}).getBody().toString('utf8'));

        expect(deleteResponse).to.deep.equal(
            {
                'success': true,
                'message': 'User will be deleted'
            }
        );

        sleep(250);

        /*
         * Ensure that the user's credentials no longer work
         */
        try
        {
            request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('foo:bar')}}).getBody();
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


    it('is only accessible by the root user', function()
    {

        /*
         * Create a new user
         */
        let document = {'username': 'foo', 'password': 'bar', 'action': 'create'};

        request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document});

        sleep(250);

        /*
         * Ensure that the user's credentials do not allow access to the user
         * management endpoint
         */
        try
        {
            request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('foo:bar')}}).getBody();
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

        /*
         * Delete the user
         */
        let deleteDocument = {'username': 'foo', 'action': 'delete'};

        request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': deleteDocument});

    });


    it('fails with a malformed payload', function()
    {

        try
        {
            request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'body': '{"bad":"json",}'}).getBody();
        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Malformed request',
                    'success': false
                }
            );

        }

    });


    it('fails to update with a missing username', function()
    {

        try
        {

            let document = {'password': 'foo', 'action': 'update'};

            request('PUT', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody();

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Malformed request',
                    'success': false
                }
            );

        }

    });


    it('fails to update with a missing password', function()
    {

        try
        {

            let document = {'username': 'foo', 'action': 'create'};

            request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody();

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Malformed request',
                    'success': false
                }
            );

        }

    });


    it('fails to delete with a missing username', function()
    {

        try
        {

            let document = {'action': 'delete'};

            request('POST', 'http://127.0.0.1:9999/_user', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody();

        }

        catch (error)
        {

            let unauthedResponse = JSON.parse(error.body.toString('utf8'));

            expect(unauthedResponse).to.deep.equal(
                {
                    'message': 'Malformed request',
                    'success': false
                }
            );

        }

    });


});
