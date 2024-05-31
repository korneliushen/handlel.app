import prisma from '$lib/prisma';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url }) => {
	if (url.searchParams.get('search')) {
		return {
      param: url.searchParams.get('search'),
			products: await prisma.products.findMany({
				where: {
					title: {
						contains: url.searchParams.get('search') as string
					}
				}
			})
		};
	}
};
