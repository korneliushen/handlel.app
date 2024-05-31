<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
</script>


<main class="flex min-h-screen flex-col items-center overflow-y-scroll font-sans">
	<Header />
	<slot />

<main class="flex font-sans items-center min-h-screen flex-col overflow-y-scrol relative">
  <Header/>
  <div class=" mt-1 sm:mt-0">
    <slot />
  </div>

</main>
