open SvelteStore

let testStore = Writable.make(0)
let testStoreD = Derived.fromW(testStore, i => i * 2)

let reset = () => testStore.set(0)

let log = () => {
  let current = Writable.get(testStore)
  Console.log2("testStore.value", current)
}
