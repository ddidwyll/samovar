<script>
  import { make } from "../SvelteStore.res.mjs"

  const count = make(0)
</script>

<button class="counter" onclick={() => $count = $count + 1}>
  Count is {$count}
</button>
