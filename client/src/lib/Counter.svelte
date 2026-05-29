<script>
  import { store } from "../rs/StateStore.res.mjs"
</script>

<!--
{#each Object.entries($store) as [key, value]}
  {key}: {value}
  <br>
{:else}
  empty
{/each}
-->
<pre>{JSON.stringify($store, null, 2)}</pre>
