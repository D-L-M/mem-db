import { expect } from 'chai';
import * as request from 'sync-request';
import * as sleep from 'sleep-sync';
import * as btoa from 'btoa';


var documents =
    [
        {
            'id': '1',
            'document':
            {
                'id': 1,
                'name':
                    {
                        'first': 'John',
                        'last': 'Doe'
                    },
                'age': 30,
                'interests': ['surfing', 'football']
            }
        },
        {
            'id': '2',
            'document':
            {
                'id': 2,
                'name':
                    {
                        'first': 'Jane',
                        'last': 'Doe'
                    },
                'age': 32,
                'interests': ['music', 'Cryptography']
            }
        },
        {
            'id': '3',
            'document':
            {
                'id': 3,
                'name':
                    {
                        'first': 'Jean Paul',
                        'last': 'Smith'
                    },
                'age': 27,
                'interests': ['football games', 'painting']
            }
        }
    ];


describe('Search', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}});
    });


    it('returns returns an error if malformed JSON is supplied', () =>
    {

        try
        {

            request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'body': '{"bad":"json",}'}).getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let badCriteriaResponse = JSON.parse(error.body.toString('utf8'));

            expect(badCriteriaResponse).to.deep.equal(
                {
                    'id': '',
                    'message': 'Search criteria is not valid JSON',
                    'success': false
                }
            );

        }

    });


    it('returns all documents when no criteria set', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * Search for all
         */
        let allResponses = JSON.parse(request('GET', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(allResponses.results.length).to.equal(3);
        expect(allResponses.criteria).to.deep.equal({});
        expect(allResponses.information.total_matches).to.equal(3);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('returns all documents when empty criteria set', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * Search for all
         */
        let allResponses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': {}}).getBody().toString('utf8'));

        expect(allResponses.results.length).to.equal(3);
        expect(allResponses.criteria).to.deep.equal({});
        expect(allResponses.information.total_matches).to.equal(3);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('returns documents with a simple criterion', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * Search for one
         */
        let criteria =
            {
                'and':
                    [
                        {'equals': {'age': 32}}
                    ]
            };

        let responses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(responses.results).to.deep.equal([documents[1]]);
        expect(responses.criteria).to.deep.equal(criteria);
        expect(responses.information.total_matches).to.equal(1);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('returns documents with AND criteria', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * AND search
         */
        let criteria =
            {
                'and':
                    [
                        {'equals': {'age': 32}},
                        {'equals': {'name.first': "jane"}}
                    ]
            };

        let responses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(responses.results).to.deep.equal([documents[1]]);
        expect(responses.criteria).to.deep.equal(criteria);
        expect(responses.information.total_matches).to.equal(1);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('returns documents with OR criteria', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * OR search
         */
        let criteria =
            {
                'OR':
                    [
                        {'equals': {'age': 32}},
                        {'equals': {'name.first': "john"}}
                    ]
            };

        let responses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(responses.results[0]).to.deep.equal(documents[1]);
        expect(responses.results[1]).to.deep.equal(documents[0]);
        expect(responses.criteria).to.deep.equal(criteria);
        expect(responses.information.total_matches).to.equal(2);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('returns documents with nested criteria', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * Nested search
         */
        let criteria =
            {
                'OR':
                    [
                        {'contains': {'interests': "football"}},
                        {
                            'AND':
                                [
                                    {'equals': {'name.first': "John"}},
                                    {'equals': {'name.last': "DOE"}}
                                ]
                        }
                    ]
            };

        let responses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(responses.results[0]).to.deep.equal(documents[0]);
        expect(responses.results[1]).to.deep.equal(documents[2]);
        expect(responses.criteria).to.deep.equal(criteria);
        expect(responses.information.total_matches).to.equal(2);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


    it('can bulk delete documents', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': document.document})
        });

        sleep(500);

        /*
         * Delete documents by criteria
         */
        let criteria =
            {
                'OR':
                    [
                        {'contains': {'interests': "football"}},
                        {
                            'AND':
                                [
                                    {'equals': {'name.first': "John"}},
                                    {'equals': {'name.last': "DOE"}}
                                ]
                        }
                    ]
            };

        let deletionRequest = JSON.parse(request('POST', 'http://127.0.0.1:9999/_delete', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(deletionRequest).to.deep.equal(
            {
                'message': '2 document(s) will be removed',
                'success': true
            });

        sleep(500);

        /*
         * Check that the documents have been removed
         */
        let deletedResponses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': criteria}).getBody().toString('utf8'));

        expect(deletedResponses.results.length).to.equal(0);
        expect(deletedResponses.information.total_matches).to.equal(0);

        /*
         * Check that the correct records remain
         */
        let allResponses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': {}}).getBody().toString('utf8'));

        expect(allResponses.results).to.deep.equal([documents[1]]);
        expect(allResponses.criteria).to.deep.equal({});
        expect(allResponses.information.total_matches).to.equal(1);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id, {'headers': {'Authorization': 'Basic ' + btoa('root:password')}})
        });

        sleep(500);

    });


});
