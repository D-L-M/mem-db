import { expect } from 'chai';
import * as request from 'sync-request';
import * as sleep from 'sleep-sync';
import * as btoa from 'btoa';


describe('Stats', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}});
    });


    it('retrieves default values', () =>
    {

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(statsResponse.totals).to.deep.equal(
            {
                'documents': 0,
                'inverted_indices': 0
            }
        );

        expect(statsResponse.peers.length).to.equal(2);

    });


    it('sees correct values when documents are indexed', () =>
    {

        request('PUT', 'http://127.0.0.1:9999/321', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}, 'json': {'foo': 'bar baz', 'success': true}});

        sleep(500);

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999/_stats', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(statsResponse.totals).to.deep.equal(
            {
                'documents': 1,
                'inverted_indices': 5
            }
        );

        expect(statsResponse.peers.length).to.equal(2);

        request('DELETE', 'http://127.0.0.1:9999/321', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}});

        sleep(500);

    });


});
