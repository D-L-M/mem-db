import { expect } from 'chai';
import * as request from 'sync-request';
import * as sleep from 'sleep-sync';
import * as btoa from 'btoa';


describe('Documents', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}});
    });


    it('can be created, read, updated and deleted', () =>
    {

        let document =
            {
                'foo': 'bar'
            };

        let updatedDocument =
            {
                'foo': 'baz'
            };

        /*
         * Create
         */
        let createdResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody().toString('utf8'));

        expect(createdResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be stored',
                'success': true
            }
        );

        sleep(500);

        /*
         * Read
         */
        let readResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(readResponse).to.deep.equal(document);

        let replicaReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9998/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(replicaReadResponse).to.deep.equal(document);

        /*
         * Update
         */
        let updatedResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': updatedDocument}).getBody().toString('utf8'));

        expect(createdResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be stored',
                'success': true
            }
        );

        sleep(500);

        let updatedReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(updatedReadResponse).to.deep.equal(updatedDocument);

        let replicaUpdatedReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9997/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(replicaUpdatedReadResponse).to.deep.equal(updatedDocument);

        /*
         * Add a second document to ensure it is not deleted with others
         */
        request('PUT', 'http://127.0.0.1:9999/1234', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document});

        /*
         * Delete
         */
        let deletedResponse = JSON.parse(request('DELETE', 'http://127.0.0.1:9999/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(deletedResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be removed',
                'success': true
            }
        );

        sleep(500);

        try
        {

            request('GET', 'http://127.0.0.1:9997/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let deletedReadResponse = JSON.parse(error.body.toString('utf8'));

            expect(deletedReadResponse).to.deep.equal(
                {
                    'id': '123',
                    'message': 'Document does not exist',
                    'success': false
                }
            );

        }

        try
        {

            request('GET', 'http://127.0.0.1:9998/123', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let replicaDeletedReadResponse = JSON.parse(error.body.toString('utf8'));

            expect(replicaDeletedReadResponse).to.deep.equal(
                {
                    'id': '123',
                    'message': 'Document does not exist',
                    'success': false
                }
            );

        }

        /*
         * Ensure the other document is still there and then delete it
         */
        let secondReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/1234', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(secondReadResponse).to.deep.equal(document);

        let replicaSecondReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9997/1234', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(replicaSecondReadResponse).to.deep.equal(document);

        request('DELETE', 'http://127.0.0.1:9999/1234')

    });


    it('return an error if requesting non-existent resource', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999/badId', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let nonexistentResponse = JSON.parse(error.body.toString('utf8'));

            expect(nonexistentResponse).to.deep.equal(
                {
                    'id': 'badId',
                    'message': 'Document does not exist',
                    'success': false
                }
            );

        }

    });


    it('return an error if malformed JSON is supplied', () =>
    {

        try
        {

            request('PUT', 'http://127.0.0.1:9999/badBody', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'body': '{"bad":"json",}'}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let badDocumentResponse = JSON.parse(error.body.toString('utf8'));

            expect(badDocumentResponse).to.deep.equal(
                {
                    'id': 'badBody',
                    'message': 'Document is not valid JSON',
                    'success': false
                }
            );

        }

    });


    it('will be assigned a random ID if one is not provided', () =>
    {

        let document =
            {
                'foo': 'bar'
            };

        let createdResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody().toString('utf8'));

        expect(createdResponse.message).to.equal('Document will be stored');
        expect(createdResponse.success).to.be.true;
        expect(createdResponse.id).to.have.lengthOf(36);

        sleep(500);

        let readResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/' + createdResponse.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(readResponse).to.deep.equal(document);

        request('DELETE', 'http://127.0.0.1:9999/' + createdResponse.id);

        sleep(500);

        /*
         * Create another one to ensure the IDs are different
         */
        let anotherCreatedResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document}).getBody().toString('utf8'));

        expect(anotherCreatedResponse.message).to.equal('Document will be stored');
        expect(anotherCreatedResponse.success).to.be.true;
        expect(anotherCreatedResponse.id).to.have.lengthOf(36);
        expect(anotherCreatedResponse.id).to.not.equal(createdResponse.id);

        request('DELETE', 'http://127.0.0.1:9999/' + anotherCreatedResponse.id);

        sleep(500);

    });


    it('can be truncated', () =>
    {

        let document =
            {
                'foo': 'bar'
            };

        /*
         * Create some documents
         */
        for (let i = 0; i < 10; i++) {
            request('PUT', 'http://127.0.0.1:9999/', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document});
        }

        sleep(500);

        /*
         * Delete everything
         */
        let deletedResponse = JSON.parse(request('DELETE', 'http://127.0.0.1:9999/_all', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(deletedResponse).to.deep.equal(
            {
                'message': 'All documents will be removed',
                'success': true
            }
        );

        sleep(500);

        /*
         * Search for all documents
         */
        let allResponses = JSON.parse(request('GET', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(allResponses.results.length).to.equal(0);

    });


});
