<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
</script>

<main class="flex font-sans items-center min-h-screen flex-col overflow-y-scrol relative">
  <Header/>
  <div class=" mt-1 sm:mt-0">
    <slot />
  </div>

</main>
