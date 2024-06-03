<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
</script>

<main class="flex absolute inset-0 font-sans items-center min-h-screen overflow-x-hidden overflow-y-scroll flex-col">
  <Header/>
  <div class=" mt-1 sm:mt-0">
    <slot />
  </div>

</main>
