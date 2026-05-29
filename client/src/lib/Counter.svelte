<script>
  // import { store } from "../rs/TestServerStateStore.res.mjs"
  import { testStore, testStoreD, reset, log } from "../rs/TestStore.res.mjs"
  import { onMount } from "svelte"
  import { Fields } from "../rs/StateStore.res.mjs"

  onMount(() => {
    Fields.fetch(console.log, console.error)
  })
</script>

<button
  onclick={() => log()}
>
  {$testStore} / {$testStoreD}
</button>
<button
  onclick={() => $testStore = $testStore + 1}
>
  +1
</button>
<button
  onclick={() => reset()}
>
  reset
</button>

<!--
{#each Object.entries($store) as [key, value]}
  {key}: {value}
  <br>
{:else}
  empty
{/each}
-->
