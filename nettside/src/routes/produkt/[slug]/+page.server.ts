import prisma from "$lib/prisma";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
    const product = await prisma.products.findFirst(
        {
            where: {id: params.slug}
        }
    )
    if (!product) error(404)
    return {product}
        
};