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
                'message': 'Document /123 will be stored',
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
                'message': 'Document /123 will be stored',
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
                'message': 'Document /123 will be removed',
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
                    'message': 'Document is not valid JSON',
                    'success': false
                }
            );

        }

    });


});
