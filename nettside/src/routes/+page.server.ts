import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { products } from "@prisma/client";


export const load: PageServerLoad = async ({ url }) => {
  try {
    // Fetch alle kategorier
    const categories: { category: string }[] = await prisma.products.findMany({
      select: {
        category: true,
      },
      distinct: ['category']
    });

    // Fetch 4 produkter fra hver kategori og grupper i JSON med array for hver kategori
    const groupedProducts: { [key: string]: products[] } = {};
    for (const { category } of categories) {
      const products: products[] = await prisma.products.findMany({
        where: { category },
        take: 4,
      });
      groupedProducts[category] = products;


    return { products: groupedProducts };

  } catch (e) {
    console.error("Error fetching products:", e);
    error(500, "Internal Server Error")
  }
};

