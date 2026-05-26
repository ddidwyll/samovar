type source
type event

@new external makeEventSource: string => source = "EventSource"
@send external closeEventSource: source => unit = "close"
@send external addEventListener: (source, string, event => unit) => unit = "addEventListener"
@get external eventData: event => string = "data"

type patch<'value> = {key: string, value: 'value}
type patches<'value> = array<patch<'value>>

type json = JSON.t
type valueDecoder<'value> = json => option<'value>
type statePatcher<'value> = patches<'value> => unit
type errorHandler = option<string => unit>

type patchDecoder<'value> = (json, valueDecoder<'value>) => option<patch<'value>>
type patchesDecoder<'value> = (json, valueDecoder<'value>) => option<patches<'value>>

type config<'value> = {
  url: string,
  decodeValue: valueDecoder<'value>,
  onPatch: statePatcher<'value>,
  onError: errorHandler,
}

type t<'value> = {
  source: source,
  _phantom: option<'value>,
}

let decodeString = JSON.Decode.string

let decodePatch: patchDecoder<'value> = (json, decodeValue) => {
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

let decodePatches: patchesDecoder<'value> = (json, decodeValue) => {
  switch json {
  | JSON.Object(_) => decodePatch(json, decodeValue)->Option.map(p => [p])
  | JSON.Array(patches) => Array.filterMap(patches, p => decodePatch(p, decodeValue))->Some
  | _ => None
  }
}

let start = (
  ~url: string,
  ~decodeValue: valueDecoder<'value>,
  ~onPatch: statePatcher<'value>,
  ~onError: errorHandler,
): t<'value> => {
  let source = makeEventSource(url)

  let handleError = error => onError->Option.forEach(cb => cb(error))

  addEventListener(source, "change", evt => {
    switch eventData(evt)->JSON.parseOrThrow {
    | exception _ => handleError("Failed to parse patch json")
    | json =>
      switch decodePatches(json, decodeValue) {
      | Some(patches) => onPatch(patches)
      | None => handleError("Unexpected patch format")
      }
    }
  })

  {source, _phantom: None}
}

let stop = (client: t<'value>) => {
  closeEventSource(client.source)
}
