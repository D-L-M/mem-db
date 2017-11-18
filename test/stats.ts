import { expect } from 'chai';
import * as request from 'sync-request';


describe('Stats', function()
{


    it('retrieves default values', () =>
    {

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats').getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'total_documents': 0,
                'total_inverted_indices': 0
            }
        );

    });


    it('sees correct values when documents are indexed', () =>
    {

        request('PUT', 'http://127.0.0.1:9999/321', {'json': {'foo': 'bar baz', 'success': true}});

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats').getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'total_documents': 1,
                'total_inverted_indices': 4
            }
        );

        request('DELETE', 'http://127.0.0.1:9999/321');

    });


});
