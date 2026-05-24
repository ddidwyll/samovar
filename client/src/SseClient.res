module StringMap = Belt.Map.String

type patch<'value> = {
  key: string,
  value: 'value,
}

type source
type event

@new external makeEventSource: string => source = "EventSource"
@send external closeEventSource: source => unit = "close"
@send external addEventListener: (source, string, event => unit) => unit = "addEventListene"
@get external eventData: event => string = "data"

type json = JSON.t
type decoder<'value> = json => option<'value>
type patcher<'value> = patch<'value> => unit

type config<'value> = {
  url: string,
  decoder: decoder<'value>,
  patcher: patcher<'value>,
  onError: option<string => unit>,
}

type t<'value> = {
  source: source,
  _phantom: option<'value>,
}

let decodeString = JSON.Decode.string

let decodePatch = (json: json, valueDecoder: decoder<'value>): option<patch<'value>> => {
  switch JSON.Decode.object(json) {
  | None => None
  | Some(obj) =>
    switch (
      Dict.get(obj, "key")->Option.flatMap(decodeString),
      Dict.get(obj, "value")->Option.flatMap(valueDecoder),
    ) {
    | (Some(key), Some(value)) => Some({key, value})
    | _ => None
    }
  }
}
