<script>
  import { store } from "../rs/TestServerStateStore.res.mjs"
</script>

{#each Object.entries($store) as [key, value]}
  {key}: {value}
  <br>
{:else}
  empty
{/each}
