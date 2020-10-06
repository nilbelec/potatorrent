<script>
  import Logo from "./Logo.svelte";
  import NavMenu from "./NavMenu.svelte";
</script>

<style>
  header {
    position: fixed;
    width: 100vw;
    height: 3.5rem;
    background-color: white;
    z-index: 100;
    box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0, 0, 0, 0.5);
    opacity: 0.95;
    transition: opacity 0.3s;
  }

  header:hover {
    opacity: 1;
  }

  nav {
    height: 100%;
    align-items: center;
    justify-content: space-between;
    display: flex;
    padding: 0 4rem;
  }
</style>

<header>
  <nav>
    <Logo />
    <NavMenu />
  </nav>
</header>
