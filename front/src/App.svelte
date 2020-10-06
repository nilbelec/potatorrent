<script>
  import { Router } from "svelte-routing";
  import Navbar from "./navbar/Navbar.svelte";
  import Content from "./content/Content.svelte";
  export let url = "";
</script>

<Router {url}>
  <Navbar />
  <Content />
</Router>
