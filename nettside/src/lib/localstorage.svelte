<script lang="ts">
	import type { products } from '@prisma/client';
	import { handlelapp } from '$lib/stores/handlelapp';

	export function addProduct(products: products) {
		$handlelapp.push(products);
	}
	export function removeProduct(products: products, id: number) {
		let handlelapp: products[] = JSON.parse(localStorage.getItem('handlelapp') as string);
		handlelapp.splice(id, 1);
		localStorage.setItem('handlelapp', JSON.stringify(handlelapp));
	}
</script>
