import { expect } from 'chai';
import * as request from 'sync-request';
import * as btoa from 'btoa';


describe('Welcome message', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}});
    });


    it('displays as expected', () =>
    {

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999', {'headers': {'Authorization': 'Basic ' + btoa('root:password')}}).getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                'engine': 'MemDB',
                'state': 'active',
                'version': '0.0.1'
            }
        );

    });


});
