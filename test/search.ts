import { expect } from 'chai';
import * as request from 'sync-request';


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


    it('returns returns an error if malformed JSON is supplied', () =>
    {

        try
        {

            request('POST', 'http://127.0.0.1:9999/_search', {'body': '{"bad":"json",}'}).getBody();

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
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'json': document.document})
        });

        /*
         * Search for all
         */
        let allResponses = JSON.parse(request('GET', 'http://127.0.0.1:9999/_search').getBody().toString('utf8'));

        expect(allResponses.results).to.deep.equal(documents);
        expect(allResponses.criteria).to.deep.equal({});
        expect(allResponses.information.total_matches).to.equal(3);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id)
        });

    });


    it('returns all documents when empty criteria set', () =>
    {

        /*
         * Create documents
         */
        documents.forEach((document) =>
        {
            request('PUT', 'http://127.0.0.1:9999/' + document.id, {'json': document.document})
        });

        /*
         * Search for all
         */
        let allResponses = JSON.parse(request('POST', 'http://127.0.0.1:9999/_search', {'json': {}}).getBody().toString('utf8'));

        expect(allResponses.results).to.deep.equal(documents);
        expect(allResponses.criteria).to.deep.equal({});
        expect(allResponses.information.total_matches).to.equal(3);

        /*
         * Remove documents
         */
        documents.forEach((document) =>
        {
            request('DELETE', 'http://127.0.0.1:9999/' + document.id)
        });

    });


});
