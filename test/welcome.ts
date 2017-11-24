import { expect } from 'chai';
import * as request from 'sync-request';


describe('Welcome message', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all');
    });


    it('displays as expected', () =>
    {

        let statsResponse = JSON.parse(request('GET', 'http://127.0.0.1:9999').getBody().toString('utf8'));

        expect(statsResponse).to.deep.equal(
            {
                "engine": "MemDB",
                "state": "active",
                "version": "0.0.1"
            }
        );

    });


});
