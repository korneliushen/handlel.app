import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
    return {
        products: await prisma.products.findMany(
            {
                take: 20
            }
        )
    }
        
};
