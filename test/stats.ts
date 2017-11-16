import { expect } from 'chai';
import * as request from 'sync-request';


describe('Stats', function()
{


    it('retrieves default values', () =>
    {

        let response = request('GET', 'http://127.0.0.1:9999/_stats');
        let body     = response.getBody().toString('utf8');
        let stats    = JSON.parse(body);

        expect(stats).to.haveOwnProperty('total_documents');
        expect(stats).to.haveOwnProperty('total_inverted_indices');
        expect(stats.total_documents).to.equal(0);
        expect(stats.total_inverted_indices).to.equal(0);

    });


});
