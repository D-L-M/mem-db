import { expect } from 'chai';
import * as request from 'sync-request';


describe('Documents', function()
{


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
        let createdResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/123', {'json': document}).getBody().toString('utf8'));

        expect(createdResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be stored',
                'success': true
            }
        );

        /*
         * Read
         */
        let readResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/123').getBody().toString('utf8'));
        
        expect(readResponse).to.deep.equal(document);

        /*
         * Update
         */
        let updatedResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999/123', {'json': updatedDocument}).getBody().toString('utf8'));
        
        expect(createdResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be stored',
                'success': true
            }
        );

        let updatedReadResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/123').getBody().toString('utf8'));
        
        expect(updatedReadResponse).to.deep.equal(updatedDocument);

        /*
         * Delete
         */
        let deletedResponse = JSON.parse(request('DELETE', 'http://127.0.0.1:9999/123').getBody().toString('utf8'));
        
        expect(deletedResponse).to.deep.equal(
            {
                'id': '123',
                'message': 'Document will be removed',
                'success': true
            }
        );

        try
        {

            request('GET', 'http://127.0.0.1:9999/123').getBody();

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

    });


    it('return an error if requesting non-existent resource', () =>
    {

        try
        {

            request('GET', 'http://127.0.0.1:9999/badId').getBody();

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

            request('PUT', 'http://127.0.0.1:9999/badBody', {'body': '{"bad":"json",}'}).getBody();

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

        let createdResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999', {'json': document}).getBody().toString('utf8'));

        expect(createdResponse.message).to.equal('Document will be stored');
        expect(createdResponse.success).to.be.true;
        expect(createdResponse.id).to.have.lengthOf(36);

        let readResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/' + createdResponse.id).getBody().toString('utf8'));
        
        expect(readResponse).to.deep.equal(document);

        request('DELETE', 'http://127.0.0.1:9999/' + createdResponse.id);

        /*
         * Create another one to ensure the IDs are different
         */
        let anotherCreatedResponse = JSON.parse(request('PUT', 'http://127.0.0.1:9999', {'json': document}).getBody().toString('utf8'));
        
        expect(anotherCreatedResponse.message).to.equal('Document will be stored');
        expect(anotherCreatedResponse.success).to.be.true;
        expect(anotherCreatedResponse.id).to.have.lengthOf(36);
        expect(anotherCreatedResponse.id).to.not.equal(createdResponse.id);

        request('DELETE', 'http://127.0.0.1:9999/' + anotherCreatedResponse.id);

    });


});
