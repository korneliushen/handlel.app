import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { ExtendedProduct } from "$lib/types/extendedPrisma";

export const load: PageServerLoad = () => {
  return {
    streamed: {
      products: new Promise<{ category: string; products: ExtendedProduct[] }[]>(async (resolve, reject) => {
        try {
          const categories: { category: string }[] = await prisma.products.findMany({
            select: {
              category: true,
            },
            distinct: ['category']
          });

          const groupedProducts: { category: string; products: ExtendedProduct[] }[] = [];
          for (const { category } of categories) {
            const products = await prisma.products.findMany({
              where: { category },
              take: 4,
            }) as ExtendedProduct[];
            groupedProducts.push({ category, products });
          }

          return resolve(groupedProducts);
        } catch (e) {
          console.error("Error fetching products:", e);
          reject(error(500, "Internal Server Error"));
        }
      })
    }
  };
};


