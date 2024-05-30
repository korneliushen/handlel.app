import type { products } from "@prisma/client";

export function addProduct(products: products) {
    if (!localStorage.getItem("handlelapp")) {
        localStorage.setItem("handlelapp", JSON.stringify([]))
    }

    let handlelapp: products[] = JSON.parse(localStorage.getItem("handlelapp") as string)
    handlelapp.push(products)
    localStorage.setItem("handlelapp", JSON.stringify(handlelapp))
}
export function removeProduct(products: products, id: number) {
    let handlelapp: products[] = JSON.parse(localStorage.getItem("handlelapp") as string)
    handlelapp.splice(id, 1)
    localStorage.setItem("handlelapp", JSON.stringify(handlelapp))
}