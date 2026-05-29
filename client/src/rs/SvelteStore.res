type invalidator = unit => unit
type subscriber<'a> = 'a => unit
type updater<'a> = 'a => 'a

type unsubscribe = unit => unit
type subscribe<'a> = subscriber<'a> => unsubscribe

type readable<'a> = {subscribe: subscribe<'a>}

type writable<'a> = {
  subscribe: subscribe<'a>,
  set: 'a => unit,
  update: ('a => 'a) => unit,
}

@module("svelte/store")
external writable: 'a => writable<'a> = "writable"

@module("svelte/store")
external readable: ('a, subscriber<'a> => unsubscribe) => readable<'a> = "readable"

@module("svelte/store")
external derivedR: (readable<'a>, 'a => 'b) => readable<'b> = "derived"

@module("svelte/store")
external derivedRR: ((readable<'a>, readable<'b>), ('a => 'b) => 'c) => readable<'c> = "derived"

@module("svelte/store")
external derivedW: (writable<'a>, 'a => 'b) => readable<'b> = "derived"

@module("svelte/store")
external derivedWW: ((writable<'a>, writable<'b>), ('a => 'b) => 'c) => readable<'c> = "derived"

@module("svelte/store")
external derivedRW: ((readable<'a>, writable<'b>), ('a => 'b) => 'c) => readable<'c> = "derived"

@module("svelte/store")
external readonlyW: writable<'a> => readable<'a> = "get"

@module("svelte/store")
external getW: writable<'a> => 'a = "get"

@module("svelte/store")
external getR: readable<'a> => 'a = "get"

let subscribeR = (store: readable<'a>, callback: 'a => unit): unsubscribe =>
  store.subscribe(callback)

let subscribeW = (store: writable<'a>, callback: 'a => unit): unsubscribe =>
  store.subscribe(callback)

let setW = (store: writable<'a>, value: 'a) => store.set(value)

let updateW = (store: writable<'a>, fn: updater<'a>) => store.update(fn)

module Writable = {
  type t<'a> = writable<'a>

  let make = writable
  let get = getW
  let set = setW
  let update = updateW
  let subscribe = subscribeW
  let readonly = readonlyW
}

module Readable = {
  type t<'a> = readable<'a>

  let get = getR
  let subscribe = subscribeR
}

module Derived = {
  let fromR = derivedR
  let fromRR = derivedRR
  let fromW = derivedW
  let fromWW = derivedWW
  let fromRW = derivedRW
}
