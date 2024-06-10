import prisma from '$lib/prisma';
import type { PageServerLoad } from './$types';
import algoliasearch from 'algoliasearch';

const client = algoliasearch('AA8FDXU3JW', '5ebf3bd5ba51b5d6ce63cfe54ce78985');
const index = client.initIndex('test');

export const load: PageServerLoad = async ({ url }) => {
	if (url.searchParams.get('search')) {
		return {
      param: url.searchParams.get('search'),
			products: await index.search(url.searchParams.get('search') || "", {
				attributesToRetrieve: ['*'],
				hitsPerPage: 20,
			  }).then(({ hits }) => {
				console.log(hits)
				return hits
			  })
		};
	}
};
