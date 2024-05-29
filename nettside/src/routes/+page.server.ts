import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
    if (url.searchParams.get("search")) {
        return {
            products: await prisma.products.findMany({
                where: {
                  title: {
                    search: url.searchParams.get("search") as string,
                  },
                },
              })
        }
    }
    return {
        products: await prisma.products.findMany(
            {
                take: 20
            }
        )
    }
};
