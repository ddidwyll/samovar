module Err = {
  type t = {msg: string}

  let make = (msg): t => {msg: msg}

  let buildDefault = (defaultErr): string => {
    defaultErr->Option.getOr("Unexpexted")
  }

  let catch = (exn, ~defaultErr=?): t => {
    switch JsExn.message(exn) {
    | Some(msg) => make(msg)
    | None => defaultErr->buildDefault->make
    }
  }

  let parse = (json, ~defaultErr=?): t => {
    JSON.Decode.object(json)
    ->Option.flatMap(obj => Dict.get(obj, "msg"))
    ->Option.flatMap(JSON.Decode.string)
    ->Option.getOr(buildDefault(defaultErr))
    ->make
  }
}
