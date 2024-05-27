<script>
  import "../app.css"
  import Header from "$lib/components/Header.svelte";
</script>

<main class="flex items-center min-h-screen flex-col overflow-y-scroll">
  <Header/>
	<slot />
</main>

