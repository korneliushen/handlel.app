<script lang="ts">
	import type { local_products } from '$lib/types/extendedPrisma';
	import { handlelapp } from '$lib/stores/handlelapp';

	export function addProduct(products: local_products) {
		$handlelapp.push(products);
	}
	export function removeProduct(products: local_products, id: number) {
		let handlelapp: local_products[] = JSON.parse(localStorage.getItem('handlelapp') as string);
		handlelapp.splice(id, 1);
		localStorage.setItem('handlelapp', JSON.stringify(handlelapp));
	}
</script>
