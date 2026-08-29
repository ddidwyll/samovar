module Err = {
  type from = Exn | Json | Manual

  type t = {
    msg: string,
    from: from,
  }

  let make = (msg, ~from=Manual): t => {msg, from}

  let buildDefault = (defaultErr): string => {
    defaultErr->Option.getOr("Unexpexted")
  }

  let catch = (exn, ~defaultErr=?): t => {
    switch JsExn.message(exn) {
    | Some(msg) => make(msg, ~from=Exn)
    | None => defaultErr->buildDefault->make(~from=Exn)
    }
  }

  let parse = (json, ~defaultErr=?): t => {
    JSON.Decode.object(json)
    ->Option.flatMap(Dict.get(_, "msg"))
    ->Option.flatMap(JSON.Decode.string)
    ->Option.getOr(buildDefault(defaultErr))
    ->make(~from=Json)
  }

  let string = (err: t): string => err.msg
}
