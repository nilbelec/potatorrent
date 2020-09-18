<script>
  import Icon from "svelte-awesome";
  import { github } from "svelte-awesome/icons";
</script>

<style>
  ul {
    list-style: none;
    line-height: 1;
    display: flex;
    align-items: center;
  }

  li {
    display: inline;
    margin: 0 1rem;
  }

  a {
    color: inherit;
    border-bottom: none;
    transition: none;
  }
</style>

<ul>
  <li>Búsquedas Guardadas</li>
  <li>Configuración</li>
  <li>
    <a href="https://github.com/nilbelec/potatorrent/" target="_blank">
      <Icon data={github} scale="2" />
    </a>
  </li>
</ul>
