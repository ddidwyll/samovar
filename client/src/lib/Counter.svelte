<script>
  // import { onMount } from "svelte"
  // import { SvelteSet } from "svelte/reactivity"

  // let connected = $state(false)
  // let events = new SvelteSet()

  // const addEvent = (name, data) => {
  //   console.log(name)
  //   events.add({name, data})
  // }

  // onMount(() => {
  //   const es = new EventSource("//localhost:4000/feed")

  //   es.onopen = () => connected = true
  //   es.onerror = () => connected = false

  //   es.addEventListener("welcome", (e) => {
  //     addEvent("welcome", e.data)
  //   })

  //   es.addEventListener("change", (e) => {
  //     addEvent("change", JSON.parse(e.data))
  //   })
  // })

  import { make } from "../SSEStore.res.mjs"

  const events = make("//localhost:4000/store/feed")
</script>

<!-- connected: {connected}

<br>

{#each events as {name, data}}
  Name: {name}
  <pre>{JSON.stringify(data, null, 2)}</pre>
{:else}
  no events
{/each}
-->
<pre>{JSON.stringify($events, null, 2)}</pre>
