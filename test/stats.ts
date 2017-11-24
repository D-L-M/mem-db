import { expect } from 'chai';
import * as request from 'sync-request';
import * as sleep from 'sleep-sync';


describe('Stats', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all');
    });


    it('retrieves default values', () =>
    {

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats').getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'totals':
                {
                    'documents': 0,
                    'inverted_indices': 0
                }
            }
        );

    });


    it('sees correct values when documents are indexed', () =>
    {

        request('PUT', 'http://127.0.0.1:9999/321', {'json': {'foo': 'bar baz', 'success': true}});

        sleep(500);

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats').getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'totals':
                {
                    'documents': 1,
                    'inverted_indices': 4
                }
            }
        );

        request('DELETE', 'http://127.0.0.1:9999/321');

        sleep(500);

    });


});
