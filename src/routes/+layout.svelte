<script>
  import "../app.css"
  import Header from "$lib/components/Header.svelte";
</script>

<main class="flex items-center h-screen flex-col">
  <Header/>
	<slot />
</main>

