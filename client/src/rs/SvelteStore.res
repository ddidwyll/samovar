type subscriber<'a> = 'a => unit
type unsubscribe = unit => unit

module Readable = {
  type t<'a> = {subscribe: subscriber<'a> => unsubscribe}

  type setFn<'a> = 'a => unit
  type updateFn<'a> = ('a => 'a) => unit
  type startFn<'a> = (setFn<'a>, updateFn<'a>) => unsubscribe

  @module("svelte/store")
  external make: ('a, startFn<'a>) => t<'a> = "readable"

  @module("svelte/store")
  external get: t<'a> => 'a = "get"
}
