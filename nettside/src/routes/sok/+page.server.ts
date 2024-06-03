import prisma from '$lib/prisma';
import type { ExtendedProduct } from '$lib/types/extendedPrisma';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url }) => {
	if (url.searchParams.get('search')) {
		return {
			param: url.searchParams.get('search'),
			products: (await prisma.products.findMany({
				where: {
					title: {
						contains: url.searchParams.get('search') as string,
						mode: 'insensitive'
					}
				}
			})) as ExtendedProduct[]
		};
	}
};
