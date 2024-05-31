import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { products } from '@prisma/client';

export const handlelapp = writable<products[]>(
	browser ? JSON.parse(window.localStorage.getItem('handlelapp') ?? '[]') : []
);

if (browser) {
	handlelapp.subscribe((value) => {
		console.log('Store updated:', value);
		window.localStorage.setItem('handlelapp', JSON.stringify(value));
	});
}
