import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { products } from "@prisma/client";
import type { ExtendedProduct } from "$lib/types/extendedPrisma";

export const load: PageServerLoad = async ({ setHeaders }) => {
  setHeaders({
      'cache-control': 'max-age=86400'
  });

  try {
    // Fetch all categories
    const categories: { category: string }[] = await prisma.products.findMany({
      select: {
        category: true,
      },
      distinct: ['category']
    });

    // Fetch 4 products from each category and group them in JSON with an array for each category
    const groupedProducts: { [key: string]: ExtendedProduct[] } = {};
    for (const { category } of categories) {
      const products = await prisma.products.findMany({
        where: { category },
        take: 4,
      }) as ExtendedProduct[];
      groupedProducts[category] = products;
    }

    return { products: groupedProducts };

  } catch (e) {
    console.error("Error fetching products:", e);
    throw error(500, "Internal Server Error");
  }
};

