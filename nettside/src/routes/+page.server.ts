import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { products } from "@prisma/client";

export const load: PageServerLoad = async ({ url }) => {
  try {
    // Fetch all categories
    const categories: { category: string }[] = await prisma.products.findMany({
      select: {
        category: true,
      },
      distinct: ['category']
    });

    // Fetch 4 products from each category and group them in JSON with an array for each category
    const groupedProducts: { [key: string]: products[] } = {};
    for (const { category } of categories) {
      const products: products[] = await prisma.products.findMany({
        where: { category },
        take: 4,
      });
      groupedProducts[category] = products;
    }

    return { products: groupedProducts };

  } catch (e) {
    console.error("Error fetching products:", e);
    throw error(500, "Internal Server Error");
  }
};

