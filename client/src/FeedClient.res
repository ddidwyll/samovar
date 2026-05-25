module StringMap = Belt.Map.String

type patch<'value> = {
  key: string,
  value: 'value,
}

type source
type event

@new external makeEventSource: string => source = "EventSource"
@send external closeEventSource: source => unit = "close"
@send external addEventListener: (source, string, event => unit) => unit = "addEventListener"
@get external eventData: event => string = "data"

type json = JSON.t
type decoder<'value> = json => option<'value>
type patcher<'value> = patch<'value> => unit
type errorHandler = option<string => unit>

type config<'value> = {
  url: string,
  decodeValue: decoder<'value>,
  onPatch: patcher<'value>,
  onError: errorHandler,
}

type t<'value> = {
  source: source,
  _phantom: option<'value>,
}

let decodeString = JSON.Decode.string

let decodePatch = (json: json, decodeValue: decoder<'value>): option<patch<'value>> => {
  switch JSON.Decode.object(json) {
  | None => None
  | Some(obj) =>
    switch (
      Dict.get(obj, "key")->Option.flatMap(decodeString),
      Dict.get(obj, "value")->Option.flatMap(decodeValue),
    ) {
    | (Some(key), Some(value)) => Some({key, value})
    | _ => None
    }
  }
}

let start = (
  ~url: string,
  ~decodeValue: decoder<'value>,
  ~onPatch: patcher<'value>,
  ~onError: errorHandler,
): t<'value> => {
  let source = makeEventSource(url)

  let handleError = error => onError->Option.forEach(cb => cb(error))

  addEventListener(source, "change", evt => {
    switch eventData(evt)->JSON.parseOrThrow {
    | exception _ => handleError("Failed to parse patch json")
    | json =>
      switch decodePatch(json, decodeValue) {
      | Some(patch) => onPatch(patch)
      | None => handleError("Unexpected patch format")
      }
    }
  })

  {source, _phantom: None}
}

let stop = (client: t<'value>) => {
  closeEventSource(client.source)
}
