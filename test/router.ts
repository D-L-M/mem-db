import { expect } from 'chai';
import * as request from 'sync-request';


describe('Router', function()
{


    this.timeout(5000);


    /*
     * Truncate the database
     */
    beforeEach(() =>
    {
        request('DELETE', 'http://127.0.0.1:9999/_all');
    });


    it('displays an error message if no matches are found', () =>
    {

        try
        {

            request('POST', 'http://127.0.0.1:9999/_bad_route').getBody();

            expect(true).to.equal(false);

        }

        catch (error)
        {

            let noRouteResponse = JSON.parse(error.body.toString('utf8'));

            expect(noRouteResponse).to.deep.equal(
                {
                    'message': 'Unknown request',
                    'success': false
                }
            );

        }

    });


});
